package game

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Game errors reserve 1100 ~ 1199.
const (
	DefaultCodespace sdk.CodespaceType = 11

	CodeNotFound  sdk.CodeType = 1101
	CodeDuplicate sdk.CodeType = 1102
)

// ErrNotFound creates an error when the searched entity is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeNotFound,
		"Game with id "+fmt.Sprintf("%d", id)+" not found")
}

// ErrExists throws when more than one current game is attempted on a story
func ErrExists(storyID int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeDuplicate,
		"Story with id "+fmt.Sprintf("%d", storyID)+" has an active game")
}
