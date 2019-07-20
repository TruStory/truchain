package app

import (
	"os"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/account"
	trubank "github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	truslashing "github.com/TruStory/truchain/x/slashing"
	trustaking "github.com/TruStory/truchain/x/staking"
)

const (
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.trucli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.truchaind")
	// The ModuleBasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics sdk.ModuleBasicManager
)

func init() {
	ModuleBasics = sdk.NewModuleBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		params.AppModuleBasic{},
		community.AppModuleBasic{},
		claim.AppModuleBasic{},
		mint.AppModuleBasic{},
		account.AppModuleBasic{},
		trubank.AppModuleBasic{},
		trustaking.AppModuleBasic{},
		truslashing.AppModuleBasic{},
	)
}

// TruChain implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type TruChain struct {
	*bam.BaseApp
	codec *codec.Codec

	// keys to access the multistore
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyIBC           *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey

	// keys to access trustory state
	keyTruBank     *sdk.KVStoreKey
	keyMint        *sdk.KVStoreKey
	keyAppAccount  *sdk.KVStoreKey
	keyCommunity   *sdk.KVStoreKey
	keyClaim       *sdk.KVStoreKey
	keyTruStaking  *sdk.KVStoreKey
	keyTruSlashing *sdk.KVStoreKey

	// manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	ibcMapper           ibc.Mapper
	distrKeeper         distr.Keeper
	paramsKeeper        params.Keeper

	// access truchain multistore
	mintKeeper        mint.Keeper
	appAccountKeeper  account.Keeper
	communityKeeper   community.Keeper
	claimKeeper       claim.Keeper
	truBankKeeper     trubank.Keeper
	truStakingKeeper  trustaking.Keeper
	truSlashingKeeper truslashing.Keeper

	// the module manager
	mm *sdk.ModuleManager
}

// NewTruChain returns a reference to a new TruChain. Internally,
// a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewTruChain(logger log.Logger, db dbm.DB, loadLatest bool, options ...func(*bam.BaseApp)) *TruChain {
	// create and register app-level codec for TXs and accounts
	codec := MakeCodec()

	bApp := bam.NewBaseApp(types.AppName, logger, db, auth.DefaultTxDecoder(codec), options...)
	bApp.SetAppVersion(version.Version)

	// create your application type
	var app = &TruChain{
		BaseApp: bApp,
		codec:   codec,

		keyParams:  sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams: sdk.NewTransientStoreKey(params.TStoreKey),

		keyMain:          sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyIBC:           sdk.NewKVStoreKey("ibc"),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:         sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:        sdk.NewTransientStoreKey(distr.TStoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyCommunity:     sdk.NewKVStoreKey(community.StoreKey),
		keyClaim:         sdk.NewKVStoreKey(claim.StoreKey),
		keyMint:          sdk.NewKVStoreKey(mint.StoreKey),
		keyAppAccount:    sdk.NewKVStoreKey(account.StoreKey),
		keyTruStaking:    sdk.NewKVStoreKey(trustaking.StoreKey),
		keyTruBank:       sdk.NewKVStoreKey(trubank.StoreKey),
		keyTruSlashing:   sdk.NewKVStoreKey(truslashing.StoreKey),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.codec, app.keyParams, app.tkeyParams, params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	appAccountSubspace := app.paramsKeeper.Subspace(account.DefaultParamspace)
	trubank2Subspace := app.paramsKeeper.Subspace(trubank.DefaultParamspace)
	truStakingSubspace := app.paramsKeeper.Subspace(trustaking.DefaultParamspace)
	truSlashingSubspace := app.paramsKeeper.Subspace(truslashing.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.codec, app.keyAccount, authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, bank.DefaultCodespace)
	app.ibcMapper = ibc.NewMapper(app.codec, app.keyIBC, ibc.DefaultCodespace)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.codec, app.keyFeeCollection)

	stakingKeeper := staking.NewKeeper(app.codec, app.keyStaking, app.tkeyStaking, app.bankKeeper,
		stakingSubspace, staking.DefaultCodespace)

	app.mintKeeper = mint.NewKeeper(app.codec, app.keyMint, mintSubspace, &stakingKeeper, app.feeCollectionKeeper)

	app.distrKeeper = distr.NewKeeper(
		app.codec,
		app.keyDistr,
		distrSubspace,
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.stakingKeeper = stakingKeeper

	app.communityKeeper = community.NewKeeper(
		app.keyCommunity,
		app.paramsKeeper.Subspace(community.StoreKey),
		codec,
	)

	app.truBankKeeper = trubank.NewKeeper(codec, app.keyTruBank, app.bankKeeper,
		trubank2Subspace, trubank.DefaultCodespace)

	app.appAccountKeeper = account.NewKeeper(
		app.keyAppAccount,
		appAccountSubspace,
		codec,
		app.truBankKeeper,
		app.accountKeeper,
	)

	app.claimKeeper = claim.NewKeeper(
		app.keyClaim,
		app.paramsKeeper.Subspace(claim.StoreKey),
		codec,
		app.appAccountKeeper,
		app.communityKeeper,
	)

	app.truStakingKeeper = trustaking.NewKeeper(codec, app.keyTruStaking, app.appAccountKeeper,
		app.truBankKeeper, app.claimKeeper, truStakingSubspace, trustaking.DefaultCodespace)

	app.truSlashingKeeper = truslashing.NewKeeper(
		app.keyTruSlashing,
		truSlashingSubspace,
		codec,
		app.truBankKeeper,
		app.truStakingKeeper,
		app.appAccountKeeper,
		app.claimKeeper,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	// TODO [shanev]: see https://github.com/TruStory/truchain/issues/364
	// Add this back after fixing issues with signature verification
	// app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute(claim.RouterKey, claim.NewHandler(app.claimKeeper)).
		AddRoute(account.RouterKey, account.NewHandler(app.appAccountKeeper)).
		AddRoute(trustaking.RouterKey, trustaking.NewHandler(app.truStakingKeeper)).
		AddRoute(truslashing.RouterKey, truslashing.NewHandler(app.truSlashingKeeper)).
		AddRoute(trubank.RouterKey, trubank.NewHandler(app.truBankKeeper)).
		AddRoute(community.RouterKey, community.NewHandler(app.communityKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(community.QuerierRoute, community.NewQuerier(app.communityKeeper)).
		AddRoute(claim.QuerierRoute, claim.NewQuerier(app.claimKeeper)).
		AddRoute(trubank.QuerierRoute, trubank.NewQuerier(app.truBankKeeper)).
		AddRoute(account.QuerierRoute, account.NewQuerier(app.appAccountKeeper)).
		AddRoute(trustaking.QuerierRoute, trustaking.NewQuerier(app.truStakingKeeper)).
		AddRoute(truslashing.QuerierRoute, truslashing.NewQuerier(app.truSlashingKeeper))

	app.mm = sdk.NewModuleManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper, app.feeCollectionKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper),
		staking.NewAppModule(app.stakingKeeper, app.feeCollectionKeeper, app.distrKeeper, app.accountKeeper),
		community.NewAppModule(app.communityKeeper),
		claim.NewAppModule(app.claimKeeper),
		mint.NewAppModule(app.mintKeeper),
		trubank.NewAppModule(app.truBankKeeper),
		account.NewAppModule(app.appAccountKeeper),
		trustaking.NewAppModule(app.truStakingKeeper),
		truslashing.NewAppModule(app.truSlashingKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName, trustaking.ModuleName, truslashing.ModuleName)

	// genutils must occur after staking so that pools are properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, account.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName,
		genutil.ModuleName, community.ModuleName, claim.ModuleName,
		trubank.ModuleName, trustaking.ModuleName, truslashing.ModuleName,
		mint.ModuleName)

	app.SetInitChainer(app.InitChainer)

	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// mount the multistore and load the latest state
	app.MountStores(
		app.keyAccount,
		app.keyParams,
		app.keyStaking,
		app.keyDistr,
		app.keyFeeCollection,
		app.keyIBC,
		app.keyMain,
		app.tkeyParams,
		app.tkeyStaking,
		app.tkeyDistr,
		app.keyCommunity,
		app.keyClaim,
		app.keyMint,
		app.keyAppAccount,
		app.keyTruBank,
		app.keyTruStaking,
		app.keyTruSlashing,
	)

	if loadLatest {
		err := app.LoadLatestVersion(app.keyMain)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

// MakeCodec creates a new codec codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *TruChain) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *TruChain) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *TruChain) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.codec.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads the app at a particular height
func (app *TruChain) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}
