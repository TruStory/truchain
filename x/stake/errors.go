package stake

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Stake errors reserve 1500 ~ 1599.
const (
	DefaultCodespace sdk.CodespaceType = "stake"

	CodeArgumentTooShort  sdk.CodeType = 1501
	CodeArgumentTooLong   sdk.CodeType = 1502
	CodeInvalidStoryState sdk.CodeType = 1503
)

// ErrArgumentTooShortMsg throws for an invalid argument
func ErrArgumentTooShortMsg(argument string, len int) sdk.Error {
	msg := "Argument too short: %s. Must be greater than %d characters."

	return sdk.NewError(
		DefaultCodespace,
		CodeArgumentTooShort,
		fmt.Sprintf(msg, argument, len))
}

// ErrArgumentTooLongMsg throws for an invalid argument
func ErrArgumentTooLongMsg(len int) sdk.Error {
	msg := "Argument too long. Must be less than %d characters."

	return sdk.NewError(
		DefaultCodespace,
		CodeArgumentTooLong,
		fmt.Sprintf(msg, len))
}

// ErrInvalidStoryState throws when story not pending
func ErrInvalidStoryState(state string) sdk.Error {
	msg := "Story can only be staked when it's pending. Current state is: %s"

	return sdk.NewError(
		DefaultCodespace, CodeInvalidStoryState, fmt.Sprintf(msg, state))
}
