package challenge

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story errors reserve 1000 ~ 1099.
const (
	DefaultCodespace sdk.CodespaceType = 10

	CodeNotFound   sdk.CodeType = 1001
	CodeInvalidMsg sdk.CodeType = 1002
)

// ErrNotFound throws an error when the searched story is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeNotFound,
		"Challenge with id "+fmt.Sprintf("%d", id)+" not found")
}

// ErrInvalidMsg throws an error when `Msg` validation fails
func ErrInvalidMsg(value interface{}) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeInvalidMsg,
		"Invalid message field: "+fmt.Sprintf("%s", reflect.TypeOf(value).String()))
}
