package backing

import (
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================

// Queue is a list of all backings
type Queue []int64

// IsEmpty checks if the queue is empty
func (asq Queue) IsEmpty() bool {
	if len(asq) == 0 {
		return true
	}
	return false
}

// ============================================================================

// Params holds data for backing interest calculations
type Params struct {
	AmountWeight    sdk.Dec
	PeriodWeight    sdk.Dec
	MinPeriod       time.Duration
	MaxPeriod       time.Duration
	MinInterestRate sdk.Dec
	MaxInterestRate sdk.Dec
}

// DefaultParams creates a new BackingParams type with defaults
func DefaultParams() Params {
	return Params{
		AmountWeight:    sdk.NewDecWithPrec(333, 3), // 33.3%
		PeriodWeight:    sdk.NewDecWithPrec(667, 3), // 66.7%
		MinPeriod:       3 * 24 * time.Hour,         // 3 days
		MaxPeriod:       90 * 24 * time.Hour,        // 90 days
		MinInterestRate: sdk.ZeroDec(),              // 0%
		MaxInterestRate: sdk.NewDecWithPrec(10, 2),  // 10%
	}
}

// Backing type
type Backing struct {
	ID        int64          `json:"id"`
	StoryID   int64          `json:"story_id"`
	Principal sdk.Coin       `json:"principal"`
	Interest  sdk.Coin       `json:"interest"`
	Expires   time.Time      `json:"expires"`
	Params    Params         `json:"params"`
	Period    time.Duration  `json:"period"`
	User      sdk.AccAddress `json:"user"`
	Timestamp app.Timestamp  `json:"timestamp"`
}
