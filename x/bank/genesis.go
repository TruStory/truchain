package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Transactions []Transaction `json:"transactions"`
	Params       Params        `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, transactions []Transaction) GenesisState {
	return GenesisState{
		Params:       params,
		Transactions: transactions,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:       DefaultParams(),
		Transactions: make([]Transaction, 0),
	}
}

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	for _, tx := range data.Transactions {
		keeper.Set(ctx, tx.ID, tx)
		keeper.PushWithAddress(ctx, keeper.storeKey, accountKey, tx.ID, tx.AppAccountAddress)
	}
	keeper.SetLen(ctx, uint64(len(data.Transactions)))
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Params:       keeper.GetParams(ctx),
		Transactions: keeper.Transactions(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.RewardBrokerAddress.Empty() {
		return fmt.Errorf("param: RewardBrokerAddress, a valid address must be provided")
	}
	return nil
}
