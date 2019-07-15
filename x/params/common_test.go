package params

import (
	"net/url"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/staking"
	"github.com/tendermint/tendermint/crypto"

	trubank "github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/slashing"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkParams "github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	communityKey := sdk.NewKVStoreKey(community.ModuleName)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	accountKey := sdk.NewKVStoreKey(account.StoreKey)
	claimKey := sdk.NewKVStoreKey(claim.ModuleName)
	bankKey := sdk.NewKVStoreKey(trubank.ModuleName)
	slashKey := sdk.NewKVStoreKey(slashing.ModuleName)
	stakingKey := sdk.NewKVStoreKey(staking.ModuleName)
	sdkParamsKey := sdk.NewKVStoreKey(sdkParams.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(sdkParams.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(slashKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(stakingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(sdkParamsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(bankKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(communityKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(claimKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	account.RegisterCodec(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	RegisterCodec(codec)

	sdkParamsKeeper := sdkParams.NewKeeper(codec, sdkParamsKey, transientParamsKey, sdkParams.DefaultCodespace)

	authKeeper := auth.NewAccountKeeper(
		codec,
		authKey,
		sdkParamsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount)

	bankKeeper := bank.NewBaseKeeper(
		authKeeper,
		sdkParamsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace)

	trubankKeeper := trubank.NewKeeper(
		codec,
		bankKey,
		bankKeeper,
		sdkParamsKeeper.Subspace(trubank.DefaultParamspace),
		trubank.DefaultCodespace)
	trubank.InitGenesis(ctx, trubankKeeper, trubank.DefaultGenesisState())

	communityKeeper := community.NewKeeper(
		communityKey,
		sdkParamsKeeper.Subspace(community.ModuleName),
		codec)
	communityID := "furry"
	community.InitGenesis(ctx, communityKeeper, community.DefaultGenesisState())
	_, err := communityKeeper.NewCommunity(ctx, communityID, "Furries", "")
	if err != nil {
		panic(err)
	}

	accountKeeper := account.NewKeeper(
		accountKey,
		sdkParamsKeeper.Subspace(account.DefaultParamspace),
		codec,
		trubankKeeper,
		authKeeper,
	)
	account.InitGenesis(ctx, accountKeeper, account.DefaultGenesisState())

	_, publicKey, creator, coins := getFakeAppAccountParams()
	_, err = accountKeeper.CreateAppAccount(ctx, creator, coins, publicKey)
	if err != nil {
		panic(err)
	}

	claimKeeper := claim.NewKeeper(
		claimKey,
		sdkParamsKeeper.Subspace(claim.DefaultParamspace),
		codec,
		accountKeeper,
		communityKeeper,
	)
	claim.InitGenesis(ctx, claimKeeper, claim.DefaultGenesisState())

	claim1, err := claimKeeper.SubmitClaim(ctx, "blockchains will allow communities to self governance and manage their own value", communityID, creator, url.URL{})
	if err != nil {
		panic(err)
	}

	stakingKeeper := staking.NewKeeper(
		codec,
		stakingKey,
		accountKeeper,
		trubankKeeper,
		claimKeeper,
		sdkParamsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	staking.InitGenesis(ctx, stakingKeeper, staking.DefaultGenesisState())

	_, err = stakingKeeper.SubmitArgument(ctx, "argument", "summary", creator, claim1.ID, staking.StakeBacking)
	if err != nil {
		panic(err)
	}

	slashKeeper := slashing.NewKeeper(slashKey, sdkParamsKeeper.Subspace(slashing.ModuleName), codec, trubankKeeper, stakingKeeper, accountKeeper, claimKeeper)
	// create fake admins
	_, pubKey, addr1, coins := getFakeAppAccountParams()
	accountKeeper.CreateAppAccount(ctx, addr1, coins, pubKey)
	_, pubKey, addr2, coins := getFakeAppAccountParams()
	accountKeeper.CreateAppAccount(ctx, addr2, coins, pubKey)
	genesis := slashing.DefaultGenesisState()
	genesis.Params.SlashAdmins = append(genesis.Params.SlashAdmins, addr1, addr2)
	slashing.InitGenesis(ctx, slashKeeper, genesis)

	paramsKeeper := NewKeeper(accountKeeper, communityKeeper, claimKeeper, trubankKeeper, stakingKeeper, slashKeeper)

	return ctx, paramsKeeper
}

func getFakeAppAccountParams() (privateKey crypto.PrivKey, publicKey crypto.PubKey, address sdk.AccAddress, coins sdk.Coins) {
	privateKey, publicKey, address = getFakeKeyPubAddr()
	coins = getFakeCoins()

	return
}

func getFakeCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin(app.StakeDenom, 100000000000),
	}
}

func getFakeKeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
