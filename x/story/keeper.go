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

// WriteKeeper defines a module interface that facilities write only access
// to truchain data
type WriteKeeper interface {
	NewStory(ctx sdk.Context, body string, categoryID int64, creator sdk.AccAddress, kind Kind) (int64, sdk.Error)
	StartChallenge(ctx sdk.Context, storyID int64) sdk.Error
	UpdateStory(ctx sdk.Context, story Story)
}

// ReadWriteKeeper defines a module interface that facilities read/write access
// to truchain data
type ReadWriteKeeper interface {
	ReadKeeper
	WriteKeeper
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	categoryKeeper category.ReadKeeper
	catKey         sdk.StoreKey
	challengeKey   sdk.StoreKey
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	catKey sdk.StoreKey,
	challengeKey sdk.StoreKey,
	ck category.ReadKeeper,
	codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec, storeKey), ck, catKey, challengeKey}
}

// ============================================================================

// TODO: run this from end blocker

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
	k.appendChallengedCategoryStoriesList(ctx, story)

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
		ID:           k.GetNextID(ctx),
		Body:         body,
		CategoryID:   categoryID,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		State:        Created,
		Kind:         kind,
	}

	k.setStory(ctx, story)
	k.appendCategoryStoriesList(ctx, story)

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

	// get bytes stored at "categories:id:[catID]:stories"
	store := k.GetStore(ctx)
	bz := store.Get(getCategoryStoriesKey(k, catID))
	if bz == nil {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}

	// return list of stories
	return getStories(ctx, k, bz)
}

// GetChallengedStoriesWithCategory gets all challenged stories for a category
func (k Keeper) GetChallengedStoriesWithCategory(
	ctx sdk.Context,
	catID int64) (stories []Story, err sdk.Error) {

	// get bytes stored at "challenges:categories:id:[catID]:stories"
	store := k.GetStore(ctx)
	bz := store.Get(getChallengedStoriesKey(k, catID))
	if bz == nil {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}

	// return list of stories
	return getStories(ctx, k, bz)
}

// GetFeedWithCategory gets stories ordered by challenged stories first
func (k Keeper) GetFeedWithCategory(
	ctx sdk.Context,
	catID int64) (stories []Story, err sdk.Error) {

	// get story store
	store := k.GetStore(ctx)

	// get bytes stored at "categories:id:[catID]:stories"
	bz := store.Get(getCategoryStoriesKey(k, catID))
	if bz == nil {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}

	var (
		all          List
		challenged   List
		unchallenged List
	)

	// unmarshal bytes to story ids
	k.GetCodec().MustUnmarshalBinary(bz, &all)

	// get bytes stored at "challenges:categories:id:[catID]:stories"
	bz = store.Get(getChallengedStoriesKey(k, catID))
	if bz == nil {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}
	// unmarshal challenged story id list
	k.GetCodec().MustUnmarshalBinary(bz, &challenged)

	// make a list of all unchallenged story ids
	for _, sid := range all {
		isMatch := false
		for _, cid := range challenged {
			isMatch = sid == cid
			if isMatch {
				break
			}
		}
		if !isMatch {
			unchallenged = append(unchallenged, sid)
		}
	}

	// concat challenged story ids with unchallenged story ids
	feed := append(challenged, unchallenged...)

	return getStoriesFromIDList(ctx, k, feed)
}

// UpdateStory updates an existing story in the store
func (k Keeper) UpdateStory(ctx sdk.Context, story Story) {
	newStory := NewStory(
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
		story.Users)

	k.setStory(ctx, newStory)
}

// ============================================================================

// setStory saves a `Story` type to the KVStore
func (k Keeper) setStory(ctx sdk.Context, story Story) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(story.ID),
		k.GetCodec().MustMarshalBinary(story))
}

// adds a story id to key "categories:id:[ID]:stories"
func (k Keeper) appendCategoryStoriesList(ctx sdk.Context, story Story) {
	k.appendList(ctx, getCategoryStoriesKey(k, story.CategoryID), story.ID)
}

// adds a story id to key "challenges:categories:id:[ID]:stories"
func (k Keeper) appendChallengedCategoryStoriesList(ctx sdk.Context, story Story) {
	k.appendList(ctx, getChallengedStoriesKey(k, story.CategoryID), story.ID)
}

// appendList adds a story id to a list of story id
func (k Keeper) appendList(ctx sdk.Context, key []byte, storyID int64) {
	// get list of story ids from category store for a given key
	store := k.GetStore(ctx)
	bz := store.Get(key)
	var storyList List
	if bz == nil {
		bz = k.GetCodec().MustMarshalBinary(storyList)
		store.Set(key, bz)
	}
	k.GetCodec().MustUnmarshalBinary(bz, &storyList)

	// add the new story id and marshal it back to the store
	storyList = append(storyList, storyID)
	store.Set(key, k.GetCodec().MustMarshalBinary(storyList))
}

// getCategoryStoriesKey returns "categories:id:[`catID`]:stories"
func getCategoryStoriesKey(k Keeper, catID int64) []byte {
	return []byte(fmt.Sprintf(
		"%s:id:%d:%s",
		k.catKey.Name(),
		catID,
		k.GetStoreKey().Name()))
}

// getChallengedStoriesKey returns "challenges:categories:id:[`catID`]:stories"
func getChallengedStoriesKey(k Keeper, catID int64) []byte {
	return []byte(
		fmt.Sprintf(
			"%s:%s:id:%d:%s",
			k.challengeKey.Name(),
			k.catKey.Name(),
			catID,
			k.GetStoreKey().Name()))
}

func getStories(ctx sdk.Context, k Keeper, bz []byte) (stories []Story, err sdk.Error) {
	// deserialize bytes to story ids
	var storyIDs List
	k.GetCodec().MustUnmarshalBinary(bz, &storyIDs)

	// return list of stories
	return getStoriesFromIDList(ctx, k, storyIDs)
}

// extract each story and add to a list
func getStoriesFromIDList(ctx sdk.Context, k Keeper, storyIDs List) (stories []Story, err sdk.Error) {
	// extract each story and add to a list
	for _, id := range storyIDs {
		story, err := k.GetStory(ctx, id)
		if err != nil {
			return stories, ErrStoryNotFound(id)
		}
		stories = append(stories, story)
	}

	// return list of stories
	return
}
