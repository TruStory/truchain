package bank

import (
	"github.com/TruStory/truchain/x/bank/exported"
	"github.com/TruStory/truchain/x/bank/types"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace

	TransactionRegistration             = exported.TransactionRegistration
	TransactionBacking                  = exported.TransactionBacking
	TransactionBackingReturned          = exported.TransactionBackingReturned
	TransactionChallenge                = exported.TransactionChallenge
	TransactionChallengeReturned        = exported.TransactionChallengeReturned
	TransactionUpvote                   = exported.TransactionUpvote
	TransactionUpvoteReturned           = exported.TransactionUpvoteReturned
	TransactionInterestArgumentCreation = exported.TransactionInterestArgumentCreation
	TransactionInterestUpvoteReceived   = exported.TransactionInterestUpvoteReceived
	TransactionInterestUpvoteGiven      = exported.TransactionInterestUpvoteGiven
	TransactionRewardPayout             = exported.TransactionRewardPayout
	SortAsc                             = exported.SortAsc
	SortDesc                            = exported.SortDesc
	QueryTransactionsByAddress          = exported.QueryTransactionsByAddress
	RouterKey                           = exported.RouterKey
)

var (
	AccountKey              = types.AccountKey
	GetFilters              = exported.GetFilters
	FilterByTransactionType = exported.FilterByTransactionType
	SortOrder               = exported.SortOrder
	Limit                   = exported.Limit
	Offset                  = exported.Offset
)

type (
	// RecordKeeper alias
	RecordKeeper                     = recordkeeper.RecordKeeper
	TransactionType                  = exported.TransactionType
	Filter                           = exported.Filter
	SortOrderType                    = exported.SortOrderType
	Transaction                      = exported.Transaction
	QueryTransactionsByAddressParams = exported.QueryTransactionsByAddressParams
)
