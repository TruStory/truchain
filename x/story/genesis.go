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
		keeper.appendStoriesList(ctx, storyIDsByCategoryKey(keeper, s.CategoryID, s.Timestamp, false), s)
	}
	for _, storyID := range data.StoryQueue {
		keeper.storyQueue(ctx).Push(storyID)
	}
	keeper.SetLen(ctx, int64(len(data.Stories)))
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var storyIDs []int64
	var storyID int64
	keeper.storyQueue(ctx).List.Iterate(&storyID, func(uint64) bool {
		storyIDs = append(storyIDs, storyID)
		return false
	})

	return GenesisState{
		Stories:    keeper.Stories(ctx),
		StoryQueue: storyIDs,
		Params:     keeper.GetParams(ctx),
	}
}


func ValidateGenesis(data GenesisState) error {
	return nil
}
