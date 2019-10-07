package app

import (
	"os"

	"github.com/cosmos/cosmos-sdk/x/crisis"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/account"
	trubank "github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	trudist "github.com/TruStory/truchain/x/distribution"
	truslashing "github.com/TruStory/truchain/x/slashing"
	trustaking "github.com/TruStory/truchain/x/staking"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.truchaincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.truchaind")

	// The ModuleBasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		// trustory modules
		community.AppModuleBasic{},
		claim.AppModuleBasic{},
		account.AppModuleBasic{},
		trubank.AppModuleBasic{},
		trustaking.AppModuleBasic{},
		truslashing.AppModuleBasic{},
		trudist.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		// trustory module accounts
		trudist.UserGrowthPoolName:    {supply.Minter, supply.Burner},
		trudist.StakeholderPoolName:   nil,
		trudist.UserRewardPoolName:    {supply.Minter, supply.Burner},
		trustaking.UserStakesPoolName: {supply.Minter, supply.Burner},
	}
)

// TruChain implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type TruChain struct {
	*bam.BaseApp
	codec *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// cosmos keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper

	// trustory keepers
	appAccountKeeper      account.Keeper
	communityKeeper       community.Keeper
	claimKeeper           claim.Keeper
	truBankKeeper         trubank.Keeper
	truStakingKeeper      trustaking.Keeper
	truSlashingKeeper     truslashing.Keeper
	truDistributionKeeper trudist.Keeper

	// the module manager
	mm *module.Manager
}

// NewTruChain returns a reference to a new TruChain. Internally,
// a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewTruChain(logger log.Logger, db dbm.DB, loadLatest bool,
	invCheckPeriod uint, options ...func(*bam.BaseApp)) *TruChain {
	// create and register app-level codec for TXs and accounts
	codec := MakeCodec()

	bApp := bam.NewBaseApp(types.AppName, logger, db, auth.DefaultTxDecoder(codec), options...)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey,
		community.StoreKey, claim.StoreKey, account.StoreKey, trustaking.StoreKey,
		trubank.StoreKey, truslashing.StoreKey, trudist.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// create your application type
	var app = &TruChain{
		BaseApp: bApp,
		codec:   codec,
		keys:    keys,
		tkeys:   tkeys,
	}

	// init params keeper and cosmos subspaces
	app.paramsKeeper = params.NewKeeper(app.codec, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	// trustory subspaces
	appAccountSubspace := app.paramsKeeper.Subspace(account.DefaultParamspace)
	trubank2Subspace := app.paramsKeeper.Subspace(trubank.DefaultParamspace)
	truStakingSubspace := app.paramsKeeper.Subspace(trustaking.DefaultParamspace)
	truSlashingSubspace := app.paramsKeeper.Subspace(truslashing.DefaultParamspace)
	truDistSubspace := app.paramsKeeper.Subspace(trudist.DefaultParamspace)

	// add cosmos keepers
	app.accountKeeper = auth.NewAccountKeeper(app.codec, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, bank.DefaultCodespace, app.ModuleAccountAddrs())
	app.supplyKeeper = supply.NewKeeper(app.codec, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.codec, keys[staking.StoreKey], app.supplyKeeper, stakingSubspace, staking.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(app.codec, keys[mint.StoreKey], mintSubspace, &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	app.distrKeeper = distr.NewKeeper(app.codec, keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		app.supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(
		app.codec, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace, slashing.DefaultCodespace,
	)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper))
	app.govKeeper = gov.NewKeeper(
		app.codec, keys[gov.StoreKey], govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// add truchain keepers
	app.communityKeeper = community.NewKeeper(
		keys[community.StoreKey],
		app.paramsKeeper.Subspace(community.StoreKey),
		codec,
	)

	app.truBankKeeper = trubank.NewKeeper(
		codec,
		keys[trubank.StoreKey],
		app.bankKeeper,
		trubank2Subspace,
		trubank.DefaultCodespace,
		app.supplyKeeper,
	)

	app.appAccountKeeper = account.NewKeeper(
		keys[account.StoreKey],
		appAccountSubspace,
		codec,
		app.truBankKeeper,
		app.accountKeeper,
		app.supplyKeeper,
	)

	app.claimKeeper = claim.NewKeeper(
		keys[claim.StoreKey],
		app.paramsKeeper.Subspace(claim.StoreKey),
		codec,
		app.appAccountKeeper,
		app.communityKeeper,
	)

	app.truStakingKeeper = trustaking.NewKeeper(
		codec,
		keys[trustaking.StoreKey],
		app.appAccountKeeper,
		app.truBankKeeper,
		app.claimKeeper,
		app.supplyKeeper,
		truStakingSubspace,
		trustaking.DefaultCodespace,
	)

	app.truSlashingKeeper = truslashing.NewKeeper(
		keys[truslashing.StoreKey],
		truSlashingSubspace,
		codec,
		app.truBankKeeper,
		app.truStakingKeeper,
		app.appAccountKeeper,
		app.claimKeeper,
	)

	app.truDistributionKeeper = trudist.NewKeeper(
		keys[trudist.StoreKey],
		truDistSubspace,
		codec,
		app.truBankKeeper,
		app.accountKeeper,
		app.supplyKeeper,
		app.distrKeeper,
	)

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
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

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		// trustory modules
		community.NewAppModule(app.communityKeeper),
		claim.NewAppModule(app.claimKeeper),
		trubank.NewAppModule(app.truBankKeeper),
		account.NewAppModule(app.appAccountKeeper),
		trustaking.NewAppModule(app.truStakingKeeper),
		truslashing.NewAppModule(app.truSlashingKeeper),
		trudist.NewAppModule(app.truDistributionKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName, trudist.ModuleName)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName, trustaking.ModuleName, truslashing.ModuleName, account.ModuleName)

	// genutils must occur after staking so that pools are properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, supply.ModuleName,
		crisis.ModuleName, genutil.ModuleName,
		community.ModuleName, claim.ModuleName, trubank.ModuleName,
		account.ModuleName, trustaking.ModuleName, truslashing.ModuleName, trudist.ModuleName)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	//app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	// The AnteHandler handles signature verification and transaction pre-processing
	// TODO [shanev]: see https://github.com/TruStory/truchain/issues/364
	// Add this back after fixing issues with signature verification
	//app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
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
	codec.RegisterEvidences(cdc)

	return cdc.Seal()
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *TruChain) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
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
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *TruChain) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
