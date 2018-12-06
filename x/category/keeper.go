package category

import (
	"sort"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
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

	InitCategories(
		ctx sdk.Context, creator sdk.AccAddress, categories map[string]string) (err sdk.Error)

	NewCategory(ctx sdk.Context, title string, creator sdk.AccAddress, slug string, description string) (int64, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(storeKey sdk.StoreKey, codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec, storeKey)}
}

// InitCategories creates the initial set of categories
func (k Keeper) InitCategories(
	ctx sdk.Context, creator sdk.AccAddress, categories map[string]string) (err sdk.Error) {

	var keys []string
	for key := range categories {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		title := categories[key]
		slug := key
		_, err = k.NewCategory(ctx, title, creator, slug, "")
		if err != nil {
			return err
		}
	}

	return
}

// NewCategory adds a story to the key-value store
func (k Keeper) NewCategory(
	ctx sdk.Context,
	title string,
	creator sdk.AccAddress,
	slug string,
	description string) (int64, sdk.Error) {

	cat := Category{
		k.GetNextID(ctx),
		creator,
		title,
		slug,
		description,
		app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setCategory(ctx, cat)

	return cat.ID, nil
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

// GetAllCategories gets the category with the given id from the key-value store
func (k Keeper) GetAllCategories(ctx sdk.Context) (cats []Category, err sdk.Error) {
	cat := Category{}
	err = k.Each(ctx, func(val []byte) bool {
		k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &cat)
		cats = append(cats, cat)
		return true
	})
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
