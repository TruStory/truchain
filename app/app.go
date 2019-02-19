package app

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	params "github.com/TruStory/truchain/parameters"
	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/game"
	clientParams "github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/truapi"
	"github.com/TruStory/truchain/x/users"
	"github.com/TruStory/truchain/x/vote"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
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
	keyAccount          *sdk.KVStoreKey
	keyBacking          *sdk.KVStoreKey
	keyCategory         *sdk.KVStoreKey
	keyChallenge        *sdk.KVStoreKey
	keyFee              *sdk.KVStoreKey
	keyGame             *sdk.KVStoreKey
	keyIBC              *sdk.KVStoreKey
	keyMain             *sdk.KVStoreKey
	keyStory            *sdk.KVStoreKey
	keyStoryQueue       *sdk.KVStoreKey
	keyVotingStoryQueue *sdk.KVStoreKey
	keyVote             *sdk.KVStoreKey
	keyParams           *sdk.KVStoreKey
	tkeyParams          *sdk.TransientStoreKey

	// manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	paramsKeeper        sdkparams.Keeper

	// access truchain database
	storyKeeper     story.WriteKeeper
	categoryKeeper  category.WriteKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	gameKeeper      game.WriteKeeper
	voteKeeper      vote.WriteKeeper

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
		BaseApp: bam.NewBaseApp(params.AppName, logger, db, auth.DefaultTxDecoder(codec), options...),
		codec:   codec,

		keyParams:  sdk.NewKVStoreKey("params"),
		tkeyParams: sdk.NewTransientStoreKey("transient_params"),

		keyMain:             sdk.NewKVStoreKey("main"),
		keyAccount:          sdk.NewKVStoreKey("acc"),
		keyIBC:              sdk.NewKVStoreKey("ibc"),
		keyStory:            sdk.NewKVStoreKey(story.StoreKey),
		keyStoryQueue:       sdk.NewKVStoreKey(story.QueueStoreKey),
		keyCategory:         sdk.NewKVStoreKey(category.StoreKey),
		keyBacking:          sdk.NewKVStoreKey(backing.StoreKey),
		keyChallenge:        sdk.NewKVStoreKey(challenge.StoreKey),
		keyFee:              sdk.NewKVStoreKey("fee_collection"),
		keyGame:             sdk.NewKVStoreKey(game.StoreKey),
		keyVotingStoryQueue: sdk.NewKVStoreKey(story.VotingQueueStoreKey),
		keyVote:             sdk.NewKVStoreKey(vote.StoreKey),
		api:                 nil,
		apiStarted:          false,
		blockCtx:            nil,
		blockHeader:         abci.Header{},
		registrarKey:        loadRegistrarKey(),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = sdkparams.NewKeeper(app.codec, app.keyParams, app.tkeyParams)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.codec,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	app.coinKeeper = bank.NewBaseKeeper(
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
		app.keyStoryQueue,
		app.keyVotingStoryQueue,
		app.categoryKeeper,
		app.paramsKeeper.Subspace(story.DefaultParamspace),
		app.codec,
	)

	app.backingKeeper = backing.NewKeeper(
		app.keyBacking,
		app.keyStoryQueue,
		app.keyVotingStoryQueue,
		app.storyKeeper,
		app.coinKeeper,
		app.categoryKeeper,
		codec,
	)

	// app.gameKeeper = game.NewKeeper(
	// 	app.keyGame,
	// 	app.keyVotingStoryQueue,
	// 	app.storyKeeper,
	// 	app.backingKeeper,
	// 	app.coinKeeper,
	// 	codec,
	// )

	app.challengeKeeper = challenge.NewKeeper(
		app.keyChallenge,
		app.keyVotingStoryQueue,
		app.backingKeeper,
		app.coinKeeper,
		app.gameKeeper,
		app.storyKeeper,
		codec,
	)

	app.voteKeeper = vote.NewKeeper(
		app.keyVote,
		app.keyVotingStoryQueue,
		app.accountKeeper,
		app.backingKeeper,
		app.challengeKeeper,
		app.storyKeeper,
		app.gameKeeper,
		app.coinKeeper,
		codec,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	// TODO [shanev]: see https://github.com/TruStory/truchain/issues/364
	// Add this back after fixing issues with signature verification
	// app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
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
		AddRoute(game.QueryPath, game.NewQuerier(app.gameKeeper, app.backingKeeper)).
		AddRoute(backing.QueryPath, backing.NewQuerier(app.backingKeeper)).
		AddRoute(challenge.QueryPath, challenge.NewQuerier(app.challengeKeeper)).
		AddRoute(vote.QueryPath, vote.NewQuerier(app.voteKeeper)).
		AddRoute(clientParams.QueryPath, clientParams.NewQuerier())

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
		app.keyFee,
		app.keyGame,
		app.keyVotingStoryQueue,
		app.keyIBC,
		app.keyMain,
		app.keyStory,
		app.keyStoryQueue,
		app.keyVote)

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
	app.storyKeeper.NewResponseEndBlock(ctx)
	// app.backingKeeper.NewResponseEndBlock(ctx)
	// app.challengeKeeper.NewResponseEndBlock(ctx)
	// app.voteKeeper.NewResponseEndBlock(ctx)

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
		Accounts: accounts,
		AuthData: auth.DefaultGenesisState(),
		BankData: bank.DefaultGenesisState(),
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
