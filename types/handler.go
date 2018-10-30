package types

import (
	"encoding/binary"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Utilities for all module handlers

// ErrMsgHandler returns an unknown Msg request error result
func ErrMsgHandler(msg sdk.Msg) sdk.Result {
	if mType := reflect.TypeOf(msg); mType != nil {
		errMsg := "Unrecognized Msg type: " + mType.Name()
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	return sdk.ErrUnknownRequest("Unrecognized Msg").Result()
}

// Result returns a successful handler result with id of the type
// encoded as binary data
func Result(id int64) sdk.Result {
	return sdk.Result{Data: i2b(id)}
}

// i2b converts an int64 into a byte array
func i2b(x int64) []byte {
	var b [binary.MaxVarintLen64]byte
	return b[:binary.PutVarint(b[:], x)]
}
