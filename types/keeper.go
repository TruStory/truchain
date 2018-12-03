package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// Utilities for all module keepers

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data. Modules keepers should implement this base interface.
type ReadKeeper interface {
	GetCodec() *amino.Codec
	GetStoreKey() sdk.StoreKey
	Each(sdk.Context, func([]byte) bool) sdk.Error
}

// WriteKeeper defines an interface for read/write operations on a KVStore
type WriteKeeper interface {
	ReadKeeper

	GetStore(ctx sdk.Context) sdk.KVStore
}

// Keeper data type with a default codec
type Keeper struct {
	codec    *amino.Codec
	storeKey sdk.StoreKey
}

// NewKeeper creates a new parent keeper for module keepers to embed
func NewKeeper(codec *amino.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{codec, storeKey}
}

// GetCodec returns the base keeper's underlying codec
func (k Keeper) GetCodec() *amino.Codec {
	return k.codec
}

// GetStoreKey returns the default store key for the keeper
func (k Keeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}

// GetStore returns the default KVStore for the keeper
func (k Keeper) GetStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.GetStoreKey())
}

// GetNextID increments and returns the next available id by 1
func (k Keeper) GetNextID(ctx sdk.Context) (id int64) {
	store := k.GetStore(ctx)
	lenKey := []byte(k.storeKey.Name() + ":len")

	bz := store.Get(lenKey)
	if bz == nil {
		one := k.GetCodec().MustMarshalBinaryBare(int64(1))
		store.Set(lenKey, one)

		return 1
	}

	k.GetCodec().MustUnmarshalBinaryBare(bz, &id)
	nextID := id + 1
	bz = k.GetCodec().MustMarshalBinaryBare(nextID)
	store.Set(lenKey, bz)

	return nextID
}

// EachPrefix calls `fn` for each record in a store with a given prefix. Iteration will stop if `fn` returns false
func (k Keeper) EachPrefix(ctx sdk.Context, prefix string, fn func([]byte) bool) (err sdk.Error) {
	var val []byte
	store := k.GetStore(ctx)
	iter := store.Iterator(nil, nil)
	if prefix != "" {
		iter = sdk.KVStorePrefixIterator(store, []byte(prefix))
	}

	for iter.Valid() {
		val = iter.Value()
		if len(val) > 1 {
			if fn(val) == false {
				break
			}
		}
		iter.Next()
	}

	iter.Close()
	return
}

// Each calls `EachPrefix` with an empty prefix
func (k Keeper) Each(ctx sdk.Context, fn func([]byte) bool) (err sdk.Error) {
	return k.EachPrefix(ctx, "", fn)
}

// GetIDKey returns the key for a given index
func (k Keeper) GetIDKey(id int64) []byte {
	return []byte(fmt.Sprintf("%s:id:%d", k.GetStoreKey().Name(), id))
}
