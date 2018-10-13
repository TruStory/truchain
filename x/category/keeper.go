package category

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	GetCategory(ctx sdk.Context, id int64) (Category, sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
// to truchain data
type WriteKeeper interface {
	NewCategory(ctx sdk.Context, name string, slug string, description string) (int64, sdk.Error)
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

	catKey   sdk.StoreKey
	storyKey sdk.StoreKey
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(catKey sdk.StoreKey, storyKey sdk.StoreKey, codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec), catKey, storyKey}
}

// NewCategory adds a story to the key-value store
func (k Keeper) NewCategory(
	ctx sdk.Context,
	title string,
	slug string,
	description string) (int64, sdk.Error) {

	cat := NewCategory(
		k.GetNextID(ctx, k.catKey),
		title,
		slug,
		description)

	k.setCategory(ctx, cat)

	return cat.ID, nil
}

// GetCategory gets the category with the given id from the key-value store
func (k Keeper) GetCategory(ctx sdk.Context, id int64) (cat Category, err sdk.Error) {
	store := ctx.KVStore(k.catKey)
	val := store.Get(getCategoryIDKey(k, id))
	if val == nil {
		return cat, ErrCategoryNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinary(val, &cat)

	return
}

// ============================================================================

// setCategory saves a `Category` type to the KVStore
func (k Keeper) setCategory(ctx sdk.Context, cat Category) {
	store := ctx.KVStore(k.catKey)
	store.Set(
		getCategoryIDKey(k, cat.ID),
		k.GetCodec().MustMarshalBinary(cat))
}

// getCategoryIDKey returns byte array for "categories:id:[ID]"
func getCategoryIDKey(k Keeper, id int64) []byte {
	return app.GetIDKey(k.catKey, id)
}
