package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

// Utilities for all module keepers

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data. Modules keepers should implement this base interface.
type ReadKeeper interface {
	GetCodec() *amino.Codec
}

// Keeper data type with a default codec
type Keeper struct {
	codec *amino.Codec
}

// NewKeeper creates a new parent keeper for module keepers to embed
func NewKeeper(codec *amino.Codec) Keeper {
	return Keeper{codec: codec}
}

// GetCodec returns the base keeper's underlying codec
func (k Keeper) GetCodec() *amino.Codec {
	return k.codec
}

// GetNextID increments and returns the next available id by 1
func (k Keeper) GetNextID(ctx sdk.Context, storeKey sdk.StoreKey) (id int64) {
	store := ctx.KVStore(storeKey)
	lenKey := []byte(storeKey.Name() + ":len")

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

// Marshal marshals a data type into a byte array using the codec
func (k Keeper) Marshal(value interface{}) []byte {
	return k.codec.MustMarshalBinary(value)
}

// GetIDKey returns a store key of form name:id:[ID]
func GetIDKey(storeKey sdk.StoreKey, id int64) []byte {
	return []byte(fmt.Sprintf("%s:id:%d", storeKey.Name(), id))
}
