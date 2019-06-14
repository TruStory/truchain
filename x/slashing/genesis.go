package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Slashes []Slash `json:"slashes"`
	Params  Params  `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Slashes: nil,
		Params:  DefaultParams(),
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
		//		Slashings: keeper.Slashings(ctx),
		Params: keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MaxStakeSlashCount < 1 {
		return fmt.Errorf("Param: MaxStakeSlashCount, must have a positive value")
	}

	if !data.Params.SlashMagnitude.IsPositive() {
		return fmt.Errorf("Param: SlashMagnitude, must have a positive value")
	}

	if !data.Params.SlashMinStake.IsPositive() {
		return fmt.Errorf("Param: SlashMinStake, must have a positive value")
	}

	if len(data.Params.SlashAdmins) < 1 {
		return fmt.Errorf("Param: SlashAdmins, must have atleast one admin")
	}

	if data.Params.JailTime.Seconds() < 1 {
		return fmt.Errorf("Param: JailTime, must have a positive value")
	}

	return nil
}
