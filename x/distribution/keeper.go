package distribution

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cosmosDist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	bankKeeper       BankKeeper
	accountKeeper    auth.AccountKeeper
	supplyKeeper     supply.Keeper
	cosmosDistKeeper cosmosDist.Keeper
}

// NewKeeper creates a new keeper of the auth Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, bankKeeper BankKeeper,
	accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper, cosmosDistKeeper cosmosDist.Keeper) Keeper {

	// ensure distribution module accounts are set
	if addr := supplyKeeper.GetModuleAddress(UserGrowthPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", UserGrowthPoolName))
	}
	if addr := supplyKeeper.GetModuleAddress(UserRewardPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", UserRewardPoolName))
	}
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		bankKeeper,
		accountKeeper,
		supplyKeeper,
		cosmosDistKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
