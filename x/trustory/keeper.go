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
