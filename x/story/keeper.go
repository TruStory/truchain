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

	GetStory(ctx sdk.Context, storyID int64) (Story, sdk.Error)
	GetStoriesWithCategory(ctx sdk.Context, catID int64) (stories []Story, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
// to truchain data
type WriteKeeper interface {
	NewStory(ctx sdk.Context, body string, categoryID int64, creator sdk.AccAddress, kind Kind) (int64, sdk.Error)
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
	baseKeeper     app.Keeper
	categoryKeeper category.ReadKeeper
	storyKey       sdk.StoreKey
	catKey         sdk.StoreKey
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(storyKey sdk.StoreKey, catKey sdk.StoreKey, ck category.ReadKeeper, codec *amino.Codec) Keeper {
	return Keeper{
		baseKeeper:     app.NewKeeper(codec),
		storyKey:       storyKey,
		catKey:         catKey,
		categoryKeeper: ck,
	}
}

// ============================================================================

// GetCodec returns the base keeper's underlying codec
func (k Keeper) GetCodec() *amino.Codec {
	return k.baseKeeper.Codec
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
		ID:           k.baseKeeper.GetNextID(ctx, k.storyKey),
		Body:         body,
		CategoryID:   categoryID,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		State:        Created,
		Kind:         kind,
	}

	k.setStory(ctx, story)
	k.addStoryToCategory(ctx, story)

	return story.ID, nil
}

// GetStory gets the story with the given id from the key-value store
func (k Keeper) GetStory(ctx sdk.Context, storyID int64) (story Story, err sdk.Error) {
	store := ctx.KVStore(k.storyKey)
	key := getStoryIDKey(k, storyID)
	val := store.Get(key)
	if val == nil {
		return story, ErrStoryNotFound(storyID)
	}
	k.GetCodec().MustUnmarshalBinary(val, &story)

	return
}

// GetStoriesWithCategory gets the stories for a given category id
func (k Keeper) GetStoriesWithCategory(ctx sdk.Context, catID int64) (stories []Story, err sdk.Error) {

	// get bytes stored at "categories:id:[catID]:stories"
	store := ctx.KVStore(k.catKey)
	bz := store.Get(getCategoryStoriesKey(k, catID))
	if bz == nil {
		return stories, ErrStoriesWithCategoryNotFound(catID)
	}

	// deserialize bytes to story ids
	storyIDs := []int64{}
	k.GetCodec().MustUnmarshalBinary(bz, &storyIDs)

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

// UpdateStory updates an existing story in the store
func (k Keeper) UpdateStory(ctx sdk.Context, story Story) {
	newStory := NewStory(
		story.ID,
		story.BackIDs,
		story.EvidenceIDs,
		story.Thread,
		story.Body,
		story.CategoryID,
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
	store := ctx.KVStore(k.storyKey)
	store.Set(
		getStoryIDKey(k, story.ID),
		k.GetCodec().MustMarshalBinary(story))
}

// addStoryToCategory adds a story id to key "categories:id:[ID]"
func (k Keeper) addStoryToCategory(ctx sdk.Context, story Story) {
	store := ctx.KVStore(k.catKey)
	key := getCategoryStoriesKey(k, story.CategoryID)
	bz := store.Get(key)
	if bz == nil {
		bz = k.GetCodec().MustMarshalBinary([]int64{})
		store.Set(key, bz)
	}

	// get list of story ids from category store
	storyIDs := []int64{}
	k.GetCodec().MustUnmarshalBinary(bz, &storyIDs)

	storyIDs = append(storyIDs, story.ID)
	store.Set(key, k.GetCodec().MustMarshalBinary(storyIDs))
}

// getStoryIDKey returns byte array for "stories:id:[ID]"
func getStoryIDKey(k Keeper, id int64) []byte {
	return app.GetIDKey(k.storyKey, id)
}

// getCategoryStoriesKey returns "categories:id:[`catID`]:stories"
func getCategoryStoriesKey(k Keeper, catID int64) []byte {
	return []byte(fmt.Sprintf("%s:id:%d:%s", k.catKey.Name(), catID, k.storyKey.Name()))

}
