package staking

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ModuleName is the name of this module
const ModuleName = "trustaking"

// AppModuleBasic defines the internal data for the module
// ----------------------------------------------------------------------------
type AppModuleBasic struct{}

// Name define the name of the module
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the types needed for amino encoding/decoding
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis creates the default genesis state for testing
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCodec.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis validates the genesis state
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCodec.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the staking module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	// no REST routes, use GraphQL API endpoint
	// i.e: rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the staking module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the staking module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// AppModule defines external data for the module
// ----------------------------------------------------------------------------
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a NewAppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// RegisterInvariants enforces registering of invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route defines the key for the route
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler creates the handler for the staking module
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute defines the querier route
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler creates a new querier handler
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis enforces the creation of the genesis state for the staking module
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCodec.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis enforces exporting this module's data to a genesis file
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCodec.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the supply module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the supply module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
