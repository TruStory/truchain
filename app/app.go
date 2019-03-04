package app

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	clientParams "github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/truapi"
	"github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/users"
	"github.com/TruStory/truchain/x/vote"
	"github.com/TruStory/truchain/x/voting"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
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
	keyAccount              *sdk.KVStoreKey
	keyBacking              *sdk.KVStoreKey
	keyCategory             *sdk.KVStoreKey
	keyChallenge            *sdk.KVStoreKey
	keyExpiration           *sdk.KVStoreKey
	keyFee                  *sdk.KVStoreKey
	keyIBC                  *sdk.KVStoreKey
	keyMain                 *sdk.KVStoreKey
	keyStake                *sdk.KVStoreKey
	keyStory                *sdk.KVStoreKey
	keyPendingStoryQueue    *sdk.KVStoreKey
	keyTruBank              *sdk.KVStoreKey
	keyChallengedStoryQueue *sdk.KVStoreKey
	keyExpiredStoryQueue    *sdk.KVStoreKey
	keyVote                 *sdk.KVStoreKey
	keyVoting               *sdk.KVStoreKey
	keyParams               *sdk.KVStoreKey
	tkeyParams              *sdk.TransientStoreKey

	// manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	paramsKeeper        params.Keeper

	// access truchain multistore
	backingKeeper      backing.WriteKeeper
	categoryKeeper     category.WriteKeeper
	challengeKeeper    challenge.WriteKeeper
	clientParamsKeeper clientParams.Keeper
	expirationKeeper   expiration.Keeper
	storyKeeper        story.WriteKeeper
	stakeKeeper        stake.Keeper
	truBankKeeper      trubank.WriteKeeper
	voteKeeper         vote.WriteKeeper
	votingKeeper       voting.WriteKeeper

	// state to run api
	blockCtx     *sdk.Context
	blockHeader  abci.Header
	api          *truapi.TruAPI
	apiStarted   bool
	registrarKey secp256k1.PrivKeySecp256k1
}

// NewTruChain returns a reference to a new TruChain. Internally,
// a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewTruChain(logger log.Logger, db dbm.DB, options ...func(*bam.BaseApp)) *TruChain {
	// create and register app-level codec for TXs and accounts
	codec := MakeCodec()

	loadEnvVars()

	// create your application type
	var app = &TruChain{
		BaseApp: bam.NewBaseApp(types.AppName, logger, db, auth.DefaultTxDecoder(codec), options...),
		codec:   codec,

		keyParams:  sdk.NewKVStoreKey("params"),
		tkeyParams: sdk.NewTransientStoreKey("transient_params"),

		keyMain:                 sdk.NewKVStoreKey("main"),
		keyAccount:              sdk.NewKVStoreKey("acc"),
		keyIBC:                  sdk.NewKVStoreKey("ibc"),
		keyStory:                sdk.NewKVStoreKey(story.StoreKey),
		keyPendingStoryQueue:    sdk.NewKVStoreKey(story.PendingQueueStoreKey),
		keyCategory:             sdk.NewKVStoreKey(category.StoreKey),
		keyBacking:              sdk.NewKVStoreKey(backing.StoreKey),
		keyChallenge:            sdk.NewKVStoreKey(challenge.StoreKey),
		keyExpiration:           sdk.NewKVStoreKey(expiration.StoreKey),
		keyFee:                  sdk.NewKVStoreKey("fee_collection"),
		keyStake:                sdk.NewKVStoreKey(stake.StoreKey),
		keyTruBank:              sdk.NewKVStoreKey(trubank.StoreKey),
		keyChallengedStoryQueue: sdk.NewKVStoreKey(story.ChallengedQueueStoreKey),
		keyExpiredStoryQueue:    sdk.NewKVStoreKey(story.ExpiredQueueStoreKey),
		keyVote:                 sdk.NewKVStoreKey(vote.StoreKey),
		keyVoting:               sdk.NewKVStoreKey(voting.StoreKey),
		api:                     nil,
		apiStarted:              false,
		blockCtx:                nil,
		blockHeader:             abci.Header{},
		registrarKey:            loadRegistrarKey(),
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

	// wire up keepers
	app.categoryKeeper = category.NewKeeper(
		app.keyCategory,
		codec,
	)

	app.storyKeeper = story.NewKeeper(
		app.keyStory,
		app.keyPendingStoryQueue,
		app.keyExpiredStoryQueue,
		app.keyChallengedStoryQueue,
		app.categoryKeeper,
		app.paramsKeeper.Subspace(story.StoreKey),
		app.codec,
	)

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
		app.stakeKeeper,
		app.storyKeeper,
		app.bankKeeper,
		app.categoryKeeper,
		codec,
	)

	app.voteKeeper = vote.NewKeeper(
		app.keyVote,
		app.keyChallengedStoryQueue,
		app.stakeKeeper,
		app.accountKeeper,
		app.backingKeeper,
		app.challengeKeeper,
		app.storyKeeper,
		app.bankKeeper,
		app.paramsKeeper.Subspace(vote.StoreKey),
		codec,
	)

	app.challengeKeeper = challenge.NewKeeper(
		app.keyChallenge,
		app.stakeKeeper,
		app.backingKeeper,
		app.bankKeeper,
		app.storyKeeper,
		app.paramsKeeper.Subspace(challenge.StoreKey),
		codec,
	)

	app.expirationKeeper = expiration.NewKeeper(
		app.keyExpiration,
		app.keyExpiredStoryQueue,
		app.stakeKeeper,
		app.storyKeeper,
		app.backingKeeper,
		app.challengeKeeper,
		app.paramsKeeper.Subspace(expiration.StoreKey),
		codec,
	)

	app.votingKeeper = voting.NewKeeper(
		app.keyVoting,
		app.keyChallengedStoryQueue,
		app.accountKeeper,
		app.backingKeeper,
		app.challengeKeeper,
		app.stakeKeeper,
		app.storyKeeper,
		app.voteKeeper,
		app.bankKeeper,
		app.truBankKeeper,
		app.paramsKeeper.Subspace(voting.StoreKey),
		codec,
	)

	app.clientParamsKeeper = clientParams.NewKeeper(
		app.backingKeeper,
		app.challengeKeeper,
		app.expirationKeeper,
		app.stakeKeeper,
		app.storyKeeper,
		app.votingKeeper,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	// TODO [shanev]: see https://github.com/TruStory/truchain/issues/364
	// Add this back after fixing issues with signature verification
	// app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute("story", story.NewHandler(app.storyKeeper)).
		AddRoute("category", category.NewHandler(app.categoryKeeper)).
		AddRoute("backing", backing.NewHandler(app.backingKeeper)).
		AddRoute("challenge", challenge.NewHandler(app.challengeKeeper)).
		AddRoute("vote", vote.NewHandler(app.voteKeeper)).
		AddRoute("users", users.NewHandler(app.accountKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute(story.QueryPath, story.NewQuerier(app.storyKeeper)).
		AddRoute(category.QueryPath, category.NewQuerier(app.categoryKeeper)).
		AddRoute(users.QueryPath, users.NewQuerier(codec, app.accountKeeper)).
		AddRoute(backing.QueryPath, backing.NewQuerier(app.backingKeeper)).
		AddRoute(challenge.QueryPath, challenge.NewQuerier(app.challengeKeeper)).
		AddRoute(vote.QueryPath, vote.NewQuerier(app.voteKeeper)).
		AddRoute(clientParams.QueryPath, clientParams.NewQuerier(app.clientParamsKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)

	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// mount the multistore and load the latest state
	app.MountStores(
		app.keyAccount,
		app.keyParams,
		app.keyBacking,
		app.keyCategory,
		app.keyChallenge,
		app.keyExpiration,
		app.keyFee,
		app.keyIBC,
		app.keyMain,
		app.keyStory,
		app.keyPendingStoryQueue,
		app.keyExpiredStoryQueue,
		app.keyChallengedStoryQueue,
		app.keyTruBank,
		app.keyVote,
		app.keyVoting,
	)

	app.MountStoresTransient(app.tkeyParams)

	err := app.LoadLatestVersion(app.keyMain)

	if err != nil {
		cmn.Exit(err.Error())
	}

	// build HTTP api
	app.api = app.makeAPI()

	return app
}

func loadEnvVars() {
	rootdir := viper.GetString(cli.HomeFlag)
	if rootdir == "" {
		rootdir = DefaultNodeHome
	}

	envPath := filepath.Join(rootdir, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		panic("Error loading .env file")
	}
}

// MakeCodec creates a new codec codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)

	// register msg types
	story.RegisterAmino(cdc)
	backing.RegisterAmino(cdc)
	category.RegisterAmino(cdc)
	challenge.RegisterAmino(cdc)
	vote.RegisterAmino(cdc)
	users.RegisterAmino(cdc)

	// register other types
	cdc.RegisterConcrete(&types.AppAccount{}, "types/AppAccount", nil)

	codec.RegisterCrypto(cdc)

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *TruChain) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.blockCtx = &ctx
	app.blockHeader = req.Header

	if !(app.apiStarted) {
		go app.startAPI()
		app.apiStarted = true
	}

	if app.apiStarted == false && ctx.BlockHeight() > int64(1) {
		panic("API server not started.")
	}

	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *TruChain) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	app.storyKeeper.EndBlock(ctx)
	app.expirationKeeper.EndBlock(ctx)
	app.votingKeeper.EndBlock(ctx)

	return abci.ResponseEndBlock{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *TruChain) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*auth.BaseAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenesisState{
		Accounts:       accounts,
		AuthData:       auth.DefaultGenesisState(),
		BankData:       bank.DefaultGenesisState(),
		Categories:     category.DefaultCategories(),
		ExpirationData: expiration.DefaultGenesisState(),
		StoryData:      story.DefaultGenesisState(),
		VotingData:     voting.DefaultGenesisState(),
	}

	appState, err = codec.MarshalJSONIndent(app.codec, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

func loadRegistrarKey() secp256k1.PrivKeySecp256k1 {
	rootdir := viper.GetString(cli.HomeFlag)
	if rootdir == "" {
		rootdir = DefaultNodeHome
	}

	keypath := filepath.Join(rootdir, "registrar.key")
	fileBytes, err := ioutil.ReadFile(keypath)

	if err != nil {
		panic(err)
	}

	keyBytes, err := hex.DecodeString(string(fileBytes))

	if err != nil {
		panic(err)
	}

	if len(keyBytes) != 32 {
		panic("Invalid registrar key: " + string(fileBytes))
	}

	key := secp256k1.PrivKeySecp256k1{}

	copy(key[:], keyBytes)

	return key
}
