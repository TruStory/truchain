package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Stake errors reserve 1500 ~ 1599.
const (
	DefaultCodespace sdk.CodespaceType = "stake"

	CodeInvalidArgument sdk.CodeType = 1501
)

// ErrInvalidArgumentMsg creates an error when the searched entity is not found
func ErrInvalidArgumentMsg(argument string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace, CodeInvalidArgument, "Invalid argument "+argument)
}
