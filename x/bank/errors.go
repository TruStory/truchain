package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 400 ~ 499.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidTransactionType     sdk.CodeType = 401
	ErrorCodeInvalidRewardBrokerAddress sdk.CodeType = 402
	ErrorCodeInvalidQueryParams         sdk.CodeType = 403
	ErrorCodeUnknownTransaction         sdk.CodeType = 404
)

// ErrInvalidRewardBrokerAddress throws an error when the address doesn't match with genesis param address.
func ErrInvalidRewardBrokerAddress(address sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidRewardBrokerAddress,
		fmt.Sprintf("Invalid reward broker address %s", address.String()),
	)
}

// ErrInvalidTransactionType throws an error when the transaction type is invalid.
func ErrInvalidTransactionType(txType TransactionType) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidTransactionType,
		fmt.Sprintf("Invalid transaction type %s", txType.String()),
	)
}

// ErrInvalidQueryParams throws an error when the transaction type is invalid.
func ErrInvalidQueryParams(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidQueryParams,
		fmt.Sprintf("Invalid query params  %s", err.Error()),
	)
}

// ErrCodeUnknownTransaction throws an error when an invalid transaction id
func ErrCodeUnknownTransaction(transactionID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeUnknownTransaction,
		fmt.Sprintf("Unknown transaction id %d", transactionID),
	)
}
