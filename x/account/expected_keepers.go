package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransactionType defines the type of transaction.
// NOTE: This needs to stay in sync with x/bank/types.go
type TransactionType int8

// Types of transactions
const (
	TransactionRegistration TransactionType = iota
	TransactionBacking
	TransactionBackingReturned
	TransactionChallenge
	TransactionChallengeReturned
	TransactionUpvote
	TransactionUpvoteReturned
	TransactionInterest
	TransactionRewardPayout
)

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType TransactionType) (sdk.Coins, sdk.Error)
}
