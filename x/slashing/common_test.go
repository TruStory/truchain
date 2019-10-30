package slashing

import (
	"net/url"

	"github.com/TruStory/truchain/x/distribution"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/staking"
	"github.com/tendermint/tendermint/crypto"

	trubank "github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	communityKey := sdk.NewKVStoreKey(community.ModuleName)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	accountKey := sdk.NewKVStoreKey(account.StoreKey)
	claimKey := sdk.NewKVStoreKey(claim.ModuleName)
	bankKey := sdk.NewKVStoreKey(trubank.ModuleName)
	slashKey := sdk.NewKVStoreKey(ModuleName)
	stakingKey := sdk.NewKVStoreKey(staking.ModuleName)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(slashKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(stakingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(bankKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(communityKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(claimKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(supplyKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	account.RegisterCodec(codec)
	codec.RegisterInterface((*authexported.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)
	RegisterCodec(codec)
	supply.RegisterCodec(codec)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:           nil,
		mint.ModuleName:                 {supply.Minter},
		gov.ModuleName:                  {supply.Burner},
		distribution.UserGrowthPoolName: {supply.Burner, supply.Staking},
		distribution.UserRewardPoolName: {supply.Burner},
		staking.UserStakesPoolName:      {supply.Minter, supply.Burner},
	}

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)

	authKeeper := auth.NewAccountKeeper(
		codec,
		authKey,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount)

	bankKeeper := bank.NewBaseKeeper(
		authKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace, nil)

	userRewardAcc := supply.NewEmptyModuleAccount(staking.UserRewardPoolName, supply.Burner, supply.Staking)
	userGrowthAcc := supply.NewEmptyModuleAccount(account.UserGrowthPoolName, supply.Minter, supply.Burner, supply.Staking)
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, sdk.NewInt(10000000000)))
	err := userRewardAcc.SetCoins(initCoins)
	supplyKeeper := supply.NewKeeper(codec, supplyKey, authKeeper, bankKeeper, maccPerms)
	supplyKeeper.SetModuleAccount(ctx, userGrowthAcc)
	supplyKeeper.SetModuleAccount(ctx, userRewardAcc)
	totalSupply := initCoins
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	trubankKeeper := trubank.NewKeeper(
		codec,
		bankKey,
		bankKeeper,
		paramsKeeper.Subspace(trubank.DefaultParamspace),
		trubank.DefaultCodespace,
		supplyKeeper)

	trubank.InitGenesis(ctx, trubankKeeper, trubank.DefaultGenesisState())

	communityKeeper := community.NewKeeper(
		communityKey,
		paramsKeeper.Subspace(community.ModuleName),
		codec)
	_, _, cAdmin1, _ := getFakeAppAccountParams()
	_, _, cAdmin2, _ := getFakeAppAccountParams()
	cGenesis := community.DefaultGenesisState()
	cGenesis.Params.CommunityAdmins = append(cGenesis.Params.CommunityAdmins, cAdmin1, cAdmin2)
	community.InitGenesis(ctx, communityKeeper, cGenesis)
	communityID := "furry"
	_, err = communityKeeper.NewCommunity(ctx, communityID, "Furries", "", cAdmin1)
	if err != nil {
		panic(err)
	}

	accountKeeper := account.NewKeeper(
		accountKey,
		paramsKeeper.Subspace(account.DefaultParamspace),
		codec,
		trubankKeeper,
		authKeeper,
		supplyKeeper,
	)
	account.InitGenesis(ctx, accountKeeper, account.DefaultGenesisState())

	_, publicKey, creator, coins := getFakeAppAccountParams()
	_, err = accountKeeper.CreateAppAccount(ctx, creator, coins, publicKey)
	if err != nil {
		panic(err)
	}

	claimKeeper := claim.NewKeeper(
		claimKey,
		paramsKeeper.Subspace(claim.DefaultParamspace),
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
		supplyKeeper,
		paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	staking.InitGenesis(ctx, stakingKeeper, staking.DefaultGenesisState())

	_, err = stakingKeeper.SubmitArgument(ctx, "argument", "summary", creator, claim1.ID, staking.StakeBacking)
	if err != nil {
		panic(err)
	}

	slashKeeper := NewKeeper(slashKey, paramsKeeper.Subspace(ModuleName), codec, trubankKeeper, stakingKeeper, accountKeeper, claimKeeper)
	// create fake admins
	_, pubKey, addr1, coins := getFakeAppAccountParams()
	accountKeeper.CreateAppAccount(ctx, addr1, coins, pubKey)
	_, pubKey, addr2, coins := getFakeAppAccountParams()
	accountKeeper.CreateAppAccount(ctx, addr2, coins, pubKey)
	genesis := DefaultGenesisState()
	genesis.Params.SlashAdmins = append(genesis.Params.SlashAdmins, addr1, addr2)
	InitGenesis(ctx, slashKeeper, genesis)

	return ctx, slashKeeper
}

func getFakeAppAccountParams() (privateKey crypto.PrivKey, publicKey crypto.PubKey, address sdk.AccAddress, coins sdk.Coins) {
	privateKey, publicKey, address = getFakeKeyPubAddr()
	coins = getFakeCoins()

	return
}

func getFakeCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin(app.StakeDenom, 300000000000),
	}
}

func getFakeKeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
