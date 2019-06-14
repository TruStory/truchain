package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Slashing errors reserve 500 ~ 599.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeSlashNotFound  sdk.CodeType = 501
	ErrorCodeInvalidStake   sdk.CodeType = 502
	ErrorCodeInvalidCreator sdk.CodeType = 503
)

// ErrSlashNotFound throws an error when the searched slash is not found
func ErrSlashNotFound(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeSlashNotFound, fmt.Sprintf("Slash not found with ID: %d", id))
}

// ErrInvalidStake throws an error when the stake is invalid
func ErrInvalidStake(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidStake, fmt.Sprintf("Invalid stake with ID: %d", id))
}

// ErrInvalidCreator throws an error when the creator is not an admin
func ErrInvalidCreator(address sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidCreator, fmt.Sprintf("Creator: %d is not an admin", address))
}
