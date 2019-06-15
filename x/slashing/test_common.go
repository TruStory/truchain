package slashing

import (
	"fmt"

	"github.com/TruStory/truchain/x/auth"

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
var _ ClaimKeeper = claimKeeper{}
var _ StakeKeeper = stakeKeeper{}
var _ auth.BankKeeper = bankKeeper{}

type claimKeeper struct {
	Claims []Claim
}

// Claim ...
func (ck claimKeeper) Claim(ctx sdk.Context, id uint64) (claim Claim, err sdk.Error) {
	for _, claim := range ck.Claims {
		if claim.ID == id {
			return claim, nil
		}
	}
	return claim, sdk.NewError(DefaultCodespace, 404, fmt.Sprintf("Claim not found with ID: %d", id))
}

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

type bankKeeper struct {
}

func (bk bankKeeper) NewTransaction(ctx sdk.Context, to sdk.AccAddress, coins sdk.Coins) bool {
	return true
}

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	authKey := sdk.NewKVStoreKey(auth.ModuleName)
	slashKey := sdk.NewKVStoreKey(ModuleName)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(slashKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)

	bankKeeper := bankKeeper{}
	authKeeper := auth.NewKeeper(authKey, paramsKeeper.Subspace(auth.ModuleName), codec, bankKeeper)

	claimKeeper := claimKeeper{
		Claims: []Claim{
			Claim{ID: 1, CommunityID: 1},
			Claim{ID: 2, CommunityID: 2},
		},
	}
	stakeKeeper := stakeKeeper{
		Stakes: []Stake{
			Stake{ID: 1, ClaimID: 1},
			Stake{ID: 2, ClaimID: 2},
			Stake{ID: 3, ClaimID: 2},
		},
	}
	slashKeeper := NewKeeper(slashKey, paramsKeeper.Subspace(ModuleName), codec, claimKeeper, stakeKeeper, authKeeper)

	InitGenesis(ctx, slashKeeper, DefaultGenesisState())

	return ctx, slashKeeper
}
