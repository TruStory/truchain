package staking

import (
	"github.com/TruStory/truchain/x/bank/exported"
	"github.com/TruStory/truchain/x/distribution"
)

// Aliased constants
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

	UserRewardPoolName = distribution.UserRewardPoolName
)

type (
	TransactionType = exported.TransactionType
)

// Transaction setters
var (
	WithCommunityID   = exported.WithCommunityID
	FromModuleAccount = exported.FromModuleAccount
	ToModuleAccount   = exported.ToModuleAccount
)
