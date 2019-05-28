package app

import (
	"fmt"
	"os"
	"sort"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	clientParams "github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	trubank "github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/users"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/ibc"
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
)

// TruChain implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type TruChain struct {
	*bam.BaseApp
	codec *codec.Codec

	// keys to access the multistore
	keyAccount  *sdk.KVStoreKey
	keyDistr    *sdk.KVStoreKey
	tkeyDistr   *sdk.TransientStoreKey
	keyStaking  *sdk.KVStoreKey
	tkeyStaking *sdk.TransientStoreKey
	keyFee      *sdk.KVStoreKey
	keyIBC      *sdk.KVStoreKey
	keyMain     *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey

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
}

// NewTruChain returns a reference to a new TruChain. Internally,
// a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewTruChain(logger log.Logger, db dbm.DB, loadLatest bool, options ...func(*bam.BaseApp)) *TruChain {
	// create and register app-level codec for TXs and accounts
	codec := MakeCodec()

	// create your application type
	var app = &TruChain{
		BaseApp: bam.NewBaseApp(types.AppName, logger, db, auth.DefaultTxDecoder(codec), options...),
		codec:   codec,

		keyParams:  sdk.NewKVStoreKey("params"),
		tkeyParams: sdk.NewTransientStoreKey("transient_params"),

		keyMain:       sdk.NewKVStoreKey("main"),
		keyAccount:    sdk.NewKVStoreKey("acc"),
		keyIBC:        sdk.NewKVStoreKey("ibc"),
		keyStaking:    sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:   sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:      sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:     sdk.NewTransientStoreKey(distr.TStoreKey),
		keyArgument:   sdk.NewKVStoreKey(argument.StoreKey),
		keyStory:      sdk.NewKVStoreKey(story.StoreKey),
		keyStoryQueue: sdk.NewKVStoreKey(story.QueueStoreKey),
		keyCategory:   sdk.NewKVStoreKey(category.StoreKey),
		keyBacking:    sdk.NewKVStoreKey(backing.StoreKey),
		keyChallenge:  sdk.NewKVStoreKey(challenge.StoreKey),
		keyExpiration: sdk.NewKVStoreKey(expiration.StoreKey),
		keyFee:        sdk.NewKVStoreKey("fee_collection"),
		keyStake:      sdk.NewKVStoreKey(stake.StoreKey),
		keyTruBank:    sdk.NewKVStoreKey(trubank.StoreKey),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.codec, app.keyParams, app.tkeyParams)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.codec,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	app.ibcMapper = ibc.NewMapper(app.codec, app.keyIBC, ibc.DefaultCodespace)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.codec, app.keyFee)

	stakingKeeper := staking.NewKeeper(
		app.codec,
		app.keyStaking, app.tkeyStaking,
		app.bankKeeper, app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	app.distrKeeper = distr.NewKeeper(
		app.codec,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.stakingKeeper = stakingKeeper

	// wire up keepers
	app.categoryKeeper = category.NewKeeper(
		app.keyCategory,
		codec,
	)

	app.storyKeeper = story.NewKeeper(
		app.keyStory,
		app.keyStoryQueue,
		app.categoryKeeper,
		app.paramsKeeper.Subspace(story.StoreKey),
		app.codec,
	)

	app.argumentKeeper = argument.NewKeeper(
		app.keyArgument,
		app.storyKeeper,
		app.paramsKeeper.Subspace(argument.StoreKey),
		app.codec)

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
		AddRoute("trubank", trubank.NewHandler(app.truBankKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute("acc", auth.NewQuerier(app.accountKeeper)).
		AddRoute(argument.QueryPath, argument.NewQuerier(app.argumentKeeper)).
		AddRoute(story.QueryPath, story.NewQuerier(app.storyKeeper)).
		AddRoute(category.QueryPath, category.NewQuerier(app.categoryKeeper)).
		AddRoute(users.QueryPath, users.NewQuerier(codec, app.accountKeeper)).
		AddRoute(backing.QueryPath, backing.NewQuerier(app.backingKeeper)).
		AddRoute(challenge.QueryPath, challenge.NewQuerier(app.challengeKeeper)).
		AddRoute(clientParams.QueryPath, clientParams.NewQuerier(app.clientParamsKeeper)).
		AddRoute(trubank.QueryPath, trubank.NewQuerier(app.truBankKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)

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
		app.keyFee,
		app.keyIBC,
		app.keyMain,
		app.keyStory,
		app.keyStoryQueue,
		app.keyTruBank,
		app.keyArgument,
		app.tkeyParams,
		app.tkeyStaking,
		app.tkeyDistr,
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

	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)

	// register msg types
	story.RegisterAmino(cdc)
	backing.RegisterAmino(cdc)
	category.RegisterAmino(cdc)
	challenge.RegisterAmino(cdc)
	users.RegisterAmino(cdc)
	trubank.RegisterAmino(cdc)

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

// initialize store from a genesis state
func (app *TruChain) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) []abci.ValidatorUpdate {
	genesisState.Sanitize()

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc = app.accountKeeper.NewAccount(ctx, acc) // set account number
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)

	// load the initial staking information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err)
	}

	// initialize module-specific stores
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)

	// trustory-specific modules
	argument.InitGenesis(ctx, app.argumentKeeper, genesisState.ArgumentData)
	category.InitGenesis(ctx, app.categoryKeeper, genesisState.CategoryData)
	challenge.InitGenesis(ctx, app.challengeKeeper, genesisState.ChallengeData)
	backing.InitGenesis(ctx, app.backingKeeper, genesisState.BackingData)
	expiration.InitGenesis(ctx, app.expirationKeeper, genesisState.ExpirationData)
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	story.InitGenesis(ctx, app.storyKeeper, genesisState.StoryData)
	trubank.InitGenesis(ctx, app.truBankKeeper, genesisState.TrubankData)

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.codec.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.codec.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	return validators
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *TruChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisState GenesisState
	err := app.codec.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc.ToAccount())
	}

	// initialize module-specific stores
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)

	// trustory-specific modules
	argument.InitGenesis(ctx, app.argumentKeeper, genesisState.ArgumentData)
	category.InitGenesis(ctx, app.categoryKeeper, genesisState.CategoryData)
	challenge.InitGenesis(ctx, app.challengeKeeper, genesisState.ChallengeData)
	backing.InitGenesis(ctx, app.backingKeeper, genesisState.BackingData)
	expiration.InitGenesis(ctx, app.expirationKeeper, genesisState.ExpirationData)
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	story.InitGenesis(ctx, app.storyKeeper, genesisState.StoryData)
	trubank.InitGenesis(ctx, app.truBankKeeper, genesisState.TrubankData)

	validators := app.initFromGenesisState(ctx, genesisState)

	// sanity check
	if len(req.Validators) > 0 {
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	return abci.ResponseInitChain{
		Validators: validators,
	}
}

// LoadHeight loads the app at a particular height
func (app *TruChain) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}
