package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewCategory adds a story to the key-value store
func (k TruKeeper) NewCategory(
	ctx sdk.Context,
	name string,
	slug string,
	description string) (int64, sdk.Error) {

	cat := ts.NewCategory(
		k.id(ctx, k.categoryKey),
		name,
		slug,
		description)

	k.setCategory(ctx, cat)

	return cat.ID, nil
}

// GetCategory gets the category with the given id from the key-value store
func (k TruKeeper) GetCategory(ctx sdk.Context, id int64) (ts.Category, sdk.Error) {
	store := ctx.KVStore(k.categoryKey)
	key := key(k.categoryKey.Name(), id)
	val := store.Get(key)
	if val == nil {
		return ts.Category{}, ts.ErrCategoryNotFound(id)
	}
	cat := &ts.Category{}
	k.cdc.MustUnmarshalBinary(val, cat)

	return *cat, nil
}

// ============================================================================

// setCategory saves a `Category` type to the KVStore
func (k TruKeeper) setCategory(ctx sdk.Context, cat ts.Category) {
	store := ctx.KVStore(k.categoryKey)
	store.Set(
		key(k.categoryKey.Name(), cat.ID),
		k.cdc.MustMarshalBinary(cat))
}
