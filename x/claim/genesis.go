package claim

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Claims []Claim `json:"claims"`
	Params Params  `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Claims: nil,
		Params: DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, claim := range data.Claims {
		keeper.Set(ctx, claim.ID, claim)
	}
	keeper.SetLen(ctx, uint64(len(data.Claims)))
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Claims: keeper.Claims(ctx),
		Params: keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinClaimLength < 1 {
		return fmt.Errorf("Param: MinClaimLength must have a positive value")
	}
	if data.Params.MaxClaimLength < 1 {
		return fmt.Errorf("Param: MaxClaimLength must have a positive value")
	}

	return nil
}
