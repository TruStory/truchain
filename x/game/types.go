package game

import (
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// game states
// started
// * challenge threshold met
// * voting begins
// ended
// * voting period ended (24 hrs)
// * AND received min quorum (9+ votes)
// expired
// * voting period ended (24 hrs)
// * NOT received min quorum (9+ votes)
// * stake returned

// Game defines a validation game on a story
type Game struct {
	ID            int64          `json:"id"`
	StoryID       int64          `json:"story_id"`
	Creator       sdk.AccAddress `json:"creator"`
	ExpiresTime   time.Time      `json:"expires_time,omitempty"`
	EndTime       time.Time      `json:"end_time,omitempty"`
	ChallengePool sdk.Coin       `json:"challenge_pool,omitempty"`
	Started       bool           `json:"started,omitempty"`
	VoteQuorum    int64          `json:"vote_quorum,omitempty"`
	Timestamp     app.Timestamp  `json:"timestamp"`
}

// Ended returns true if game is over time and quorum is reached
func (g Game) Ended(time time.Time) bool {
	if time.After(g.EndTime) &&
		g.VoteQuorum >= DefaultParams().VoteQuorum {

		return true
	}

	return false
}

// Expired returns true if game is over and quorum is not reached
func (g Game) Expired(time time.Time) bool {
	if time.After(g.EndTime) &&
		g.VoteQuorum < DefaultParams().VoteQuorum {

		return true
	}

	return false
}

// Params holds default parameters for a game
type Params struct {
	ChallengeThreshold sdk.Dec       // % backings at which game begins
	MinChallengeStake  sdk.Int       // min amount required to challenge
	Expires            time.Duration // time to expire if threshold not met
	VotingPeriod       time.Duration // length of challenge game / voting period
	VoteQuorum         int64         // num voters required
}

// DefaultParams creates a new MsgParams type with defaults
func DefaultParams() Params {
	return Params{
		ChallengeThreshold: sdk.NewDecWithPrec(33, 2), // 33%
		MinChallengeStake:  sdk.NewInt(10),
		Expires:            10 * 24 * time.Hour,
		VotingPeriod:       1 * 24 * time.Hour,
		VoteQuorum:         7,
	}
}
