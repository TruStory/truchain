package app

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	AuthData   auth.GenesisState   `json:"auth"`
	BankData   bank.GenesisState   `json:"bank"`
	Accounts   []*auth.BaseAccount `json:"accounts"`
	Categories []category.Category `json:"categories"`
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
	bank.InitGenesis(ctx, app.coinKeeper, genesisState.BankData)
	category.InitGenesis(ctx, app.categoryKeeper, genesisState.Categories)

	return abci.ResponseInitChain{}
}
