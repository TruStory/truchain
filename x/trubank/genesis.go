package trubank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState contains all history of transactions
type GenesisState struct {
	Transactions []Transaction `json:"transactions"`
}

// InitGenesis initializes state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, t := range data.Transactions {
		keeper.setTransaction(ctx, t)
		keeper.trubankList.AppendToUser(ctx, keeper, t.Creator, t.ID)
	}
	keeper.SetLen(ctx, int64(len(data.Transactions)))
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	transactions := make([]Transaction, 0)
	prefix := fmt.Sprintf("%s:id:", StoreKey)
	err := keeper.EachPrefix(ctx, prefix, func(bz []byte) bool {
		var tx Transaction
		keeper.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &tx)
		transactions = append(transactions, tx)
		return true
	})
	if err != nil {
		panic(err)
	}
	return GenesisState{
		Transactions: transactions,
	}
}
