package slashing

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

// interface conformance check
var _ StakeKeeper = stakeKeeper{}

type stakeKeeper struct {
	Stakes []Stake
}

// Stake ...
func (sk stakeKeeper) Stake(ctx sdk.Context, id uint64) (stake Stake, err sdk.Error) {
	for _, stake := range sk.Stakes {
		if stake.ID == id {
			return stake, nil
		}
	}
	return stake, sdk.NewError(DefaultCodespace, 404, fmt.Sprintf("Stake not found with ID: %d", id))
}

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	slashKey := sdk.NewKVStoreKey(ModuleName)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(slashKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)
	stakeKeeper := stakeKeeper{
		Stakes: []Stake{
			Stake{ID: 1},
			Stake{ID: 2},
			Stake{ID: 3},
		},
	}
	slashKeeper := NewKeeper(slashKey, paramsKeeper.Subspace(ModuleName), codec, stakeKeeper)

	InitGenesis(ctx, slashKeeper, DefaultGenesisState())

	return ctx, slashKeeper
}