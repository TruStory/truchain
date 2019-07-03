package bank

import (
	"github.com/TruStory/truchain/x/bank/exported"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
)

const (
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
	ModuleName                          = exported.ModuleName
	StoreKey                            = exported.StoreKey
	RouterKey                           = exported.RouterKey
	QuerierRoute                        = exported.QuerierRoute
	DefaultParamspace                   = exported.DefaultParamspace
)

var (
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
