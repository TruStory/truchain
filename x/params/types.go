package params

import (
	"time"

	"github.com/TruStory/truchain/x/backing"
)

// Params defines defaults for a story
type Params struct {
	BackingAmountWeight    string        `json:"backing_amount_weight,omitempty"`
	BackingPeriodWeight    string        `json:"backing_period_weight,omitempty"`
	BackingMinInterestRate string        `json:"backing_min_intereset_rate,omitempty"`
	BackingMaxInterestRate string        `json:"backing_max_intereset_rate,omitempty"`
	BackingMinDuration     time.Duration `json:"backing_min_duration,omitempty"`
	BackingMaxDuration     time.Duration `json:"backing_max_duration,omitempty"`
}

// DefaultParams creates the default params
func DefaultParams() Params {

	return Params{
		BackingAmountWeight:    backing.DefaultParams().AmountWeight.String(),
		BackingPeriodWeight:    backing.DefaultParams().PeriodWeight.String(),
		BackingMinInterestRate: backing.DefaultParams().MinInterestRate.String(),
		BackingMaxInterestRate: backing.DefaultParams().MaxInterestRate.String(),
		BackingMinDuration:     backing.DefaultMsgParams().MinPeriod,
		BackingMaxDuration:     backing.DefaultMsgParams().MaxPeriod,
	}
}
