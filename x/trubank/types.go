package trubank

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transaction struct stores the transactions for the trubank
type Transaction struct {
	ID              int64           `json:"id"`
	TransactionType TransactionType `json:"transaction_type"`
	ReferenceID     int64           `json:"reference_id"`
	Creator         sdk.AccAddress  `json:"creator"`
	Status          Status          `json:"status"`
	Timestamp       app.Timestamp   `json:"timestamp"`
}

// TransactionDetails pulls back both the transaction and the vote object around it
type TransactionDetails struct {
	Transaction
	Details *app.Vote
}

// TransactionType defines the type of transaction
type TransactionType int8

// List of transaction types
const (
	Backing TransactionType = iota
	Challenge
	// StoryCreation
)

// Status defines the status of the transaction
type Status int8

// List of the allowed statuses for a transaction
const (
	Pending Status = iota
	Completed
)
