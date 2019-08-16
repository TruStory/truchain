package account

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryAppAccount      = "account"
	QueryAppAccounts     = "accounts"
	QueryPrimaryAccount  = "primary_account"
	QueryPrimaryAccounts = "primary_accounts"
	QueryParams          = "params"
)

// QueryAppAccountParams are params for querying app accounts by address queries
type QueryAppAccountParams struct {
	Address sdk.AccAddress `json:"address"`
}

// QueryAppAccountsParams are params for querying app accounts by address queries
type QueryAppAccountsParams struct {
	Addresses []sdk.AccAddress `json:"addresses"`
}

// QueryPrimaryAccountParams are params for querying app accounts by address queries
type QueryPrimaryAccountParams struct {
	Address sdk.AccAddress `json:"address"`
}

// QueryPrimaryAccountParams are params for querying app accounts by address queries
type QueryPrimaryAccountsParams struct {
	Addresses []sdk.AccAddress `json:"addresses"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAppAccount:
			return queryAppAccount(ctx, request, keeper)
		case QueryAppAccounts:
			return queryAppAccounts(ctx, request, keeper)
		case QueryPrimaryAccount:
			return queryPrimaryAccount(ctx, request, keeper)
		case QueryPrimaryAccounts:
			return queryPrimaryAccounts(ctx, request, keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
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

	appAccount, ok := k.getAppAccount(ctx, params.Address)
	if !ok {
		return nil, ErrAppAccountNotFound(params.Address)
	}

	result, jsonErr := codec.MarshalJSONIndent(k.codec, appAccount)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}

	return result, nil
}

func queryPrimaryAccount(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryPrimaryAccountParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	primaryAccount, err := k.PrimaryAccount(ctx, params.Address)
	if err != nil {
		return nil, ErrAppAccountNotFound(params.Address)
	}

	result, jsonErr := codec.MarshalJSONIndent(k.codec, primaryAccount)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}

	return result, nil
}

func queryAppAccounts(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryAppAccountsParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	accounts := make([]AppAccount, 0, len(params.Addresses))

	for _, addr := range params.Addresses {
		appAccount, ok := k.getAppAccount(ctx, addr)
		if !ok {
			return nil, ErrAppAccountNotFound(addr)
		}
		accounts = append(accounts, appAccount)
	}

	result, jsonErr := codec.MarshalJSONIndent(k.codec, accounts)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return result, nil
}

func queryPrimaryAccounts(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryPrimaryAccountsParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}
	accounts := make([]PrimaryAccount, 0, len(params.Addresses))
	for _, addr := range params.Addresses {
		primaryAccount, err := k.PrimaryAccount(ctx, addr)
		if err != nil {
			return nil, ErrAppAccountNotFound(addr)
		}
		accounts = append(accounts, primaryAccount)
	}
	result, jsonErr := codec.MarshalJSONIndent(k.codec, accounts)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return result, nil
}

func queryParams(ctx sdk.Context, keeper Keeper) (result []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)

	result, jsonErr := ModuleCodec.MarshalJSON(params)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marsal result to JSON", jsonErr.Error()))
	}

	return result, nil
}

func unmarshalQueryParams(request abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	err := ModuleCodec.UnmarshalJSON(request.Data, params)
	if err != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
		return
	}
	return
}
