package staking

import "github.com/TruStory/truchain/x/bank/exported"

const (
	TransactionInterestArgumentCreation = exported.TransactionInterestArgumentCreation
	TransactionInterestUpvoteReceived   = exported.TransactionInterestUpvoteReceived
	TransactionInterestUpvoteGiven      = exported.TransactionInterestUpvoteGiven
	TransactionBacking                  = exported.TransactionBacking
	TransactionChallenge                = exported.TransactionChallenge
	TransactionUpvote                   = exported.TransactionUpvote
	TransactionBackingReturned          = exported.TransactionBackingReturned
	TransactionChallengeReturned        = exported.TransactionChallengeReturned
	TransactionUpvoteReturned           = exported.TransactionUpvoteReturned
)

type (
	TransactionType = exported.TransactionType
)
