package story

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	GetChallengedStoriesWithCategory(
		ctx sdk.Context,
		catID int64) (stories []Story, err sdk.Error)
	GetCoinName(ctx sdk.Context, id int64) (name string, err sdk.Error)
	GetFeedWithCategory(
		ctx sdk.Context,
		catID int64) (stories []Story, err sdk.Error)
	GetStoriesWithCategory(ctx sdk.Context, catID int64) (stories []Story, err sdk.Error)
	GetStory(ctx sdk.Context, storyID int64) (Story, sdk.Error)
}

// WriteKeeper defines a module interface that facilities read/write access
type WriteKeeper interface {
	ReadKeeper

	NewStory(ctx sdk.Context, body string, categoryID int64, creator sdk.AccAddress, kind Kind) (int64, sdk.Error)
	StartChallenge(ctx sdk.Context, storyID int64) sdk.Error
	UpdateStory(ctx sdk.Context, story Story)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	categoryKeeper                 category.ReadKeeper
	storiesByCategoryKey           sdk.StoreKey
	challengedStoriesByCategoryKey sdk.StoreKey
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	storiesByCategoryKey sdk.StoreKey,
	challengedStoriesByCategoryKey sdk.StoreKey,
	categoryKeeper category.ReadKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		categoryKeeper,
		storiesByCategoryKey,
		challengedStoriesByCategoryKey}
}

// ============================================================================

// StartChallenge records challenging a story
func (k Keeper) StartChallenge(ctx sdk.Context, storyID int64) sdk.Error {
	// get story
	story, err := k.GetStory(ctx, storyID)
	if err != nil {
		return err
	}
	// update story state
	story.State = Challenged
	k.UpdateStory(ctx, story)

	// add story to challenged list
	k.appendStoriesList(ctx, k.challengedStoriesByCategoryKey, story)

	return nil
}

// NewStory adds a story to the key-value store
func (k Keeper) NewStory(
	ctx sdk.Context,
	body string,
	categoryID int64,
	creator sdk.AccAddress,
	kind Kind) (int64, sdk.Error) {

	_, err := k.categoryKeeper.GetCategory(ctx, categoryID)
	if err != nil {
		return 0, category.ErrInvalidCategory(categoryID)
	}

	story := Story{
		k.GetNextID(ctx),
		nil,
		nil,
		nil,
		body,
		categoryID,
		0,
		ctx.BlockHeight(),
		creator,
		0,
		Created,
		kind,
		ctx.BlockHeight(),
		nil,
	}

	k.setStory(ctx, story)
	k.appendStoriesList(ctx, k.storiesByCategoryKey, story)

	return story.ID, nil
}

// GetCoinName returns the name of the category coin for the story
func (k Keeper) GetCoinName(ctx sdk.Context, id int64) (name string, err sdk.Error) {
	story, err := k.GetStory(ctx, id)
	if err != nil {
		return
	}
	cat, err := k.categoryKeeper.GetCategory(ctx, story.CategoryID)
	if err != nil {
		return
	}

	return cat.CoinName(), nil
}

// GetStory gets the story with the given id from the key-value store
func (k Keeper) GetStory(ctx sdk.Context, storyID int64) (story Story, err sdk.Error) {
	store := k.GetStore(ctx)
	val := store.Get(k.GetIDKey(storyID))
	if val == nil {
		return story, ErrStoryNotFound(storyID)
	}
	k.GetCodec().MustUnmarshalBinary(val, &story)

	return
}

// GetStoriesWithCategory gets the stories for a given category id
func (k Keeper) GetStoriesWithCategory(
	ctx sdk.Context,
	catID int64) (stories []Story, err sdk.Error) {

	return k.storiesWithCategory(ctx, k.storiesByCategoryKey, catID)
}

// GetChallengedStoriesWithCategory gets all challenged stories for a category
func (k Keeper) GetChallengedStoriesWithCategory(
	ctx sdk.Context,
	catID int64) (stories []Story, err sdk.Error) {

	return k.storiesWithCategory(ctx, k.challengedStoriesByCategoryKey, catID)
}

// GetFeedWithCategory gets stories ordered by challenged stories first
func (k Keeper) GetFeedWithCategory(
	ctx sdk.Context,
	catID int64) (stories []Story, err sdk.Error) {

	// get all story ids by category
	storyIDs, err := k.storyIDsWithCategory(ctx, k.storiesByCategoryKey, catID)
	if err != nil {
		return
	}

	// get all challenged story ids by category
	challengedStoryIDs, err := k.storyIDsWithCategory(ctx, k.challengedStoriesByCategoryKey, catID)
	if err != nil {
		return
	}

	// make a list of all unchallenged story ids
	var unchallengedStoryIDs []int64
	for _, sid := range storyIDs {
		isMatch := false
		for _, cid := range challengedStoryIDs {
			isMatch = sid == cid
			if isMatch {
				break
			}
		}
		if !isMatch {
			unchallengedStoryIDs = append(unchallengedStoryIDs, sid)
		}
	}

	// concat challenged story ids with unchallenged story ids
	feedIDs := append(challengedStoryIDs, unchallengedStoryIDs...)

	return k.stories(ctx, feedIDs)
}

// UpdateStory updates an existing story in the store
func (k Keeper) UpdateStory(ctx sdk.Context, story Story) {
	newStory := Story{
		story.ID,
		story.BackIDs,
		story.EvidenceIDs,
		story.Thread,
		story.Body,
		story.CategoryID,
		story.ChallengeID,
		story.CreatedBlock,
		story.Creator,
		story.Round,
		story.State,
		story.Kind,
		ctx.BlockHeight(),
		story.Users,
	}

	k.setStory(ctx, newStory)
}

// ============================================================================

func (k Keeper) appendStoriesList(ctx sdk.Context, storeKey sdk.StoreKey, story Story) {
	// get storiesByCategory store
	store := ctx.KVStore(storeKey)

	// create key "categories:id:[CategoryID]:stories:id:[StoryID]"
	key := []byte(fmt.Sprintf(
		"%s:id:%d:%s:id:%d",
		k.categoryKeeper.GetStoreKey().Name(),
		story.CategoryID,
		k.GetStoreKey().Name(),
		story.ID))

	// marshal story id to list
	store.Set(
		key,
		k.GetCodec().MustMarshalBinary(story.ID))
}

// setStory saves a `Story` type to the KVStore
func (k Keeper) setStory(ctx sdk.Context, story Story) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(story.ID),
		k.GetCodec().MustMarshalBinary(story))
}

func (k Keeper) storiesWithCategory(
	ctx sdk.Context,
	storeKey sdk.StoreKey,
	catID int64) (stories []Story, err sdk.Error) {

	storyIDs, err := k.storyIDsWithCategory(ctx, storeKey, catID)
	if err != nil {
		return
	}

	if len(storyIDs) == 0 {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}

	return k.stories(ctx, storyIDs)
}

func (k Keeper) stories(ctx sdk.Context, storyIDs []int64) (stories []Story, err sdk.Error) {
	for _, storyID := range storyIDs {
		story, err := k.GetStory(ctx, storyID)
		if err != nil {
			return stories, ErrStoryNotFound(storyID)
		}
		stories = append(stories, story)
	}

	return
}

func (k Keeper) storyIDsWithCategory(
	ctx sdk.Context,
	storeKey sdk.StoreKey,
	catID int64) (storyIDs []int64, err sdk.Error) {

	store := ctx.KVStore(storeKey)

	// create subspace prefix "categories:id:[CategoryID]:stories:id:"
	prefix := []byte(fmt.Sprintf(
		"%s:id:%d:%s:id:",
		k.categoryKeeper.GetStoreKey().Name(),
		catID,
		k.GetStoreKey().Name()))

	// iterate over subspace, creating a list of stories
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		if bz == nil {
			return storyIDs, ErrStoriesWithCategoryNotFound(catID)
		}
		var storyID int64
		k.GetCodec().MustUnmarshalBinary(bz, &storyID)
		storyIDs = append(storyIDs, storyID)
	}

	return storyIDs, nil
}
