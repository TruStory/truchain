package app

import (
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// GenesisState struct defines the current state
type GenesisState struct {
	Accounts   []GenesisAccount      `json:"accounts"`
	Stories    []story.Story         `json:"stories"`
	Categories []category.Category   `json:"categories"`
	Backings   backing.GenesisState  `json:"backings"`
	Challenges []challenge.Challenge `json:"challenges"`
	Games      game.GenesisState     `json:"games"`
	Votes      []vote.TokenVote      `json:"votes"`
}

// NewGenesisState retturns the current GenesisState
func NewGenesisState(accounts []GenesisAccount,
	storyData []story.Story,
	categoryData []category.Category,
	backingData backing.GenesisState,
	challengeData []challenge.Challenge,
	gameData game.GenesisState,
	voteData []vote.TokenVote) GenesisState {

	return GenesisState{
		Accounts:   accounts,
		Stories:    storyData,
		Categories: categoryData,
		Backings:   backingData,
		Challenges: challengeData,
		Games:      gameData,
		Votes:      voteData,
	}
}

// GenesisAccount reflects a genesis account the application expects in it's
// genesis state.
type GenesisAccount struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccountI gets the sate addresses and cins
func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	gacc := GenesisAccount{
		Address: acc.GetAddress(),
		Coins:   acc.GetCoins(),
	}

	return gacc
}
