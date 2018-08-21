package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TruStory errors reserve 700 ~ 799.
const (
	DefaultCodespace sdk.CodespaceType = 7

	CodeInvalidOption     sdk.CodeType = 701
	CodeInvalidBody       sdk.CodeType = 702
	CodeInvalidStoryID    sdk.CodeType = 703
	CodeStoryNotFound     sdk.CodeType = 704
	CodeInvalidAmount     sdk.CodeType = 705
	CodeInvalidBondPeriod sdk.CodeType = 706
	CodeInvalidURL        sdk.CodeType = 707
	CodeInvalidCategory   sdk.CodeType = 708
	CodeInvalidStoryType  sdk.CodeType = 709
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

// ErrInvalidOption throws an error on invalid option
func ErrInvalidOption(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidOption, msg)
}

// ErrInvalidBody throws an error on invalid title
func ErrInvalidBody(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidBody, msg)
}

// ErrInvalidCategory throws an error when the category is invalid
func ErrInvalidCategory(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidCategory, msg)
}

// ErrInvalidStoryID throws an error on invalid storyID
func ErrInvalidStoryID(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidStoryID, msg)
}

// ErrInvalidStoryType throws an error on invalid story type
func ErrInvalidStoryType(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidStoryType, msg)
}

// ErrInvalidAmount throws an error when a stake amount is invalid
func ErrInvalidAmount(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidAmount, msg)
}

// ErrInvalidBondPeriod throws an error when bond period is invalid
func ErrInvalidBondPeriod(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidBondPeriod, msg)
}

// ErrInvalidURL throws an error when a URL in invalid
func ErrInvalidURL(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidURL, msg)
}

// ErrStoryNotFound throws an error when the searched story is not found
func ErrStoryNotFound(storyID int64) sdk.Error {
	return newError(DefaultCodespace, CodeStoryNotFound, "Story with id "+
		strconv.Itoa(int(storyID))+" not found")
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
