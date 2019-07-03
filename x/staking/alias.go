package staking

import "github.com/TruStory/truchain/x/bank/exported"

const (
	TransactionInterestArgumentCreation = exported.TransactionInterestArgumentCreation
	TransactionInterestUpvoteReceived = exported.TransactionInterestUpvoteReceived
	TransactionInterestUpvoteGiven = exported.TransactionInterestUpvoteGiven
	TransactionBacking = exported.TransactionBacking
	TransactionChallenge = exported.TransactionChallenge
	TransactionUpvote = exported.TransactionUpvote
)

type (
	TransactionType = exported.TransactionType
)