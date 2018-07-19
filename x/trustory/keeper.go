package trustory

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// Keeper data type
type Keeper struct {
	TruStory  sdk.StoreKey
	codespace sdk.CodespaceType
	cdc       *wire.Codec

	ck bank.Keeper
	sm stake.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(TruStory sdk.StoreKey, ck bank.Keeper, sm stake.Keeper, codespace sdk.CodespaceType) Keeper {
	cdc := wire.NewCodec()

	return Keeper{
		TruStory:  TruStory,
		cdc:       cdc,
		ck:        ck,
		sm:        sm,
		codespace: codespace,
	}
}

// NewStoryID created a new id for a story: not sure how this works
func (k Keeper) NewStoryID(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.TruStory)
	bid := store.Get([]byte("TotalID"))
	if bid == nil {
		return 0
	}

	totalID := new(int64)
	err := k.cdc.UnmarshalBinary(bid, totalID)
	if err != nil {
		panic(err)
	}

	return (*totalID + 1)
}
