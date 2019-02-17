package story

import (
	"fmt"
	"net/url"
	"sort"
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for stories
	StoreKey = "stories"
	// QueueStoreKey is string representation of the store key for backing list
	QueueStoreKey = "storyQueue"
)

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	ChallengedStoriesWithCategory(
		ctx sdk.Context,
		catID int64) (stories []Story, err sdk.Error)
	CategoryDenom(ctx sdk.Context, id int64) (name string, err sdk.Error)
	FeedByCategoryID(
		ctx sdk.Context,
		catID int64) (stories []Story, err sdk.Error)
	Stories(ctx sdk.Context) (stories []Story)
	StoriesByCategoryID(ctx sdk.Context, catID int64) (stories []Story, err sdk.Error)
	Story(ctx sdk.Context, storyID int64) (Story, sdk.Error)
	ExpireDuration(ctx sdk.Context) (res time.Duration)
}

// WriteKeeper defines a module interface that facilities read/write access
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context,
		argument string,
		body string,
		categoryID int64,
		creator sdk.AccAddress,
		source url.URL,
		storyType Type) (int64, sdk.Error)
	StartGame(ctx sdk.Context, storyID int64) sdk.Error
	EndGame(ctx sdk.Context, storyID int64, confirmed bool) sdk.Error
	ExpireGame(ctx sdk.Context, storyID int64) sdk.Error
	UpdateStory(ctx sdk.Context, story Story)
	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
	SetParams(ctx sdk.Context, params Params)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	storyQueueKey  sdk.StoreKey
	categoryKeeper category.ReadKeeper
	paramStore     params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	storyQueueKey sdk.StoreKey,
	categoryKeeper category.ReadKeeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		storyQueueKey,
		categoryKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

// ============================================================================

// StartGame records challenging a story
func (k Keeper) StartGame(
	ctx sdk.Context, storyID int64) sdk.Error {

	// get story
	story, err := k.Story(ctx, storyID)
	if err != nil {
		return err
	}

	// update story state
	story.State = Challenged
	k.UpdateStory(ctx, story)

	// add story to challenged list
	k.appendStoriesList(
		ctx, storyIDsByCategoryKey(k, story.CategoryID, story.Timestamp, true), story)

	return nil
}

// EndGame records the end of a validation game on a story
func (k Keeper) EndGame(ctx sdk.Context, storyID int64, confirmed bool) sdk.Error {
	// get story
	story, err := k.Story(ctx, storyID)
	if err != nil {
		return err
	}

	// update story state
	if confirmed {
		story.State = Confirmed
	} else {
		story.State = Rejected
	}
	k.UpdateStory(ctx, story)

	return nil
}

// ExpireGame resets a story after a game has expired
func (k Keeper) ExpireGame(ctx sdk.Context, storyID int64) sdk.Error {
	// get story
	story, err := k.Story(ctx, storyID)
	if err != nil {
		return err
	}

	// update story state
	story.State = Expired
	k.UpdateStory(ctx, story)

	return nil
}

// Create adds a story to the key-value store
func (k Keeper) Create(
	ctx sdk.Context,
	argument string,
	body string,
	categoryID int64,
	creator sdk.AccAddress,
	source url.URL,
	storyType Type) (int64, sdk.Error) {

	logger := ctx.Logger().With("module", "story")

	_, err := k.categoryKeeper.GetCategory(ctx, categoryID)
	if err != nil {
		return 0, category.ErrInvalidCategory(categoryID)
	}

	story := Story{
		ID:         k.GetNextID(ctx),
		Argument:   argument,
		Body:       body,
		CategoryID: categoryID,
		Creator:    creator,
		ExpireTime: ctx.BlockHeader().Time.Add(k.ExpireDuration(ctx)),
		Flagged:    false,
		GameID:     0,
		Source:     source,
		State:      Unconfirmed,
		Type:       storyType,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setStory(ctx, story)
	k.appendStoriesList(
		ctx, storyIDsByCategoryKey(k, categoryID, story.Timestamp, false), story)

	k.storyQueue(ctx).Push(story.ID)

	logger.Info("Created " + story.String())

	return story.ID, nil
}

// CategoryDenom returns the name of the category coin for the story
func (k Keeper) CategoryDenom(ctx sdk.Context, id int64) (name string, err sdk.Error) {
	story, err := k.Story(ctx, id)
	if err != nil {
		return
	}
	cat, err := k.categoryKeeper.GetCategory(ctx, story.CategoryID)
	if err != nil {
		return
	}

	return cat.Denom(), nil
}

// Story gets the story with the given id from the key-value store
func (k Keeper) Story(
	ctx sdk.Context, storyID int64) (story Story, err sdk.Error) {

	store := k.GetStore(ctx)
	val := store.Get(k.GetIDKey(storyID))
	if val == nil {
		return story, ErrStoryNotFound(storyID)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &story)

	return
}

// StoriesByCategoryID gets the stories for a given category id
func (k Keeper) StoriesByCategoryID(
	ctx sdk.Context, catID int64) (stories []Story, err sdk.Error) {

	return k.storiesByCategoryID(
		ctx, storyIDsByCategorySubspaceKey(k, catID, false), catID)
}

// ChallengedStoriesWithCategory gets all challenged stories for a category
func (k Keeper) ChallengedStoriesWithCategory(
	ctx sdk.Context, catID int64) (stories []Story, err sdk.Error) {

	return k.storiesByCategoryID(
		ctx, storyIDsByCategorySubspaceKey(k, catID, true), catID)
}

// FeedByCategoryID gets stories ordered by challenged stories first
func (k Keeper) FeedByCategoryID(
	ctx sdk.Context,
	catID int64) (stories []Story, err sdk.Error) {

	// get all story ids by category
	storyIDs, err := k.storyIDsByCategoryID(
		ctx, storyIDsByCategorySubspaceKey(k, catID, false), catID)
	if err != nil {
		return
	}

	// get all challenged story ids by category
	challengedStoryIDs, err := k.storyIDsByCategoryID(
		ctx, storyIDsByCategorySubspaceKey(k, catID, true), catID)
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

	return k.storiesByID(ctx, feedIDs)
}

// Stories returns all stories in reverse chronological order
func (k Keeper) Stories(ctx sdk.Context) (stories []Story) {

	// get store
	store := k.GetStore(ctx)

	// builds prefix "stories:id:"
	searchKey := fmt.Sprintf("%s:id:", k.GetStoreKey().Name())
	searchPrefix := []byte(searchKey)

	// setup iterator
	iter := sdk.KVStorePrefixIterator(store, searchPrefix)
	defer iter.Close()

	// iterates through keyspace to find all stories
	for ; iter.Valid(); iter.Next() {
		var story Story
		k.GetCodec().MustUnmarshalBinaryLengthPrefixed(
			iter.Value(), &story)
		stories = append(stories, story)
	}

	// sort in reverse chronological order
	sort.Slice(stories, func(i, j int) bool {
		return stories[i].ID > stories[j].ID
	})

	return stories
}

// UpdateStory updates an existing story in the store
func (k Keeper) UpdateStory(ctx sdk.Context, story Story) {
	newStory := Story{
		story.ID,
		story.Argument,
		story.Body,
		story.CategoryID,
		story.Creator,
		story.ExpireTime,
		story.Flagged,
		story.GameID,
		story.Source,
		story.State,
		story.Type,
		story.Timestamp.Update(ctx.BlockHeader()),
	}

	k.setStory(ctx, newStory)
}

// ============================================================================

func (k Keeper) appendStoriesList(
	ctx sdk.Context, key []byte, story Story) {

	// get stories store
	store := k.GetStore(ctx)

	// marshal story id to list
	store.Set(
		key,
		k.GetCodec().MustMarshalBinaryBare(story.ID))
}

// setStory saves a `Story` type to the KVStore
func (k Keeper) setStory(ctx sdk.Context, story Story) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(story.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(story))
}

func (k Keeper) storiesByCategoryID(
	ctx sdk.Context,
	prefix []byte,
	catID int64) (stories []Story, err sdk.Error) {

	storyIDs, err := k.storyIDsByCategoryID(ctx, prefix, catID)
	if err != nil {
		return
	}

	if len(storyIDs) == 0 {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}

	return k.storiesByID(ctx, storyIDs)
}

func (k Keeper) storiesByID(
	ctx sdk.Context, storyIDs []int64) (stories []Story, err sdk.Error) {

	for _, storyID := range storyIDs {
		story, err := k.Story(ctx, storyID)
		if err != nil {
			return stories, ErrStoryNotFound(storyID)
		}
		stories = append(stories, story)
	}

	return
}

func (k Keeper) storyIDsByCategoryID(
	ctx sdk.Context, prefix []byte, catID int64) (storyIDs []int64, err sdk.Error) {

	store := k.GetStore(ctx)

	// iterate over subspace, creating a list of stories
	iter := sdk.KVStoreReversePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var storyID int64
		k.GetCodec().MustUnmarshalBinaryBare(iter.Value(), &storyID)
		storyIDs = append(storyIDs, storyID)
	}

	return storyIDs, nil
}

func (k Keeper) storyQueue(ctx sdk.Context) queue.Queue {
	store := ctx.KVStore(k.storyQueueKey)
	return queue.NewQueue(k.GetCodec(), store)
}
