package db

import (
	"fmt"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data
type ReadKeeper interface {
	BackingQueueHead(ctx sdk.Context) (ts.Backing, sdk.Error)
	BackQueueLen(ctx sdk.Context) int
	GetBacking(ctx sdk.Context, id int64) (ts.Backing, sdk.Error)
	GetStory(ctx sdk.Context, storyID int64) (ts.Story, sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
// to truchain data
type WriteKeeper interface {
	BackingQueuePop(ctx sdk.Context) (ts.Backing, sdk.Error)
	BackingQueuePush(ctx sdk.Context, id int64)
	NewBacking(
		ctx sdk.Context,
		storyID int64,
		amount sdk.Coin,
		creator sdk.AccAddress,
		duration time.Duration,
	) (int64, sdk.Error)
	NewResponseEndBlock(ctx sdk.Context) abci.ResponseEndBlock
	NewStory(
		ctx sdk.Context,
		body string,
		category ts.StoryCategory,
		creator sdk.AccAddress,
		storyType ts.StoryType) (int64, sdk.Error)
}

// Keeper defines a module interface that facilities read/write access
// to truchain data
type Keeper interface {
	ReadKeeper
	WriteKeeper
}

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

// id creates a new id for a key by incrementing the last one by 1
func (k TruKeeper) id(ctx sdk.Context, storeKey sdk.StoreKey) int64 {
	store := ctx.KVStore(storeKey)

	// create key of form "keyName:len", i.e: "stories:len"
	key := storeKey.Name() + ":len"
	keyVal := []byte(key)
	lastID := store.Get(keyVal)

	// if we don't have an ID yet, set it to one
	if lastID == nil {
		one := k.cdc.MustMarshalBinary(int64(1))
		store.Set(keyVal, one)

		return 1
	}

	// convert from binary to int64
	ID := new(int64)
	k.cdc.MustUnmarshalBinary(lastID, ID)

	// increment id by 1 and update value in the kvstore
	newID := *ID + 1
	newIDVal := k.cdc.MustMarshalBinary(newID)
	store.Set(keyVal, newIDVal)

	return newID
}

// key creates a key of the form keyName:id
func key(keyName string, id int64) []byte {
	return []byte(fmt.Sprintf("%s:%d", keyName, id))
}
