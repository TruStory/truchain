package bank

import (
	"github.com/TruStory/truchain/x/bank/exported"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
)

const (
	TransactionRewardPayout = exported.TransactionRewardPayout
	TransactionBacking = exported.TransactionBacking
	TransactionChallenge = exported.TransactionChallenge
	TransactionUpvote = exported.TransactionUpvote
	TransactionRegistration = exported.TransactionRegistration
	TransactionBackingReturned = exported.TransactionBackingReturned
	TransactionUpvoteReturned = exported.TransactionUpvoteReturned
	SortAsc = exported.SortAsc
	SortDesc = exported.SortDesc
)

var (
	GetFilters = exported.GetFilters
	FilterByTransactionType = exported.FilterByTransactionType
	SortOrder = exported.SortOrder
	Limit = exported.Limit
	Offset = exported.Offset
)

type (
	// RecordKeeper alias
	RecordKeeper = recordkeeper.RecordKeeper
	TransactionType = exported.TransactionType
	Filter = exported.Filter
	SortOrderType = exported.SortOrderType
	Transaction = exported.Transaction
)
