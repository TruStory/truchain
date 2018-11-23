package game

import (
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// game states
// created
// * at least one challenge
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
	VoteQuorum    int64          `json:"vote_quorum,omitempty"`
	Timestamp     app.Timestamp  `json:"timestamp"`
}

// Started returns true if challenge pool meets threshold
func (g Game) Started() bool {
	return g.ChallengePool.Amount.GT(DefaultParams().ChallengeThreshold)
}

// Ended returns true if game is over time and quorum is reached
func (g Game) Ended(time time.Time) bool {
	if time.After(g.EndTime) &&
		g.VoteQuorum >= DefaultParams().VoteQuorum {

		return true
	}

	return false
}

// Params holds default parameters for a game
type Params struct {
	MinChallengeStake  sdk.Int       // min amount required to challenge
	Expires            time.Duration // time to expire if threshold not met
	VotingPeriod       time.Duration // length of challenge game / voting period
	ChallengeThreshold sdk.Int       // amount at which game begins
	VoteQuorum         int64         // num voters required
}

// DefaultParams creates a new MsgParams type with defaults
func DefaultParams() Params {
	return Params{
		MinChallengeStake:  sdk.NewInt(10),
		Expires:            10 * 24 * time.Hour,
		VotingPeriod:       1 * 24 * time.Hour,
		ChallengeThreshold: sdk.NewInt(10),
		VoteQuorum:         7,
	}
}
