package community

import (
	"fmt"

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
	communities := []Community{
		// {ID: 1, Name: "Cryptocurrency", Slug: "crypto", TotalEarnedStake: sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())},
		// {ID: 2, Name: "Memes", Slug: "meme", TotalEarnedStake: sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())},
	}

	return NewGenesisState(communities, DefaultParams())
}

// InitGenesis initializes community state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, community := range data.Communities {
		keeper.Set(ctx, community.ID, community)
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

	if data.Params.MinSlugLength < 1 {
		return fmt.Errorf("Param: MinSlugLength, must have a positive value")
	}

	if data.Params.MaxSlugLength < 1 || data.Params.MaxSlugLength < data.Params.MinSlugLength {
		return fmt.Errorf("Param: MaxSlugLength, must have a positive value and be larger than MinSlugLength")
	}

	if data.Params.MaxDescriptionLength < 1 {
		return fmt.Errorf("Param: MaxDescriptionLength, must have a positive value")
	}

	return nil
}
