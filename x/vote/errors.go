package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Vote errors reserve 1100 ~ 1199.
const (
	DefaultCodespace sdk.CodespaceType = "vote"

	CodeNotFound         sdk.CodeType = 1201
	CodeDuplicate        sdk.CodeType = 1202
	CodeVotingNotStarted sdk.CodeType = 1203
)

// ErrNotFound creates an error when the searched entity is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeNotFound,
		"Vote with id "+fmt.Sprintf("%d", id)+" not found")
}

// ErrDuplicateVote throws when a vote already has been cast
func ErrDuplicateVote(
	gameID int64, user sdk.AccAddress) sdk.Error {

	return sdk.NewError(
		DefaultCodespace,
		CodeDuplicate,
		"Vote for game "+fmt.Sprintf("%d", gameID)+" has already been cast by user "+user.String())
}

// ErrVotingNotStarted is thrown when a vote is attempted on a story
// that hasn't begun the validation game yet
func ErrVotingNotStarted(storyID int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeVotingNotStarted,
		"Validation game not started for story: "+
			fmt.Sprintf("%d", storyID))
}
