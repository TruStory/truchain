package claim

import (
	"net/url"

	"github.com/TruStory/truchain/x/community"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

// interface conformance check
var _ AccountKeeper = accKeeper{}

type accKeeper struct {
	Jailed bool
}

// IsJailed ...
func (ak accKeeper) IsJailed(ctx sdk.Context, addr sdk.AccAddress) (bool, sdk.Error) {
	return ak.Jailed, nil
}

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	claimKey := sdk.NewKVStoreKey("claim")
	communityKey := sdk.NewKVStoreKey("community")
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

	accountKeeper := accKeeper{
		Jailed: false,
	}

	keeper := NewKeeper(
		claimKey,
		pk.Subspace(ModuleName),
		codec,
		accountKeeper,
		communityKeeper,
	)
	InitGenesis(ctx, keeper, DefaultGenesisState())

	return ctx, keeper
}

func fakeClaim(ctx sdk.Context, keeper Keeper) Claim {
	body := "body string ajsdkhfakjsdfhd"
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	source := url.URL{}
	claim, err := keeper.SubmitClaim(ctx, body, "crypto", creator, source)
	if err != nil {
		panic(err)
	}

	return claim
}
