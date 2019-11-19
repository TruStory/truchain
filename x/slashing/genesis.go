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
		Slashes: []Slash{},
		Params:  DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes slashing state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	counter := make(map[uint64]uint64)
	for _, slash := range data.Slashes {
		keeper.setSlash(ctx, slash)
		count, ok := counter[slash.ArgumentID]
		if !ok {
			count = 0
		}
		count = count + 1
		counter[slash.ArgumentID] = count

		keeper.setCreatorSlash(ctx, slash.Creator, slash.ID)
		keeper.setSlashCount(ctx, slash.ArgumentID, count)
		keeper.setArgumentSlash(ctx, slash.ArgumentID, slash.ID)
		keeper.setArgumentSlasherSlash(ctx, slash.ArgumentID, slash.ID, slash.Creator)

	}
	keeper.setSlashID(ctx, uint64(len(data.Slashes)+1))
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Slashes: keeper.Slashes(ctx),
		Params:  keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinSlashCount < 1 {
		return fmt.Errorf("Param: MinSlashCount, must have a positive value")
	}

	if data.Params.SlashMagnitude < 1 {
		return fmt.Errorf("Param: SlashMagnitude, must have a positive value")
	}

	if data.Params.SlashMinStake.IsNegative() {
		return fmt.Errorf("Param: SlashMinStake, cannot be a negative value")
	}

	if len(data.Params.SlashAdmins) < 1 {
		return fmt.Errorf("Param: SlashAdmins, must have atleast one admin")
	}

	if data.Params.CuratorShare.IsNegative() {
		return fmt.Errorf("Param: CuratorShare, cannot be a negative value")
	}

	return nil
}
