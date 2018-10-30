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
		one := k.GetCodec().MustMarshalBinary(int64(1))
		store.Set(lenKey, one)

		return 1
	}

	k.GetCodec().MustUnmarshalBinary(bz, &id)
	nextID := id + 1
	bz = k.GetCodec().MustMarshalBinary(nextID)
	store.Set(lenKey, bz)

	return nextID
}

// GetIDKey returns the key for a given index
func (k Keeper) GetIDKey(id int64) []byte {
	return []byte(fmt.Sprintf("%s:id:%d", k.GetStoreKey().Name(), id))
}
