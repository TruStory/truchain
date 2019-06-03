package trubank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState contains all history of transactions
type GenesisState struct {
	Transactions        []Transaction  `json:"transactions"`
	RewardBrokerAddress sdk.AccAddress `json:"reward_broker_address"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	rewardBrokerAddress, err := sdk.AccAddressFromBech32("cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w")
	if err != nil {
		panic(err)
	}
	return GenesisState{
		Transactions:        make([]Transaction, 0),
		RewardBrokerAddress: rewardBrokerAddress,
	}
}

// InitGenesis initializes state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, t := range data.Transactions {
		keeper.setTransaction(ctx, t)
		keeper.trubankList.AppendToUser(ctx, keeper, t.Creator, t.ID)
	}
	keeper.SetLen(ctx, int64(len(data.Transactions)))
	keeper.setRewardBrokerAddress(ctx, data.RewardBrokerAddress)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	transactions := make([]Transaction, 0)
	err := keeper.EachPrefix(ctx, keeper.StorePrefix(), func(bz []byte) bool {
		var tx Transaction
		keeper.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &tx)
		transactions = append(transactions, tx)
		return true
	})
	if err != nil {
		panic(err)
	}
	rewardBrokerAddress, err2 := keeper.GetRewardBrokerAddress(ctx)
	if err2 != nil {
		rewardBrokerAddress = DefaultGenesisState().RewardBrokerAddress
	}
	return GenesisState{
		Transactions:        transactions,
		RewardBrokerAddress: rewardBrokerAddress,
	}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
