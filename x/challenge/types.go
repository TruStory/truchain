package challenge

import (
	"net/url"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenger defines a user's challenge on a story
type Challenger struct {
	Amount    sdk.Coin       `json:"amount"`
	Argument  string         `json:"argument"`
	Creator   sdk.AccAddress `json:"creator"`
	Evidence  []url.URL      `json:"evidence,omitempty"`
	Timestamp app.Timestamp  `json:"timestamp"`
}

// Challenge defines a challenge on a story
type Challenge struct {
	ID              int64          `json:"id"`
	StoryID         int64          `json:"story_id"`
	Creator         sdk.AccAddress `json:"creator"`
	ExpiresTime     time.Time      `json:"expires_time,omitempty"`
	Pool            sdk.Coin       `json:"pool"`
	Started         bool           `json:"started"`
	ThresholdAmount sdk.Int        `json:"threshold_amount"`
	Timestamp       app.Timestamp  `json:"timestamp"`
}

// Params holds default parameters for a challenge
type Params struct {
	MinArgumentLength int           // min number of chars for argument
	MaxArgumentLength int           // max number of chars for argument
	MinEvidenceCount  int           // min number of evidence URLs
	MaxEvidenceCount  int           // max number of evidence URLs
	MinChallengeStake sdk.Int       // min amount required to challenge
	Expires           time.Duration // time to expire if threshold not met
	Period            time.Duration // length of challenge game / voting period
	Threshold         int64         // amount at which game begins
}

// DefaultParams creates a new Params type with defaults
func DefaultParams() Params {
	return Params{
		MinArgumentLength: 10,
		MaxArgumentLength: 340,
		MinEvidenceCount:  0,
		MaxEvidenceCount:  10,
		MinChallengeStake: sdk.NewInt(10),
		Expires:           10 * 24 * time.Hour,
		Period:            30 * 24 * time.Hour,
		Threshold:         10,
	}
}
