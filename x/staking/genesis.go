package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Arguments []Argument `json:"arguments"`
	Params    Params     `json:"params"`
	Stakes    []Stake    `json:"stakes"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(arguments []Argument, params Params, stakes []Stake) GenesisState {
	return GenesisState{
		Arguments: arguments,
		Params:    params,
		Stakes:    stakes,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:    DefaultParams(),
		Stakes:    make([]Stake, 0),
		Arguments: make([]Argument, 0),
	}
}

// InitGenesis initializes staking state from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, a := range data.Arguments {
		k.setArgument(ctx, a)
		k.setClaimArgument(ctx, a.ClaimID, a.ID)
		k.serUserArgument(ctx, a.Creator, a.ID)
	}
	for _, s := range data.Stakes {
		k.setStake(ctx, s)
		if !s.Expired {
			k.InsertActiveStakeQueue(ctx, s.ID, s.EndTime)
		}
		k.setArgumentStake(ctx, s.ArgumentID, s.ID)
		k.setUserStake(ctx, s.Creator, s.CreatedTime, s.ID)
	}
	k.setArgumentID(ctx, uint64(len(data.Arguments)+1))
	k.setStakeID(ctx, uint64(len(data.Stakes)+1))
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
	return nil
}
