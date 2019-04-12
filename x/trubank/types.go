package trubank

import (
	app "github.com/TruStory/truchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transaction struct stores the transactions for the trubank
type Transaction struct {
	ID              int64           `json:"id"`
	TransactionType TransactionType `json:"transaction_type"`
	GroupID         int64           `json:"group_id"`
	ReferenceID     int64           `json:"reference_id,omitempty"`
	Amount          sdk.Coin        `json:"amount"`
	Creator         sdk.AccAddress  `json:"creator"`
	Timestamp       app.Timestamp   `json:"timestamp"`
}

// TransactionType defines the type of transaction
type TransactionType int8

// List of transaction types
const (
	Backing TransactionType = iota
	Challenge
	BackingReturned
	ChallengeReturned
	RewardPool
	Interest
	// Deprecated
	Like
	BackingLike
	ChallengeLike
)
