package bank

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Define keys
var (
	TransactionsKeyPrefix = []byte{0x00}

	// ID Keys
	TransactionIDKey = []byte{0x10}

	// AssociationKeys
	UserTransactionKeyPrefix = []byte{0x20}
)

// stakeKey gets a key for a stake.
// 0x00<transaction_id>
func transactionKey(id uint64) []byte {
	bz := sdk.Uint64ToBigEndian(id)
	return append(TransactionsKeyPrefix, bz...)
}

// userTransactionsPrefix
// 0x20<creator>
func userTransactionsPrefix(creator sdk.AccAddress) []byte {
	return append(UserTransactionKeyPrefix, creator.Bytes()...)
}

// userTransactionKey builds the key for user->transaction association
// 0x20<creator><created_time><transaction_id>
func userTransactionKey(creator sdk.AccAddress, createdTime time.Time, transactionID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(transactionID)
	timeBz := sdk.FormatTimeBytes(createdTime)
	return append(userTransactionsPrefix(creator), append(timeBz, bz...)...)
}
