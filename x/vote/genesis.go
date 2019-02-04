package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, vk WriteKeeper) []TokenVote {

	votes := vk.Votes(ctx)

	return votes
}
