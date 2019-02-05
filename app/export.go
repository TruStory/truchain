package app

import (
	"encoding/json"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/vote"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *TruChain) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	// iterate over accounts
	// accounts := []*types.GenesisAccount{}
	accounts := []GenesisAccount{}
	appendAccountsFn := func(acc auth.Account) bool {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	// Export current state to each Keeper (TESTING)
	app.storyKeeper.ExportState(ctx, DefaultNodeHome, app.LastBlockHeight())
	app.categoryKeeper.ExportState(ctx, DefaultNodeHome, app.LastBlockHeight())
	app.challengeKeeper.ExportState(ctx, DefaultNodeHome, app.LastBlockHeight())
	app.gameKeeper.ExportState(ctx, DefaultNodeHome, app.LastBlockHeight())
	// app.backingKeeper.ExportState(ctx)

	genState := NewGenesisState(
		accounts,
		story.ExportGenesis(ctx, app.storyKeeper),
		category.ExportGenesis(ctx, app.categoryKeeper),
		backing.ExportGenesis(ctx, app.backingKeeper),
		challenge.ExportGenesis(ctx, app.challengeKeeper),
		game.ExportGenesis(ctx, app.gameKeeper),
		vote.ExportGenesis(ctx, app.voteKeeper),
	)

	appState, err = codec.MarshalJSONIndent(app.codec, genState)
	if err != nil {
		return nil, nil, err
	}

	// TODO export validators
	// validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, err
}
