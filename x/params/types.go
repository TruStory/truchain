package params

import (
	"time"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/voting"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinArgumentLength: 10,
// MaxArgumentLength: 1000,
// MaxInterestRate: .10,
// AmountWeight: .333,
// DurationWeight: .667,

// MinBackingAmount : '1000000000',
// MaxBackingAmount:  '100000000000',
// AddStoryStake:     '10000000000',
// VoteStake:         '10000000000',
// BackingSupply:     '100000000000000',

// Params defines defaults for a story
type Params struct {
	StoryParams      story.Params
	ChallengeParams  challenge.Params
	ExpirationParams expiration.Params
	VotingParams     voting.Params

	BackingAmountWeight    sdk.Dec       `json:"backing_amount_weight,omitempty"`
	BackingPeriodWeight    sdk.Dec       `json:"backing_period_weight,omitempty"`
	BackingMinInterestRate sdk.Dec       `json:"backing_min_intereset_rate,omitempty"`
	BackingMaxInterestRate sdk.Dec       `json:"backing_max_intereset_rate,omitempty"`
	BackingMinDuration     time.Duration `json:"backing_min_duration,omitempty"`
	BackingMaxDuration     time.Duration `json:"backing_max_duration,omitempty"`
}

// DefaultParams creates the default params
func DefaultParams() Params {

	return Params{
		BackingAmountWeight:    backing.DefaultParams().AmountWeight,
		BackingPeriodWeight:    backing.DefaultParams().PeriodWeight,
		BackingMinInterestRate: backing.DefaultParams().MinInterestRate,
		BackingMaxInterestRate: backing.DefaultParams().MaxInterestRate,
		BackingMinDuration:     backing.DefaultMsgParams().MinPeriod,
		BackingMaxDuration:     backing.DefaultMsgParams().MaxPeriod,
	}
}
