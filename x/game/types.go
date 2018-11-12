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
	ThresholdAmount sdk.Int        `json:"threshold_amount"`
	Timestamp       app.Timestamp  `json:"timestamp"`
}

// Started returns true if a validation game has started
func (g Game) Started() bool {
	if g.EndTime.After(time.Time{}) {
		return true
	}

	return false
}

// MsgParams holds default parameters for a game
type MsgParams struct {
	MinChallengeStake sdk.Int       // min amount required to challenge
	Expires           time.Duration // time to expire if threshold not met
	Period            time.Duration // length of challenge game / voting period
	Threshold         int64         // amount at which game begins
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinChallengeStake: sdk.NewInt(10),
		Expires:           10 * 24 * time.Hour,
		Period:            30 * 24 * time.Hour,
		Threshold:         10,
	}
}
