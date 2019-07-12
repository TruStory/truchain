package params

import (
	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/slashing"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/staking"
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
	AccountParams    account.Params
	CommunityParams  community.Params
	ClaimParams      claim.Params
	BankParams       bank.Params
	StakingParams    staking.Params
	SlashingParams   slashing.Params
}

// DefaultParams creates the default params
func DefaultParams() Params {
	return Params{
		ArgumentParams:   argument.DefaultParams(),
		ChallengeParams:  challenge.DefaultParams(),
		ExpirationParams: expiration.DefaultParams(),
		StakeParams:      stake.DefaultParams(),
		StoryParams:      story.DefaultParams(),
		AccountParams:    account.DefaultParams(),
		CommunityParams:  community.DefaultParams(),
		ClaimParams:      claim.DefaultParams(),
		BankParams:       bank.DefaultParams(),
		StakingParams:    staking.DefaultParams(),
		SlashingParams:   slashing.DefaultParams(),
	}
}
