package trubank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TruBank errors reserve 1400 ~ 1499.
const (
	DefaultCodespace sdk.CodespaceType = "tbank"

	CodeErrorAddingCoinsToUser     sdk.CodeType = 1401
	CodeErrorAddingCoinsToCategory sdk.CodeType = 1402
	CodeErrorTransactionNotFound   sdk.CodeType = 1403
)

// ErrTransferringCoinsToUser throws an error when the category is invalid
func ErrTransferringCoinsToUser(creator sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeErrorAddingCoinsToUser, "Coins could not be added to the user "+creator.String())
}

// ErrTransferringCoinsToCategory throws an error when a category msg is invalid
func ErrTransferringCoinsToCategory(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeErrorAddingCoinsToCategory, "Coins could not be added to category "+fmt.Sprintf("%d", id))
}

// ErrTransactionNotFound throws an error when a transaction was not found
func ErrTransactionNotFound(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeErrorTransactionNotFound, "There was no transaction found with an id of "+fmt.Sprintf("%d", id))
}
