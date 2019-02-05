package params

import (
	"time"

	"github.com/TruStory/truchain/x/backing"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params defines defaults for a story
type Params struct {
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
