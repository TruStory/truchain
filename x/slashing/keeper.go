package slashing

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
)

// Keeper is the model object for the package slashing module
type Keeper struct {
	recordkeeper.RecordKeeper

	codec      *codec.Codec
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the slashing Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		recordkeeper.NewRecordKeeper(storeKey, codec),
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}
