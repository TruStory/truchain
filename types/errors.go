package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// App errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = 1

	CodeInvalidCommentMsg  sdk.CodeType = 101
	CodeInvalidEvidenceMsg sdk.CodeType = 102
)

// ErrInvalidCommentMsg creates an error when `Msg` validation fails
func ErrInvalidCommentMsg() sdk.Error {
	return sdk.NewError(
		DefaultCodespace, CodeInvalidCommentMsg, "Invalid comment")
}

// ErrInvalidEvidenceMsg creates an error when `Msg` validation fails
func ErrInvalidEvidenceMsg() sdk.Error {
	return sdk.NewError(
		DefaultCodespace, CodeInvalidEvidenceMsg, "Invalid evidence")
}
