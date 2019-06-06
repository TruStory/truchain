package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all story state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// InitGenesis initializes state from genesis file
func InitGenesis(ctx sdk.Context, stakingKeeper Keeper, data GenesisState) {
	stakingKeeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	return GenesisState{
		Params: params,
	}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
