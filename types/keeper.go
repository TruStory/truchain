package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

// Utilities for all module keepers

// Keeper data type with a default codec
type Keeper struct {
	Codec *amino.Codec
}

// NewKeeper creates a new parent keeper for module keepers to embed
func NewKeeper(codec *amino.Codec) Keeper {
	return Keeper{Codec: codec}
}

// GetNextID increments and returns the next available id by 1
func (k Keeper) GetNextID(ctx sdk.Context, storeKey sdk.StoreKey) (id int64) {
	store := ctx.KVStore(storeKey)
	lenKey := []byte(storeKey.Name() + ":len")

	bz := store.Get(lenKey)
	if bz == nil {
		one := k.Codec.MustMarshalBinary(int64(1))
		store.Set(lenKey, one)

		return 1
	}

	k.Codec.MustUnmarshalBinary(bz, &id)
	nextID := id + 1
	bz = k.Codec.MustMarshalBinary(nextID)
	store.Set(lenKey, bz)

	return nextID
}

// GetIDKey returns a store key of form name:id:[ID]
func GetIDKey(storeKey sdk.StoreKey, id int64) []byte {
	return []byte(fmt.Sprintf("%s:id:%d", storeKey.Name(), id))
}
