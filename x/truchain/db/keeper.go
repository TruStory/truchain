package db

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

// TruKeeper data type storing keys to the key-value store
type TruKeeper struct {
	// key to access coin store
	ck bank.Keeper

	// unexposed keys to access store from context
	storyKey   sdk.StoreKey
	voteKey    sdk.StoreKey
	backingKey sdk.StoreKey

	// wire codec for binary encoding/decoding
	cdc *amino.Codec
}

// NewTruKeeper creates a new keeper with write and read access
func NewTruKeeper(storyKey sdk.StoreKey, backingKey sdk.StoreKey, ck bank.Keeper, cdc *amino.Codec) TruKeeper {
	return TruKeeper{
		ck:         ck,
		storyKey:   storyKey,
		backingKey: backingKey,
		cdc:        cdc,
	}
}

// ============================================================================

// newID creates a new id for a key by incrementing the last one by 1
func (k TruKeeper) newID(ctx sdk.Context, storeKey sdk.StoreKey) int64 {
	store := ctx.KVStore(storeKey)

	// create key of form "keyName:lastIndex", i.e: "stories:lastIndex"
	key := storeKey.Name() + ":lastIndex"
	keyVal := []byte(key)
	lastID := store.Get(keyVal)

	// if we don't have an ID yet, set it to zero
	if lastID == nil {
		zero := k.cdc.MustMarshalBinary(int64(0))
		store.Set(keyVal, zero)

		return 0
	}

	// convert from binary to int64
	ID := new(int64)
	k.cdc.MustUnmarshalBinary(lastID, ID)

	// increment id by 1 and update value in blockchain
	newID := *ID + 1
	newIDVal := k.cdc.MustMarshalBinary(newID)
	store.Set(keyVal, newIDVal)

	return newID
}

// generateKey creates a key of the form "keyName"|{id}
func generateKey(keyName string, id int64) []byte {
	var key []byte
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(id))
	key = []byte(keyName)
	key = append(key, idBytes...)
	return key
}
