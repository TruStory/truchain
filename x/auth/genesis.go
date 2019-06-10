package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Registrar sdk.AccAddress `json:"registrar"`
	Params    Params         `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	// TODO: figure out where to get this address from
	registrar, err := sdk.AccAddressFromBech32("cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w")
	if err != nil {
		panic(err)
	}
	return GenesisState{
		Registrar: registrar,
		Params:    DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Params: keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MaxSlashCount < 1 {
		return fmt.Errorf("Param: MaxSlashCount, must have a positive value")
	}

	return nil
}
