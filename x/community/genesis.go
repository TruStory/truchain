package community

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Communities []Community `json:"communities"`
	Params      Params      `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(communities []Community, params Params) GenesisState {
	return GenesisState{
		Communities: communities,
		Params:      params,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	community1 := NewCommunity("crypto", "Cryptocurrency", "description string", time.Now())
	community2 := NewCommunity("meme", "Memes", "description string", time.Now())

	return NewGenesisState(Communities{community1, community2}, DefaultParams())
}

// InitGenesis initializes community state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, community := range data.Communities {
		keeper.setCommunity(ctx, community)
	}
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Communities: keeper.Communities(ctx),
		Params:      keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinNameLength < 1 {
		return fmt.Errorf("Param: MinNameLength, must have a positive value")
	}

	if data.Params.MaxNameLength < 1 || data.Params.MaxNameLength < data.Params.MinNameLength {
		return fmt.Errorf("Param: MaxNameLength, must have a positive value and be larger than MinNameLength")
	}

	if data.Params.MinIDLength < 1 {
		return fmt.Errorf("Param: MinIDLength, must have a positive value")
	}

	if data.Params.MaxIDLength < 1 || data.Params.MaxIDLength < data.Params.MinIDLength {
		return fmt.Errorf("Param: MaxIDLength, must have a positive value and be larger than MinIDLength")
	}

	if data.Params.MaxDescriptionLength < 1 {
		return fmt.Errorf("Param: MaxDescriptionLength, must have a positive value")
	}

	if len(data.Params.CommunityAdmins) < 1 {
		return fmt.Errorf("Param: CommunityAdmins, must have atleast one admin")
	}

	return nil
}
