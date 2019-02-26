package params

import (
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/vote"
	"github.com/TruStory/truchain/x/voting"
)

// MinBackingAmount : '1000000000',
// MaxBackingAmount:  '100000000000',

// Params defines defaults for a story
type Params struct {
	ChallengeParams  challenge.Params
	ExpirationParams expiration.Params
	StakeParams      stake.Params
	StoryParams      story.Params
	VoteParams       vote.Params
	VotingParams     voting.Params
}

// DefaultParams creates the default params
func DefaultParams() Params {
	return Params{
		ChallengeParams:  challenge.DefaultParams(),
		ExpirationParams: expiration.DefaultParams(),
		StakeParams:      stake.DefaultParams(),
		StoryParams:      story.DefaultParams(),
		VoteParams:       vote.DefaultParams(),
		VotingParams:     voting.DefaultParams(),
	}
}
