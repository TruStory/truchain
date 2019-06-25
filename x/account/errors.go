package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Auth errors reserve 200 ~ 299.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeAppAccountNotFound     sdk.CodeType = 201
	ErrorCodeAppAccountCreateFailed sdk.CodeType = 202
)

// ErrAppAccountNotFound throws an error when the searched AppAccount is not found
func ErrAppAccountNotFound(address sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeAppAccountNotFound, fmt.Sprintf("AppAccount not found with Address: %s", address))
}

// ErrAppAccountCreateFailed throws an error when creating an account fails
func ErrAppAccountCreateFailed(address sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeAppAccountCreateFailed, fmt.Sprintf("Creating AppAccount failed: %s", address))
}
