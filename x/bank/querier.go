package bank

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryTransactionsByAddress = "transactions_by_address"
)

// QueryTransactionsByAddress query transactions params for a specific address.
type QueryTransactionsByAddressParams struct {
	Address   sdk.AccAddress    `json:"address"`
	Types     []TransactionType `json:"types,omitempty"`
	SortOrder SortOrderType     `json:"sortOrder,omitempty"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryTransactionsByAddress:
			return queryTransactionsByAddress(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown bank query endpoint")
		}
	}
}

func queryTransactionsByAddress(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTransactionsByAddressParams
	err := json.Unmarshal(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	sortOrder := SortAsc
	if params.SortOrder.Valid() {
		sortOrder = params.SortOrder
	}
	transactions := keeper.TransactionsByAddress(ctx,
		params.Address,
		FilterByTransactionType(params.Types...),
		SortOrder(sortOrder),
		Limit(params.Limit),
		Offset(params.Offset),
	)
	return keeper.codec.MustMarshalJSON(transactions), nil
}
