package app

import (
	"os"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	truauth "github.com/TruStory/truchain/x/auth"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/expiration"
	clientParams "github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	trubank "github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/users"
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
		argument.AppModuleBasic{},
		category.AppModuleBasic{},
		backing.AppModuleBasic{},
		challenge.AppModuleBasic{},
		expiration.AppModuleBasic{},
		stake.AppModuleBasic{},
		story.AppModuleBasic{},
		trubank.AppModuleBasic{},
		truauth.AppModuleBasic{},
		community.AppModuleBasic{},
		claim.AppModuleBasic{},
		mint.AppModuleBasic{},
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
	keyArgument   *sdk.KVStoreKey
	keyBacking    *sdk.KVStoreKey
	keyCategory   *sdk.KVStoreKey
	keyChallenge  *sdk.KVStoreKey
	keyExpiration *sdk.KVStoreKey
	keyStake      *sdk.KVStoreKey
	keyStory      *sdk.KVStoreKey
	keyStoryQueue *sdk.KVStoreKey
	keyTruBank    *sdk.KVStoreKey
	keyMint       *sdk.KVStoreKey

	// keys to access trustory V2 state
	keyTruAuth   *sdk.KVStoreKey
	keyCommunity *sdk.KVStoreKey
	keyClaim     *sdk.KVStoreKey

	// manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	ibcMapper           ibc.Mapper
	distrKeeper         distr.Keeper
	paramsKeeper        params.Keeper

	// access truchain multistore
	argumentKeeper     argument.Keeper
	backingKeeper      backing.Keeper
	categoryKeeper     category.Keeper
	challengeKeeper    challenge.Keeper
	clientParamsKeeper clientParams.Keeper
	expirationKeeper   expiration.Keeper
	storyKeeper        story.Keeper
	stakeKeeper        stake.Keeper
	truBankKeeper      trubank.Keeper
	mintKeeper         mint.Keeper

	// access truchain V2 multistore
	truAuthKeeper   truauth.Keeper
	communityKeeper community.Keeper
	claimKeeper     claim.Keeper

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
		keyArgument:      sdk.NewKVStoreKey(argument.StoreKey),
		keyStory:         sdk.NewKVStoreKey(story.StoreKey),
		keyStoryQueue:    sdk.NewKVStoreKey(story.QueueStoreKey),
		keyCategory:      sdk.NewKVStoreKey(category.StoreKey),
		keyBacking:       sdk.NewKVStoreKey(backing.StoreKey),
		keyChallenge:     sdk.NewKVStoreKey(challenge.StoreKey),
		keyExpiration:    sdk.NewKVStoreKey(expiration.StoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyStake:         sdk.NewKVStoreKey(stake.StoreKey),
		keyTruBank:       sdk.NewKVStoreKey(trubank.StoreKey),
		keyTruAuth:       sdk.NewKVStoreKey(truauth.StoreKey),
		keyCommunity:     sdk.NewKVStoreKey(community.StoreKey),
		keyClaim:         sdk.NewKVStoreKey(claim.StoreKey),
		keyMint:          sdk.NewKVStoreKey(mint.StoreKey),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.codec, app.keyParams, app.tkeyParams, params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	storySubspace := app.paramsKeeper.Subspace(story.StoreKey)
	argumentSubspace := app.paramsKeeper.Subspace(argument.StoreKey)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)

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

	// wire up keepers
	app.categoryKeeper = category.NewKeeper(app.keyCategory, codec)
	app.storyKeeper = story.NewKeeper(app.keyStory, app.keyStoryQueue,
		app.categoryKeeper, storySubspace, app.codec)
	app.argumentKeeper = argument.NewKeeper(app.keyArgument, app.storyKeeper,
		argumentSubspace, app.codec)

	app.truBankKeeper = trubank.NewKeeper(
		app.keyTruBank,
		app.bankKeeper,
		app.categoryKeeper,
		app.codec)

	app.stakeKeeper = stake.NewKeeper(
		app.storyKeeper,
		app.truBankKeeper,
		app.paramsKeeper.Subspace(stake.StoreKey),
	)

	app.backingKeeper = backing.NewKeeper(
		app.keyBacking,
		app.argumentKeeper,
		app.stakeKeeper,
		app.storyKeeper,
		app.bankKeeper,
		app.truBankKeeper,
		app.categoryKeeper,
		codec,
	)

	app.challengeKeeper = challenge.NewKeeper(
		app.keyChallenge,
		app.argumentKeeper,
		app.stakeKeeper,
		app.backingKeeper,
		app.truBankKeeper,
		app.bankKeeper,
		app.storyKeeper,
		app.categoryKeeper,
		app.paramsKeeper.Subspace(challenge.StoreKey),
		codec,
	)

	app.expirationKeeper = expiration.NewKeeper(
		app.keyExpiration,
		app.keyStoryQueue,
		app.stakeKeeper,
		app.storyKeeper,
		app.backingKeeper,
		app.challengeKeeper,
		app.paramsKeeper.Subspace(expiration.StoreKey),
		codec,
	)

	app.clientParamsKeeper = clientParams.NewKeeper(
		app.argumentKeeper,
		app.backingKeeper,
		app.challengeKeeper,
		app.expirationKeeper,
		app.stakeKeeper,
		app.storyKeeper,
	)

	app.truAuthKeeper = truauth.NewKeeper(
		app.keyTruAuth,
		app.paramsKeeper.Subspace(truauth.StoreKey),
		codec,
		nil,
		app.accountKeeper,
	)

	app.communityKeeper = community.NewKeeper(
		app.keyCommunity,
		app.paramsKeeper.Subspace(community.StoreKey),
		codec,
	)

	app.claimKeeper = claim.NewKeeper(
		app.keyClaim,
		app.paramsKeeper.Subspace(claim.StoreKey),
		codec,
		app.truAuthKeeper,
		app.communityKeeper,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	// TODO [shanev]: see https://github.com/TruStory/truchain/issues/364
	// Add this back after fixing issues with signature verification
	// app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute("story", story.NewHandler(app.storyKeeper)).
		AddRoute("category", category.NewHandler(app.categoryKeeper)).
		AddRoute("backing", backing.NewHandler(app.backingKeeper)).
		AddRoute("challenge", challenge.NewHandler(app.challengeKeeper)).
		AddRoute("users", users.NewHandler(app.accountKeeper, app.categoryKeeper)).
		AddRoute("trubank", trubank.NewHandler(app.truBankKeeper)).
		AddRoute(truauth.RouterKey, truauth.NewHandler(app.truAuthKeeper)).
		AddRoute(community.RouterKey, community.NewHandler(app.communityKeeper)).
		AddRoute(claim.RouterKey, claim.NewHandler(app.claimKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(argument.QueryPath, argument.NewQuerier(app.argumentKeeper)).
		AddRoute(story.QueryPath, story.NewQuerier(app.storyKeeper)).
		AddRoute(category.QueryPath, category.NewQuerier(app.categoryKeeper)).
		AddRoute(users.QueryPath, users.NewQuerier(codec, app.accountKeeper)).
		AddRoute(backing.QueryPath, backing.NewQuerier(app.backingKeeper)).
		AddRoute(challenge.QueryPath, challenge.NewQuerier(app.challengeKeeper)).
		AddRoute(clientParams.QueryPath, clientParams.NewQuerier(app.clientParamsKeeper)).
		AddRoute(trubank.QueryPath, trubank.NewQuerier(app.truBankKeeper)).
		AddRoute(truauth.QuerierRoute, truauth.NewQuerier(app.truAuthKeeper)).
		AddRoute(community.QuerierRoute, community.NewQuerier(app.communityKeeper)).
		AddRoute(claim.QuerierRoute, claim.NewQuerier(app.claimKeeper))

	app.mm = sdk.NewModuleManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper, app.feeCollectionKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper),
		staking.NewAppModule(app.stakingKeeper, app.feeCollectionKeeper, app.distrKeeper, app.accountKeeper),
		story.NewAppModule(app.storyKeeper),
		category.NewAppModule(app.categoryKeeper),
		argument.NewAppModule(app.argumentKeeper),
		stake.NewAppModule(app.stakeKeeper),
		backing.NewAppModule(app.backingKeeper),
		challenge.NewAppModule(app.challengeKeeper),
		expiration.NewAppModule(app.expirationKeeper),
		trubank.NewAppModule(app.truBankKeeper),
		truauth.NewAppModule(app.truAuthKeeper),
		community.NewAppModule(app.communityKeeper),
		claim.NewAppModule(app.claimKeeper),
		mint.NewAppModule(app.mintKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName, expiration.ModuleName)

	// genutils must occur after staking so that pools are properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName,
		genutil.ModuleName, category.ModuleName, story.ModuleName,
		argument.ModuleName, stake.ModuleName, backing.ModuleName, challenge.ModuleName,
		expiration.ModuleName, trubank.ModuleName, truauth.ModuleName, community.ModuleName,
		claim.ModuleName, mint.ModuleName)

	app.SetInitChainer(app.InitChainer)

	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// mount the multistore and load the latest state
	app.MountStores(
		app.keyAccount,
		app.keyParams,
		app.keyStaking,
		app.keyDistr,
		app.keyBacking,
		app.keyCategory,
		app.keyChallenge,
		app.keyExpiration,
		app.keyFeeCollection,
		app.keyIBC,
		app.keyMain,
		app.keyStory,
		app.keyStoryQueue,
		app.keyTruBank,
		app.keyArgument,
		app.tkeyParams,
		app.tkeyStaking,
		app.tkeyDistr,
		app.keyTruAuth,
		app.keyCommunity,
		app.keyClaim,
		app.keyMint,
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

	users.RegisterCodec(cdc)
	// register other types
	cdc.RegisterConcrete(&types.AppAccount{}, "types/AppAccount", nil)

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
func (app *TruChain) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := app.expirationKeeper.EndBlock(ctx)

	return abci.ResponseEndBlock{Tags: tags}
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
