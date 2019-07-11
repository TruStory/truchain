package argument

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story errors reserve 1500 ~ 1599.
const (
	DefaultCodespace sdk.CodespaceType = "argument"

	CodeInvalidID    sdk.CodeType = 1501
	CodeNotFound     sdk.CodeType = 1502
	CodeBodyTooShort sdk.CodeType = 1503
	CodeBodyTooLong  sdk.CodeType = 1504
	CodeInvalid      sdk.CodeType = 1505
)

// ErrInvalidArgumentID throws an error on invalid title
func ErrInvalidArgumentID() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidID, "Invalid argument ID")
}

// ErrNotFound throws when an argument by id is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotFound, fmt.Sprintf("Invalid argument ID %d", id))
}

// ErrArgumentTooShortMsg throws for an invalid argument
func ErrArgumentTooShortMsg(length int) sdk.Error {
	msg := "Argument body too short. Must be greater than %d characters."

	return sdk.NewError(
		DefaultCodespace,
		CodeBodyTooShort,
		fmt.Sprintf(msg, length))
}

// ErrArgumentTooLongMsg throws for an invalid argument
func ErrArgumentTooLongMsg(length int, maxLength int) sdk.Error {
	msg := "Argument is %d character%s too long. Must be less than %d characters."
	plural := "s"
	if (length - maxLength) == 1 {
		plural = ""
	}

	return sdk.NewError(
		DefaultCodespace,
		CodeBodyTooLong,
		fmt.Sprintf(msg, length-maxLength, plural, maxLength))
}

// ErrInvalidArgument throws when are argument is invalid
func ErrInvalidArgument(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalid, "Invalid argument: "+msg)
}
