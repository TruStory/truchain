package challenge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, ck WriteKeeper) []Challenge {

	return ck.Challenges(ctx)
}
