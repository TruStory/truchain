package game

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Games        []Game  `json:"game"`
	PendingGames []int64 `json:"pending_games"`
	GamesQueue   []int64 `json:"game_queue"`
}

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, gk WriteKeeper) (data GenesisState) {

	games := gk.Games(ctx)

	pendingGames := gk.PendingGames(ctx)

	gameQueue := gk.GameQueue(ctx)

	return GenesisState{
		Games:        games,
		PendingGames: pendingGames,
		GamesQueue:   gameQueue,
	}
}
