package backing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story errors reserve 700 ~ 799.
const (
	DefaultCodespace sdk.CodespaceType = "backing"

	CodeQueueEmpty sdk.CodeType = 902
	CodeNotFound   sdk.CodeType = 903
	CodeDuplicate  sdk.CodeType = 904
)

// ErrQueueEmpty throws an error when the searched Queue is not found
func ErrQueueEmpty() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeQueueEmpty, "Backing queue is empty")
}

// ErrNotFound throws an error when the searched backing is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotFound, "Backing with id "+fmt.Sprintf("%d", id)+" not found")
}

// ErrDuplicate throws an error when user has already backed the story
func ErrDuplicate(storyID int64, creator sdk.AccAddress) sdk.Error {
	msg :=
		"Story with id " + fmt.Sprintf("%d", storyID) +
			" has already been backed by user " + creator.String()

	return sdk.NewError(DefaultCodespace, CodeDuplicate, msg)
}
