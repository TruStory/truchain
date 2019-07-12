package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Slashing errors reserve 500 ~ 599.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeSlashNotFound        sdk.CodeType = 501
	ErrorCodeInvalidStake         sdk.CodeType = 502
	ErrorCodeInvalidArgument      sdk.CodeType = 503
	ErrorCodeMaxSlashCountReached sdk.CodeType = 504
	ErrorCodeInvalidCreator       sdk.CodeType = 505
	ErrorCodeNotEnoughEarnedStake sdk.CodeType = 506
	ErrorCodeAlreadySlashed       sdk.CodeType = 507
)

// ErrSlashNotFound throws an error when the searched slash is not found
func ErrSlashNotFound(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeSlashNotFound, fmt.Sprintf("Slash not found with ID: %d", id))
}

// ErrInvalidStake throws an error when the stake is invalid
func ErrInvalidStake(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidStake, fmt.Sprintf("Invalid stake with ID: %d", id))
}

// ErrInvalidArgument throws an error when the argument is invalid
func ErrInvalidArgument(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidArgument, fmt.Sprintf("Invalid argument with ID: %d", id))
}

// ErrMaxSlashCountReached throws an error when the max slash count on a stake is reached
func ErrMaxSlashCountReached(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidStake, fmt.Sprintf("Max slash count reached for stake with ID: %d", id))
}

// ErrInvalidCreator throws an error when the creator is not an admin
func ErrInvalidCreator(address sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidCreator, fmt.Sprintf("Creator: %d is not an admin", address))
}

// ErrNotEnoughEarnedStake throws an error when the creator is not an admin
func ErrNotEnoughEarnedStake(address sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeNotEnoughEarnedStake, fmt.Sprintf("Creator: %d does not have enough earned stake", address))
}

// ErrAlreadySlashed throws an error when the creator has already slashed the stake previously
func ErrAlreadySlashed() sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeAlreadySlashed, "Creator cannot slash a stake more than once")
}
