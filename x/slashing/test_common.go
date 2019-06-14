package slashing

import (
	"fmt"
	"time"

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
var _ AppAccountKeeper = appAccountKeeper{}

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

func (sk stakeKeeper) SlashCountByID(ctx sdk.Context, id uint64) (count int, err sdk.Error) {
	stake, err := sk.Stake(ctx, id)
	if err != nil {
		return
	}

	return stake.SlashCount, nil
}

func (sk stakeKeeper) IncrementSlashCount(ctx sdk.Context, id uint64) (err sdk.Error) {
	for i, stake := range sk.Stakes {
		if stake.ID == id {
			sk.Stakes[i].SlashCount++
			return nil
		}
	}

	return sdk.NewError(DefaultCodespace, 404, fmt.Sprintf("Stake not found with ID: %d", id))
}

type appAccountKeeper struct {
	AppAccounts []AppAccount
}

func (aak appAccountKeeper) AppAccount(ctx sdk.Context, address sdk.AccAddress) (appAccount AppAccount, err sdk.Error) {
	for _, appAccount := range aak.AppAccounts {
		if appAccount.Address.Equals(address) {
			return appAccount, nil
		}
	}
	return appAccount, sdk.NewError(DefaultCodespace, 404, fmt.Sprintf("AppAccount not found with address: %d", address))
}

func (aak appAccountKeeper) IsJailed(ctx sdk.Context, address sdk.AccAddress) (isJailed bool, err sdk.Error) {
	appAccount, err := aak.AppAccount(ctx, address)
	if err != nil {
		return
	}

	return appAccount.IsJailed, nil
}

func (aak appAccountKeeper) JailUntil(ctx sdk.Context, address sdk.AccAddress, until time.Time) sdk.Error {
	// nothing here..
	return nil
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
	appAccountKeeper := appAccountKeeper{
		AppAccounts: []AppAccount{
			AppAccount{
				Address: DefaultParams().SlashAdmins[0],
				IsJailed: false,
			},
			AppAccount{
				Address: sdk.AccAddress([]byte{3, 4}),
				IsJailed: true,
			},
		},
	}
	slashKeeper := NewKeeper(slashKey, paramsKeeper.Subspace(ModuleName), codec, stakeKeeper, appAccountKeeper)

	InitGenesis(ctx, slashKeeper, DefaultGenesisState())

	return ctx, slashKeeper
}
