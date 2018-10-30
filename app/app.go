package app

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	params "github.com/TruStory/truchain/parameters"
	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/registration"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/truapi"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// TruChain implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type TruChain struct {
	*bam.BaseApp
	codec *codec.Codec

	// keys to access the multistore
	keyMain      *sdk.KVStoreKey
	keyAccount   *sdk.KVStoreKey
	keyIBC       *sdk.KVStoreKey
	keyStory     *sdk.KVStoreKey
	keyCategory  *sdk.KVStoreKey
	keyBacking   *sdk.KVStoreKey
	keyChallenge *sdk.KVStoreKey
	keyFee       *sdk.KVStoreKey

	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper

	// access truchain database
	storyKeeper     story.WriteKeeper
	categoryKeeper  category.WriteKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper

	// list of initial categories
	categories map[string]string

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

	// map of initial categories (slug -> title)
	categories := map[string]string{
		"btc":        "Bitcoin",
		"conensus":   "Consensus",
		"eth":        "Ethereum",
		"founder":    "Founders",
		"ico":        "ICOs",
		"regulation": "Regulation",
		"shitcoin":   "Shitcoins",
	}

	// create your application type
	var app = &TruChain{
		categories:   categories,
		codec:        codec,
		BaseApp:      bam.NewBaseApp(params.AppName, logger, db, auth.DefaultTxDecoder(codec), options...),
		keyMain:      sdk.NewKVStoreKey("main"),
		keyAccount:   sdk.NewKVStoreKey("acc"),
		keyIBC:       sdk.NewKVStoreKey("ibc"),
		keyStory:     sdk.NewKVStoreKey("stories"),
		keyCategory:  sdk.NewKVStoreKey("categories"),
		keyBacking:   sdk.NewKVStoreKey("backings"),
		keyChallenge: sdk.NewKVStoreKey("challenges"),
		keyFee:       sdk.NewKVStoreKey("collectedFees"),
		api:          nil,
		apiStarted:   false,
		blockCtx:     nil,
		blockHeader:  abci.Header{},
		registrarKey: loadRegistrarKey(),
	}

	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(
		codec,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)
	app.coinKeeper = bank.NewBaseKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.codec, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.codec, app.keyFee)

	// wire up keepers
	app.categoryKeeper = category.NewKeeper(app.keyCategory, codec)
	app.storyKeeper = story.NewKeeper(
		app.keyStory, app.keyCategory, app.keyChallenge,
		app.categoryKeeper, app.codec)
	app.backingKeeper = backing.NewKeeper(
		app.keyBacking, app.storyKeeper, app.coinKeeper,
		app.categoryKeeper, codec)
	app.challengeKeeper = challenge.NewKeeper(
		app.keyChallenge, app.storyKeeper, app.coinKeeper, codec)

	// register message routes for modifying state
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("story", story.NewHandler(app.storyKeeper)).
		AddRoute("category", category.NewHandler(app.categoryKeeper)).
		AddRoute("backing", backing.NewHandler(app.backingKeeper)).
		AddRoute("challenge", challenge.NewHandler(app.challengeKeeper)).
		AddRoute(registration.RegisterKeyMsg{}.Type(),
			registration.NewHandler(app.accountMapper))

	// register query routes for reading state
	app.QueryRouter().
		AddRoute(story.QueryPath, story.NewQuerier(app.storyKeeper))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStory, app.keyBacking, app.keyFee, app.keyCategory, app.keyChallenge)

	err := app.LoadLatestVersion(app.keyMain)

	if err != nil {
		cmn.Exit(err.Error())
	}

	// build HTTP api
	app.api = app.makeAPI()

	return app
}

// MakeCodec creates a new codec codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)

	// register msg types
	story.RegisterAmino(cdc)
	backing.RegisterAmino(cdc)
	category.RegisterAmino(cdc)
	challenge.RegisterAmino(cdc)
	registration.RegisterAmino(cdc)

	// register other types
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&types.AppAccount{}, "truchain/Account", nil)
	cdc.RegisterConcrete(&auth.StdTx{}, "cosmos-sdk/StdTx", nil)

	cdc.Seal()

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

	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *TruChain) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	app.backingKeeper.NewResponseEndBlock(ctx)
	app.challengeKeeper.NewResponseEndBlock(ctx)

	return abci.ResponseEndBlock{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *TruChain) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*types.GenesisAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountMapper.IterateAccounts(ctx, appendAccountsFn)

	genState := types.GenesisState{Accounts: accounts}
	appState, err = codec.MarshalJSONIndent(app.codec, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

func loadRegistrarKey() secp256k1.PrivKeySecp256k1 {
	fileBytes, err := ioutil.ReadFile("registrar.key")

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
