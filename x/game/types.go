package game

import (
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Game defines a validation game on a story
type Game struct {
	ID              int64          `json:"id"`
	StoryID         int64          `json:"story_id"`
	Creator         sdk.AccAddress `json:"creator"`
	ExpiresTime     time.Time      `json:"expires_time,omitempty"`
	EndTime         time.Time      `json:"end_time,omitempty"`
	Pool            sdk.Coin       `json:"pool"`
	Started         bool           `json:"started"`
	ThresholdAmount sdk.Int        `json:"threshold_amount"`
	Timestamp       app.Timestamp  `json:"timestamp"`
}

// Params holds default parameters for a game
type Params struct {
	MinChallengeStake sdk.Int       // min amount required to challenge
	Expires           time.Duration // time to expire if threshold not met
	Period            time.Duration // length of challenge game / voting period
	Threshold         int64         // amount at which game begins
}

// DefaultParams creates a new Params type with defaults
func DefaultParams() Params {
	return Params{
		MinChallengeStake: sdk.NewInt(10),
		Expires:           10 * 24 * time.Hour,
		Period:            30 * 24 * time.Hour,
		Threshold:         10,
	}
}
