package params

import (
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
)

// TODO [shanev]: these will be added by https://github.com/TruStory/truchain/issues/399
// MinBackingAmount : '1000000000',
// MaxBackingAmount:  '100000000000',

// Params defines defaults for a story
type Params struct {
	ArgumentParams   argument.Params
	ChallengeParams  challenge.Params
	ExpirationParams expiration.Params
	StakeParams      stake.Params
	StoryParams      story.Params
}

// DefaultParams creates the default params
func DefaultParams() Params {
	return Params{
		ArgumentParams:   argument.DefaultParams(),
		ChallengeParams:  challenge.DefaultParams(),
		ExpirationParams: expiration.DefaultParams(),
		StakeParams:      stake.DefaultParams(),
		StoryParams:      story.DefaultParams(),
	}
}
