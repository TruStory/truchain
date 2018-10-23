package app

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
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
	keyMain     *sdk.KVStoreKey
	keyAccount  *sdk.KVStoreKey
	keyFee      *sdk.KVStoreKey
	keyIBC      *sdk.KVStoreKey
	keyStory    *sdk.KVStoreKey
	keyCategory *sdk.KVStoreKey
	keyBacking  *sdk.KVStoreKey

	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper

	// access truchain database
	readStoryKeeper story.ReadKeeper
	storyKeeper     story.ReadWriteKeeper
	categoryKeeper  category.ReadWriteKeeper
	backingKeeper   backing.ReadWriteKeeper

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

	// create your application type
	var app = &TruChain{
		codec:        codec,
		BaseApp:      bam.NewBaseApp(AppName, logger, db, auth.DefaultTxDecoder(codec), options...),
		keyMain:      sdk.NewKVStoreKey("main"),
		keyAccount:   sdk.NewKVStoreKey("acc"),
		keyIBC:       sdk.NewKVStoreKey("ibc"),
		keyFee:       sdk.NewKVStoreKey("collectedFees"),
		keyStory:     sdk.NewKVStoreKey("stories"),
		keyCategory:  sdk.NewKVStoreKey("categories"),
		keyBacking:   sdk.NewKVStoreKey("backings"),
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

	// wire up trustory keepers
	app.categoryKeeper = category.NewKeeper(app.keyCategory, app.keyStory, codec)
	app.readStoryKeeper = story.NewKeeper(app.keyStory, app.keyCategory, app.categoryKeeper, app.codec)
	app.storyKeeper = story.NewKeeper(app.keyStory, app.keyCategory, app.categoryKeeper, app.codec)
	app.backingKeeper = backing.NewKeeper(app.keyBacking, app.storyKeeper, app.coinKeeper, app.categoryKeeper, codec)

	// register message routes for modifying state
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("story", story.NewHandler(app.storyKeeper)).
		AddRoute("category", category.NewHandler(app.categoryKeeper)).
		AddRoute("backing", backing.NewHandler(app.backingKeeper)).
		AddRoute(registration.RegisterKeyMsg{}.Type(),
			registration.NewHandler(app.accountMapper))

	// register query routes for reading state
	app.QueryRouter().
		AddRoute("story", story.NewQuerier(app.readStoryKeeper))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStory, app.keyBacking, app.keyFee, app.keyCategory)
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
	registration.RegisterAmino(cdc)

	// register other custom types
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
	return app.backingKeeper.NewResponseEndBlock(ctx)
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *TruChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.codec.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
		panic(err)
	}

	for i, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
			panic(err)
		}

		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)

		if i == 1 { // TODO: more robust way of identifying registrar account [notduncansmith]
			acc.BaseAccount.SetPubKey(app.registrarKey.PubKey())
		}

		app.accountMapper.SetAccount(ctx, acc)
	}

	return abci.ResponseInitChain{}
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
