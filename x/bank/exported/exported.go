package exported

import (
	"time"

	"github.com/TruStory/truchain/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transaction stores data related to a transaction
type Transaction struct {
	ID                uint64          `json:"id"`
	Type              TransactionType `json:"type"`
	AppAccountAddress sdk.AccAddress  `json:"app_account_address"`
	ReferenceID       uint64          `json:"reference_id"`
	CommunityID       string          `json:"community_id"`
	Amount            sdk.Coin        `json:"amount"`
	CreatedTime       time.Time       `json:"created_time"`
}

// TransactionType defines the type of transaction.
type TransactionType int8

// Types of transactions
const (
	TransactionGift TransactionType = iota
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
	TransactionInterestArgumentCreationSlashed
	TransactionInterestUpvoteReceivedSlashed
	TransactionInterestUpvoteGivenSlashed
	TransactionStakeCreatorSlashed
	TransactionStakeCuratorSlashed
	TransactionCuratorReward
)

var TransactionTypeName = []string{
	TransactionGift:                            "TransactionGift",
	TransactionBacking:                         "TransactionBacking",
	TransactionBackingReturned:                 "TransactionBackingReturned",
	TransactionChallenge:                       "TransactionChallenge",
	TransactionChallengeReturned:               "TransactionChallengeReturned",
	TransactionUpvote:                          "TransactionUpvote",
	TransactionUpvoteReturned:                  "TransactionUpvoteReturned",
	TransactionInterestArgumentCreation:        "TransactionInterestArgumentCreation",
	TransactionInterestUpvoteReceived:          "TransactionInterestUpvoteReceived",
	TransactionInterestUpvoteGiven:             "TransactionInterestUpvoteGiven",
	TransactionRewardPayout:                    "TransactionRewardPayout",
	TransactionInterestArgumentCreationSlashed: "TransactionInterestArgumentCreationSlashed",
	TransactionInterestUpvoteReceivedSlashed:   "TransactionInterestUpvoteReceivedSlashed",
	TransactionInterestUpvoteGivenSlashed:      "TransactionInterestUpvoteGivenSlashed",
	TransactionStakeCreatorSlashed:             "TransactionStakeCreatorSlashed",
	TransactionStakeCuratorSlashed:             "TransactionStakeCuratorSlashed",
}

func (t TransactionType) String() string {
	if int(t) >= len(TransactionTypeName) {
		return "Unknown"
	}
	return TransactionTypeName[t]
}

var AllowedTransactionsForAddition = []TransactionType{
	TransactionGift,
	TransactionBackingReturned,
	TransactionChallengeReturned,
	TransactionUpvoteReturned,
	TransactionInterestArgumentCreation,
	TransactionInterestUpvoteReceived,
	TransactionInterestUpvoteGiven,
	TransactionRewardPayout,
	TransactionCuratorReward,
}

var AllowedTransactionsForEarning = []TransactionType{
	TransactionInterestArgumentCreation,
	TransactionInterestUpvoteReceived,
	TransactionInterestUpvoteGiven,
	TransactionCuratorReward,
}

var AllowedTransactionsForDeduction = []TransactionType{
	TransactionBacking,
	TransactionChallenge,
	TransactionUpvote,
	TransactionInterestArgumentCreationSlashed,
	TransactionInterestUpvoteReceivedSlashed,
	TransactionInterestUpvoteGivenSlashed,
	TransactionStakeCreatorSlashed,
	TransactionStakeCuratorSlashed,
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

type TransactionSetter func(*Transaction)

func WithCommunityID(communityID string) TransactionSetter {
	return func(tx *Transaction) {
		tx.CommunityID = communityID
	}
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

// Defines bank module constants
const (
	QueryTransactionsByAddress = "transactions_by_address"
	QueryParams                = "params"
	ModuleName                 = types.ModuleName
	StoreKey                   = ModuleName
	RouterKey                  = ModuleName
	QuerierRoute               = ModuleName
	DefaultParamspace          = ModuleName
)

// QueryTransactionsByAddress query transactions params for a specific address.
type QueryTransactionsByAddressParams struct {
	Address   sdk.AccAddress    `json:"address"`
	Types     []TransactionType `json:"types,omitempty"`
	SortOrder SortOrderType     `json:"sort_order,omitempty"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
}
