package game

import sdk "github.com/cosmos/cosmos-sdk/types"

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, gk WriteKeeper) []Game {

	games := gk.Games(ctx)

	return games
}
