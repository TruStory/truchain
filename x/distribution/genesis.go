package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Params Params `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// InitGenesis initializes state from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Params: keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if !data.Params.UserGrowthAllocation.IsPositive() {
		return fmt.Errorf("Param: UserGrowthAllocation must be positive")
	}

	if !data.Params.UserRewardAllocation.IsPositive() {
		return fmt.Errorf("Param: UserRewardAllocation must be positive")
	}

	return nil
}
