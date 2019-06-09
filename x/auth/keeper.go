package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	recordkeeper "github.com/shanev/cosmos-record-keeper/recordkeeper"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	recordkeeper.RecordKeeper
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the community Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		recordkeeper.NewRecordKeeper(storeKey, codec),
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}
