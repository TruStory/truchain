package argument

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story errors reserve 1500 ~ 1599.
const (
	DefaultCodespace sdk.CodespaceType = "argument"

	CodeInvalidArgumentID sdk.CodeType = 1501
	CodeArgumentTooShort  sdk.CodeType = 1502
	CodeArgumentTooLong   sdk.CodeType = 1503
)

// ErrInvalidArgumentID throws an error on invalid title
func ErrInvalidArgumentID() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidArgumentID, "Invalid argument ID")
}

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
