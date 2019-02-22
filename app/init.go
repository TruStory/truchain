package app

import (
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	AuthData       auth.GenesisState       `json:"auth"`
	BankData       bank.GenesisState       `json:"bank"`
	Accounts       []*auth.BaseAccount     `json:"accounts"`
	Categories     []category.Category     `json:"categories"`
	ChallengeData  challenge.GenesisState  `json:"challenge"`
	ExpirationData expiration.GenesisState `json:"expiration"`
	StoryData      story.GenesisState      `json:"story"`
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *TruChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.codec.UnmarshalJSON(stateJSON, genesisState)

	if err != nil {
		panic(err)
	}

	for i, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		if i == 1 { // TODO: more robust way of identifying registrar account [notduncansmith]
			err := acc.SetPubKey(app.registrarKey.PubKey())
			if err != nil {
				panic(err)
			}
		}

		app.accountKeeper.SetAccount(ctx, acc)
	}

	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	category.InitGenesis(ctx, app.categoryKeeper, genesisState.Categories)
	challenge.InitGenesis(ctx, app.challengeKeeper, genesisState.ChallengeData)
	expiration.InitGenesis(ctx, app.expirationKeeper, genesisState.ExpirationData)
	story.InitGenesis(ctx, app.storyKeeper, genesisState.StoryData)

	return abci.ResponseInitChain{}
}
