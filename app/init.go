package app

import (
	params "github.com/TruStory/truchain/parameters"
	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *TruChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	genesisState := new(types.GenesisState)
	err := app.codec.UnmarshalJSON(stateJSON, genesisState)

	if err != nil {
		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
		panic(err)
	}

	for i, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
			panic(err)
		}

		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)

		if i == 1 { // TODO: more robust way of identifying registrar account [notduncansmith]
			err := acc.BaseAccount.SetPubKey(app.registrarKey.PubKey())
			if err != nil {
				panic(err)
			}
		}

		app.accountKeeper.SetAccount(ctx, acc)
	}

	// get genesis account address
	genesisAddr := genesisState.Accounts[0].Address

	// persist initial categories on chain
	err = app.categoryKeeper.InitCategories(ctx, genesisAddr, app.categories)

	if err != nil {
		panic(err)
	}

	if params.Features[params.BootstrapFlag] {
		loadTestDB(
			ctx, app.storyKeeper,
			app.accountKeeper,
			app.backingKeeper,
			app.challengeKeeper,
			app.gameKeeper,
		)
	}

	return abci.ResponseInitChain{}
}
