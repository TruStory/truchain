package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Communities []Community `json:"communities"`
	Params      MsgParams   `json:"msg_params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(communities []Community, msgParams MsgParams) GenesisState {
	return GenesisState{
		Communities: communities,
		Params:      msgParams,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState([]Community{}, DefaultMsgParams()) }

// InitGenesis initializes community state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, community := range data.Communities {
		keeper.Set(ctx, community.ID, community)
	}
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Communities: keeper.Communities(ctx),
		Params:      DefaultMsgParams(),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinNameLen < 1 {
		return fmt.Errorf("Param: MinNameLen, must have a positive value")
	}

	if data.Params.MaxNameLen < 1 || data.Params.MaxNameLen < data.Params.MinNameLen {
		return fmt.Errorf("Param: MaxNameLen, must have a positive value and be larger than MinNameLen")
	}

	if data.Params.MinSlugLen < 1 {
		return fmt.Errorf("Param: MinSlugLen, must have a positive value")
	}

	if data.Params.MaxSlugLen < 1 || data.Params.MaxSlugLen < data.Params.MinSlugLen {
		return fmt.Errorf("Param: MaxSlugLen, must have a positive value and be larger than MinSlugLen")
	}

	if data.Params.MaxDescriptionLen < 1 {
		return fmt.Errorf("Param: MaxDescriptionLen, must have a positive value")
	}

	return nil
}
