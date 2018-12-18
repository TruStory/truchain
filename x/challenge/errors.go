package challenge

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenge errors reserve 1000 ~ 1099.
const (
	DefaultCodespace sdk.CodespaceType = 10

	CodeNotFound           sdk.CodeType = 1001
	CodeInvalidMsg         sdk.CodeType = 1002
	CodeInvalidVote        sdk.CodeType = 1003
	CodeDuplicateChallenge sdk.CodeType = 1004
	CodeInvalidEvidenceURL sdk.CodeType = 1005
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

// ErrDuplicateChallenge creates an error when more than one challenge is attempted on a story
func ErrDuplicateChallenge(gameID int64, user sdk.AccAddress) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeDuplicateChallenge,
		"Game with id "+fmt.Sprintf("%d", gameID)+" has already been challenged by "+user.String())
}

// ErrInvalidVote creates an error when Challenge has a true vote
func ErrInvalidVote() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeInvalidMsg,
		"Challenges cannot have a TRUE vote.")
}

// ErrInvalidEvidenceURL throws an error when a URL in invalid
func ErrInvalidEvidenceURL(url string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeInvalidEvidenceURL,
		"Invalid evidence URL: "+url)
}
