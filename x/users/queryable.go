package users

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath             = "users"
	QueryUsersByAddresses = "addresses"
)

// QueryUsersByAddressesParams are params for users by address queries
type QueryUsersByAddressesParams struct {
	Addresses []string `json:"addresses"`
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(cdc *amino.Codec, k auth.AccountKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryUsersByAddresses:
			return queryUsersByAddresses(ctx, req, cdc, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryUsersByAddresses(
	ctx sdk.Context,
	req abci.RequestQuery,
	cdc *amino.Codec,
	k auth.AccountKeeper) (res []byte, err sdk.Error) {

	// get query params
	params, err := unmarshalQueryParams(cdc, req)

	if err != nil {
		return
	}

	users := make([]User, len(params.Addresses))

	for i, a := range params.Addresses {
		addr, err := sdk.AccAddressFromBech32(a)
		if err != nil {
			return res, sdk.NewError(0, 0, "Error decoding address: "+err.Error())
		}
		account := k.GetAccount(ctx, addr)
		if account != nil {
			users[i] = NewUser(account)
		} else {
			users[i] = User{}
		}
	}

	// return users JSON bytes
	return marshalUsers(cdc, users)
}

// unmarshal query params into struct
func unmarshalQueryParams(cdc *amino.Codec, req abci.RequestQuery) (params QueryUsersByAddressesParams, sdkErr sdk.Error) {
	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return params,
			sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
	}
	return
}

func marshalUsers(cdc *amino.Codec, users []User) (res []byte, sdkErr sdk.Error) {
	res, err := cdc.MarshalJSON(users)

	if err != nil {
		panic("Could not marshal result to JSON")
	}

	return
}
