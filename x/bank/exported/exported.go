package exported

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transaction stores data related to a transaction
type Transaction struct {
	ID                uint64
	Type              TransactionType
	AppAccountAddress sdk.AccAddress
	ReferenceID       uint64
	Amount            sdk.Coin
	CreatedTime       time.Time
}

// TransactionType defines the type of transaction.
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
	TransactionInterestArgumentCreation
	TransactionInterestUpvoteReceived
	TransactionInterestUpvoteGiven
	TransactionRewardPayout
)

var TransactionTypeName = []string{
	TransactionRegistration:             "TransactionRegistration",
	TransactionBacking:                  "TransactionBacking",
	TransactionBackingReturned:          "TransactionBackingReturned",
	TransactionChallenge:                "TransactionChallenge",
	TransactionChallengeReturned:        "TransactionChallengeReturned",
	TransactionUpvote:                   "TransactionUpvote",
	TransactionUpvoteReturned:           "TransactionUpvoteReturned",
	TransactionInterestArgumentCreation: "TransactionInterestArgumentCreation",
	TransactionInterestUpvoteReceived:   "TransactionInterestUpvoteReceived",
	TransactionInterestUpvoteGiven:      "TransactionInterestUpvoteGiven",
	TransactionRewardPayout:             "TransactionRewardPayout",
}

func (t TransactionType) String() string {
	if int(t) >= len(TransactionTypeName) {
		return "Unknown"
	}
	return TransactionTypeName[t]
}

var AllowedTransactionsForAddition = []TransactionType{
	TransactionRegistration,
	TransactionBackingReturned,
	TransactionChallengeReturned,
	TransactionUpvoteReturned,
	TransactionInterestArgumentCreation,
	TransactionInterestUpvoteReceived,
	TransactionInterestUpvoteGiven,
	TransactionRewardPayout,
}

var AllowedTransactionsForDeduction = []TransactionType{
	TransactionBacking,
	TransactionChallenge,
	TransactionUpvote,
}

func (t TransactionType) AllowedForAddition() bool {
	return t.OneOf(AllowedTransactionsForAddition)
}

func (t TransactionType) AllowedForDeduction() bool {
	return t.OneOf(AllowedTransactionsForDeduction)
}

func (t TransactionType) OneOf(types []TransactionType) bool {
	for _, tType := range types {
		if tType == t {
			return true
		}
	}
	return false
}

type SortOrderType int8

const (
	SortAsc SortOrderType = iota
	SortDesc
)

func (t SortOrderType) Valid() bool {
	return t == SortAsc || t == SortDesc
}

type Filters struct {
	TransactionTypes []TransactionType
	SortOrder        SortOrderType
	Limit            int
	Offset           int
}

type Filter func(*Filters)

func FilterByTransactionType(transactionTypes ...TransactionType) Filter {
	return func(filters *Filters) {
		filters.TransactionTypes = transactionTypes
	}
}

func SortOrder(sortOrder SortOrderType) Filter {
	return func(filters *Filters) {
		if !sortOrder.Valid() {
			return
		}
		filters.SortOrder = sortOrder
	}
}

func Limit(limit int) Filter {
	return func(filters *Filters) {
		if limit <= 0 {
			return
		}
		filters.Limit = limit
	}
}

func Offset(offset int) Filter {
	return func(filters *Filters) {
		if offset <= 0 {
			return
		}
		filters.Offset = offset
	}
}

func GetFilters(filterSetters ...Filter) Filters {
	filters := Filters{
		TransactionTypes: make([]TransactionType, 0),
		SortOrder:        SortAsc,
		Limit:            0,
		Offset:           0,
	}
	for _, filter := range filterSetters {
		filter(&filters)
	}
	return filters
}
