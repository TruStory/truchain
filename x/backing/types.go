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
	MinInterestRate sdk.Dec
	MaxInterestRate sdk.Dec
}

// DefaultParams creates a new Params type with defaults
func DefaultParams() Params {
	return Params{
		AmountWeight:    sdk.NewDecWithPrec(333, 3), // 33.3%
		PeriodWeight:    sdk.NewDecWithPrec(667, 3), // 66.7%
		MinInterestRate: sdk.ZeroDec(),              // 0%
		MaxInterestRate: sdk.NewDecWithPrec(10, 2),  // 10%
	}
}

// MsgParams holds default validation params
type MsgParams struct {
	MinPeriod time.Duration
	MaxPeriod time.Duration
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinPeriod: 3 * 24 * time.Hour,  // 3 days
		MaxPeriod: 90 * 24 * time.Hour, // 90 days
	}
}

// Backing type
type Backing struct {
	app.Vote

	StoryID  int64         `json:"story_id"`
	Interest sdk.Coin      `json:"interest"`
	Expires  time.Time     `json:"expires"`
	Params   Params        `json:"params"`
	Period   time.Duration `json:"period"`
}

// AmountDenom implements `Voter`
func (b Backing) AmountDenom() string {
	return b.Amount.Denom
}
