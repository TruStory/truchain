package staking

import (
	"github.com/TruStory/truchain/x/bank/exported"
	"github.com/TruStory/truchain/x/distribution"
)

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

var (
	WithCommunityID = exported.WithCommunityID
)
