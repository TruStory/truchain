package types

import (
	"encoding/json"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Utilities for all module handlers

// IDResult is a common result struct containing the id of a created object
type IDResult struct {
	ID int64 `json:"id"`
}

// ErrMsgHandler returns an unknown Msg request error result
func ErrMsgHandler(msg sdk.Msg) sdk.Result {
	if mType := reflect.TypeOf(msg); mType != nil {
		errMsg := "Truchain: Unrecognized Msg type: " + mType.Name()
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	return sdk.ErrUnknownRequest("Unrecognized Msg").Result()
}

// Result returns a successful handler result with id of the type
// encoded as binary data
func Result(id int64) sdk.Result {
	bz, err := json.Marshal(IDResult{ID: id})

	if err != nil {
		panic(err)
	}

	return sdk.Result{Data: bz}
}
