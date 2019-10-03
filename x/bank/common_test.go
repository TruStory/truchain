package bank

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/account"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func mockDB() (sdk.Context, Keeper, auth.AccountKeeper) {
	db := dbm.NewMemDB()
	storeKey := sdk.NewKVStoreKey(ModuleName)
	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(supplyKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	// codec registration
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)

	// Keepers
	pk := params.NewKeeper(cdc, paramsKey, transientParamsKey, params.DefaultCodespace)
	accKeeper := auth.NewAccountKeeper(cdc, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		nil,
	)

	maccPerms := map[string][]string{
		account.UserGrowthPoolName: {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, supplyKey, accKeeper, bankKeeper, maccPerms)
	userGrowthAcc := supply.NewEmptyModuleAccount(account.UserGrowthPoolName, supply.Burner, supply.Staking)
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, sdk.NewInt(100000000000)))
	userGrowthAcc.SetCoins(initCoins)
	supplyKeeper.SetModuleAccount(ctx, userGrowthAcc)
	totalSupply := initCoins
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// module keeper
	keeper := NewKeeper(
		cdc,
		storeKey,
		bankKeeper,
		pk.Subspace(DefaultParamspace),
		DefaultCodespace,
		supplyKeeper,
	)

	InitGenesis(ctx, keeper, DefaultGenesisState())
	return ctx, keeper, accKeeper
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
