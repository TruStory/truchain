package story

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all story state that must be provided at genesis
type GenesisState struct {
	Stories    []Story `json:"stories"`
	StoryQueue []int64 `json:"story_queue"`
	Params     Params  `json:"params"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, s := range data.Stories {
		keeper.setStory(ctx, s)
	}
	// for _, storyID := range data.StoryQueue {
	// keeper.pendingStoryList
	// }
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	return GenesisState{
		Stories: keeper.Stories(ctx),
		Params:  params,
	}
}
