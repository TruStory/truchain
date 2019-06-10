package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Auth errors reserve 200 ~ 299.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeAppAccountNotFound sdk.CodeType = 201
)

// ErrAppAccountNotFound throws an error when the searched AppAccount is not found
func ErrAppAccountNotFound(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeAppAccountNotFound, fmt.Sprintf("AppAccount not found with ID: %d", id))
}
