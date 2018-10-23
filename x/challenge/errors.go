package challenge

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenge errors reserve 1000 ~ 1099.
const (
	DefaultCodespace sdk.CodespaceType = 10

	CodeNotFound            sdk.CodeType = 1001
	CodeInvalidMsg          sdk.CodeType = 1002
	CodeDuplicate           sdk.CodeType = 1003
	CodeNotFoundChallenger  sdk.CodeType = 1004
	CodeDuplicateChallenger sdk.CodeType = 1005
)

// ErrNotFound creates an error when the searched entity is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeNotFound,
		"Challenge with id "+fmt.Sprintf("%d", id)+" not found")
}

// ErrInvalidMsg creates an error when `Msg` validation fails
func ErrInvalidMsg(value interface{}) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeInvalidMsg,
		"Invalid message field: "+fmt.Sprintf("%s", reflect.TypeOf(value).String()))
}

// ErrDuplicate creates an error when more than one challenge is attempted on a story
func ErrDuplicate(storyID int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeDuplicate,
		"Story with id "+fmt.Sprintf("%d", storyID)+" has already been challenged")
}

// ErrNotFoundChallenger creates an error for not finding a challenger
func ErrNotFoundChallenger(id int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeNotFoundChallenger,
		"Challenger not found for challenge "+fmt.Sprintf("%d", id),
	)
}

// ErrDuplicateChallenger creates an error when more than one challenge is attempted on a story
func ErrDuplicateChallenger(id int64, user sdk.AccAddress) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeDuplicateChallenger,
		"Challenge with id "+fmt.Sprintf("%d", id)+" has already been challenged by "+user.String())
}
