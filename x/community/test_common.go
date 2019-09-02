package community

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	communityKey := sdk.NewKVStoreKey(ModuleName)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(communityKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)
	communityKeeper := NewKeeper(communityKey, paramsKeeper.Subspace(ModuleName), codec)

	admin1 := getFakeAdmin()
	admin2 := getFakeAdmin()
	genesis := DefaultGenesisState()
	genesis.Params.CommunityAdmins = append(genesis.Params.CommunityAdmins, admin1, admin2)
	InitGenesis(ctx, communityKeeper, genesis)

	return ctx, communityKeeper
}

func getFakeAdmin() (address sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return addr
}

func getFakeCommunityParams() (id string, name string, description string) {
	name, id, description = "Randomness", "randomness", "All the random quantum things happen in this community."
	return
}

func getAnotherFakeCommunityParams() (id string, name string, description string) {
	name, id, description = "Space", "space", "Come here for anything you want to learn about the space."
	return
}
