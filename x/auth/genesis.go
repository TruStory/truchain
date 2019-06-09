package auth

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState defines genesis data for the module
type GenesisState struct {
	Auths []Auth `json:"auths"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Auths: nil,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
//		Auths: keeper.Auths(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {

	return nil
}
