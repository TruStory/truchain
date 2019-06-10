package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	authKey := sdk.NewKVStoreKey(ModuleName)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)
	authKeeper := NewKeeper(authKey, paramsKeeper.Subspace(ModuleName), codec)

	InitGenesis(ctx, authKeeper, DefaultGenesisState())

	return ctx, authKeeper
}

func getFakeAppAccountParams() (
	privateKey crypto.PrivKey, publicKey crypto.PubKey, address sdk.AccAddress,
	coins sdk.Coins, earnedCoins EarnedCoins,
) {
	privateKey, publicKey, address = getFakeKeyPubAddr()
	coins = getFakeCoins()
	earnedCoins = getFakeEarnedCoins()

	return
}

func getFakeCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin("fake", 10000000),
	}
}

func getFakeEarnedCoins() EarnedCoins {
	return EarnedCoins{
		EarnedCoin{
			sdk.NewInt64Coin("fake", 10000000),
			uint64(1), // CommunityID
		},
	}
}

func getFakeKeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
