package claim

import (
	"net/url"

	"github.com/TruStory/truchain/x/community"
	truauth "github.com/TruStory/truchain/x/auth"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	claimKey := sdk.NewKVStoreKey(StoreKey)
	communityKey := sdk.NewKVStoreKey(community.StoreKey)
	truAuthKey := sdk.NewKVStoreKey(truauth.StoreKey)
	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(claimKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(communityKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)

	pk := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)

	communityKeeper := community.NewKeeper(
		communityKey,
		pk.Subspace(community.ModuleName),
		codec)
	community.InitGenesis(ctx, communityKeeper, community.DefaultGenesisState())

	_, err := communityKeeper.NewCommunity(ctx, "Furries", "furry", "")
	if err != nil {
		panic(err)
	}

	accountKeeper := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	truAuthKeeper := truauth.NewKeeper(
		truAuthKey,
		pk.Subspace(truauth.ModuleName),
		codec,
		bankKeeper,
		accountKeeper,
	)

	keeper := NewKeeper(
		claimKey,
		pk.Subspace(ModuleName),
		codec,
		truAuthKeeper,
		communityKeeper,
	)
	InitGenesis(ctx, keeper, DefaultGenesisState())

	return ctx, keeper
}

func fakeClaim(ctx sdk.Context, keeper Keeper) Claim {
	body := "body string ajsdkhfakjsdfhd"
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	source := url.URL{}
	claim, err := keeper.SubmitClaim(ctx, body, uint64(1), creator, source)
	if err != nil {
		panic(err)
	}

	return claim
}
