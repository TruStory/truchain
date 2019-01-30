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
	MinArgumentLength int // min number of chars for argument
	MaxArgumentLength int // max number of chars for argument
	MinPeriod         time.Duration
	MaxPeriod         time.Duration
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinArgumentLength: 10,
		MaxArgumentLength: 1000,
		// MinPeriod:         3 * 24 * time.Hour,  // 3 days
		// MaxPeriod:         90 * 24 * time.Hour, // 90 days
		MinPeriod: 1 * time.Hour,
		MaxPeriod: 3 * 24 * time.Hour,
	}
}

// Backing type
type Backing struct {
	app.Vote `json:"vote"`

	StoryID     int64         `json:"story_id"`
	Interest    sdk.Coin      `json:"interest"`
	MaturesTime time.Time     `json:"matures_time"`
	Params      Params        `json:"params"`
	Period      time.Duration `json:"period"`
}

// ID implements `Voter`
func (b Backing) ID() int64 {
	return b.Vote.ID
}

// Amount implements `Voter`
func (b Backing) Amount() sdk.Coin {
	return b.Vote.Amount
}

// Creator implements `Voter`
func (b Backing) Creator() sdk.AccAddress {
	return b.Vote.Creator
}

// VoteChoice implements `Voter`
func (b Backing) VoteChoice() bool {
	return b.Vote.Vote
}

// HasMatured is true if a backing has matured
func (b Backing) HasMatured(blockTime time.Time) bool {
	return blockTime.After(b.MaturesTime)
}
