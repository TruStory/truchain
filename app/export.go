package app

import (
	"encoding/json"

	"github.com/TruStory/truchain/x/trubank"

	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *TruChain) ExportAppStateAndValidators() (
	appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	accounts := []*auth.BaseAccount{}
	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address:       acc.GetAddress(),
			Coins:         acc.GetCoins(),
			PubKey:        acc.GetPubKey(),
			AccountNumber: acc.GetAccountNumber(),
			Sequence:      acc.GetSequence(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenesisState{
		ArgumentData:   argument.ExportGenesis(ctx, app.argumentKeeper),
		Accounts:       accounts,
		AuthData:       auth.ExportGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper),
		BankData:       bank.ExportGenesis(ctx, app.bankKeeper),
		BackingData:    backing.ExportGenesis(ctx, app.backingKeeper),
		CategoryData:   category.ExportGenesis(ctx, app.categoryKeeper),
		ChallengeData:  challenge.ExportGenesis(ctx, app.challengeKeeper),
		ExpirationData: expiration.ExportGenesis(ctx, app.expirationKeeper),
		StakeData:      stake.ExportGenesis(ctx, app.stakeKeeper),
		StoryData:      story.ExportGenesis(ctx, app.storyKeeper),
		TrubankData:    trubank.ExportGenesis(ctx, app.truBankKeeper),
	}

	appState, err = codec.MarshalJSONIndent(app.codec, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}
