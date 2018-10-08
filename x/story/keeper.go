package story

import (
	t "github.com/TruStory/truchain/types"
	c "github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data
type ReadKeeper interface {
	GetStory(ctx sdk.Context, storyID int64) (Story, sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
// to truchain data
type WriteKeeper interface {
	NewStory(ctx sdk.Context, body string, categoryID int64, creator sdk.AccAddress, kind Kind) (int64, sdk.Error)
}

// ReadWriteKeeper defines a module interface that facilities read/write access
// to truchain data
type ReadWriteKeeper interface {
	ReadKeeper
	WriteKeeper
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	// base keeper
	tk t.Keeper

	// read-only access to category DB
	ck c.ReadKeeper

	// key to access coin store
	bk bank.Keeper

	// unexposed keys to access store from context
	storyKey sdk.StoreKey
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(storyKey sdk.StoreKey, ck c.ReadKeeper, bk bank.Keeper, codec *amino.Codec) Keeper {
	return Keeper{
		tk:       t.NewKeeper(codec),
		storyKey: storyKey,
		ck:       ck,
		bk:       bk,
	}
}

// ============================================================================

// NewStory adds a story to the key-value store
func (k Keeper) NewStory(
	ctx sdk.Context,
	body string,
	categoryID int64,
	creator sdk.AccAddress,
	kind Kind) (int64, sdk.Error) {

	_, err := k.ck.GetCategory(ctx, categoryID)
	if err != nil {
		return 0, c.ErrInvalidCategory(categoryID)
	}

	story := Story{
		ID:           k.tk.GetNextID(ctx, k.storyKey),
		Body:         body,
		CategoryID:   categoryID,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		State:        Created,
		Kind:         kind,
	}

	k.setStory(ctx, story)

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
	k.tk.Codec.MustUnmarshalBinary(val, &story)

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
		k.tk.Codec.MustMarshalBinary(story))
}

// getStoryIDKey returns byte array for "stories:id:[ID]"
func getStoryIDKey(k Keeper, id int64) []byte {
	return t.GetIDKey(k.storyKey, id)
}
