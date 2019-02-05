package story

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportGenesis gets all the current stories
func ExportGenesis(ctx sdk.Context, sk WriteKeeper) []Story {

	return sk.StoriesNoSort(ctx)
}
