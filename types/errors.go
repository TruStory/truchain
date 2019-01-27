package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// App errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = 1

	CodeInvalidArgumentMsg sdk.CodeType = 101
)

// ErrInvalidArgumentMsg creates an error when `Msg` validation fails
func ErrInvalidArgumentMsg() sdk.Error {
	return sdk.NewError(
		DefaultCodespace, CodeInvalidArgumentMsg, "Invalid argument")
}
