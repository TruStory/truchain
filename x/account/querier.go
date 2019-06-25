package account


import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryAppAccount = "account"
)

// QueryAppAccountParams are params for querying app accounts by address queries
type QueryAppAccountParams struct {
	Address sdk.AccAddress
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAppAccount:
			return queryAppAccount(ctx, request, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown truchain query endpoint: auth/%s", path[0]))
		}
	}
}

func queryAppAccount(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryAppAccountParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	appAccount, err := k.getAccount(ctx, params.Address)
	if err != nil {
		return
	}

	result, jsonErr := k.codec.MarshalJSON(appAccount)
	if jsonErr != nil {
		panic(jsonErr)
	}
	return result, nil
}

func unmarshalQueryParams(request abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	err := json.Unmarshal(request.Data, params)
	if err != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
		return
	}
	return
}
