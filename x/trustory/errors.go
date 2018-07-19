package trustory

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = 7

	// TruStory errors reserve 700 ~ 799.
	CodeInvalidOption  CodeType = 701
	CodeInvalidBody    CodeType = 702
	CodeInvalidStoryID CodeType = 703
	CodeStoryNotFound  CodeType = 704
)

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeInvalidOption:
		return "Invalid option"
	case CodeInvalidBody:
		return "Invalid story body"
	case CodeInvalidStoryID:
		return "Invalid storyID"
	case CodeStoryNotFound:
		return "Story not found"
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

// ErrInvalidStoryID throws an error on invalid proposaID
func ErrInvalidStoryID(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidStoryID, msg)
}

// ErrStoryNotFound throws an error when the searched proposal is not found
func ErrStoryNotFound(storyID int64) sdk.Error {
	return newError(DefaultCodespace, CodeStoryNotFound, "Story with id "+
		strconv.Itoa(int(storyID))+" not found")
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
