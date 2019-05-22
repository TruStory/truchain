package category

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for categories
	StoreKey = "categories"
)

// ReadKeeper defines a module interface that facilitates read only access
type ReadKeeper interface {
	app.ReadKeeper

	GetCategory(ctx sdk.Context, id int64) (Category, sdk.Error)
	GetAllCategories(ctx sdk.Context) ([]Category, sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
type WriteKeeper interface {
	ReadKeeper

	Create(ctx sdk.Context, title string, slug string, description string) int64
	AddToTotalCred(ctx sdk.Context, id int64, amt sdk.Coin) sdk.Error
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(storeKey sdk.StoreKey, codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec, storeKey)}
}

// Create adds a story to the key-value store
func (k Keeper) Create(
	ctx sdk.Context,
	title string,
	slug string,
	description string) int64 {

	logger := ctx.Logger().With("module", "category")

	cat := Category{
		k.GetNextID(ctx),
		title,
		slug,
		description,
		sdk.NewCoin(slug, sdk.ZeroInt()),
		app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setCategory(ctx, cat)

	logger.Info("Created " + cat.String())

	return cat.ID
}

// GetCategory gets the category with the given id from the key-value store
func (k Keeper) GetCategory(ctx sdk.Context, id int64) (cat Category, err sdk.Error) {
	store := k.GetStore(ctx)
	val := store.Get(k.GetIDKey(id))
	if val == nil {
		return cat, ErrCategoryNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &cat)

	return
}

// GetAllCategories gets all categories from the key-value store
func (k Keeper) GetAllCategories(ctx sdk.Context) (cats []Category, err sdk.Error) {
	cat := Category{}
	err = k.Each(ctx, func(val []byte) bool {
		k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &cat)
		cats = append(cats, cat)
		return true
	})
	return
}

// AddToTotalCred updates the total supply of cred for the key-value store
func (k Keeper) AddToTotalCred(ctx sdk.Context, id int64, amt sdk.Coin) (err sdk.Error) {
	cat, err := k.GetCategory(ctx, id)
	if err != nil {
		return ErrCategoryNotFound(id)
	}

	if amt.Denom != cat.Denom() {
		return ErrCodeDenomMismatch(id)
	}

	cat.TotalCred = cat.TotalCred.Add(amt)
	k.setCategory(ctx, cat)

	return
}

// ============================================================================

// setCategory saves a `Category` type to the KVStore
func (k Keeper) setCategory(ctx sdk.Context, cat Category) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(cat.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(cat))
}
