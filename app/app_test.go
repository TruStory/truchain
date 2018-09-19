package app

import (
	"os"
	"testing"

	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setGenesis(baseApp *TruChain, accounts ...*types.AppAccount) (types.GenesisState, error) {
	genAccts := make([]*types.GenesisAccount, len(accounts))
	for i, appAct := range accounts {
		genAccts[i] = types.NewGenesisAccount(appAct)
	}

	genesisState := types.GenesisState{Accounts: genAccts}
	stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genesisState)
	if err != nil {
		return types.GenesisState{}, err
	}

	// initialize and commit the chain
	baseApp.InitChain(abci.RequestInitChain{
		Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	})
	baseApp.Commit()

	return genesisState, nil
}

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	baseApp := NewTruChain(logger, db)
	addr := sdk.AccAddress([]byte{2, 3})

	// construct some test coins
	coins, err := sdk.ParseCoins("77trustake,99bitcoincred")
	require.Nil(t, err)

	// create an auth.BaseAccount for the given test account and set it's coins
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	err = baseAcct.SetCoins(coins)
	require.Nil(t, err)

	// create a new test AppAccount with the given auth.BaseAccount
	appAcct := types.NewAppAccount("foobar", baseAcct)
	genState, err := setGenesis(baseApp, appAcct)
	require.Nil(t, err)

	// create a context for the BaseApp
	ctx := baseApp.BaseApp.NewContext(true, abci.Header{})
	res := baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, appAcct, res)

	// reload app and ensure the account is still there
	baseApp = NewTruChain(logger, db)

	stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genState)
	require.Nil(t, err)

	// initialize the chain with the expected genesis state
	baseApp.InitChain(abci.RequestInitChain{
		Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	})

	ctx = baseApp.BaseApp.NewContext(true, abci.Header{})
	res = baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, appAcct, res)
}
