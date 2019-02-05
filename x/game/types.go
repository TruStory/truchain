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
// finished
// * voting period ended (24 hrs)
// * AND received min quorum (9+ votes)
// expired
// * voting period ended (24 hrs)
// * NOT received min quorum (9+ votes)
// * stake returned

// // PerID contains a game and its gameID
// type PerID struct {
// 	GameID int64 `json:"game_id"`
// 	Game   Game  `json:"game"`
// }

// Game defines a validation game on a story
type Game struct {
	ID                  int64          `json:"id"`
	StoryID             int64          `json:"story_id"`
	Creator             sdk.AccAddress `json:"creator"`
	ChallengePool       sdk.Coin       `json:"challenge_pool,omitempty"`
	ChallengeExpireTime time.Time      `json:"challenge_expire_time,omitempty"`
	VotingEndTime       time.Time      `json:"voting_end_time,omitempty"`
	Timestamp           app.Timestamp  `json:"timestamp"`
}

// IsExpired if challenge threshold not met in a certain time
func (g Game) IsExpired(blockTime time.Time) bool {
	return blockTime.After(g.ChallengeExpireTime)
}

// IsVotingExpired returns true if:
// 1. passed the voting period (`VotingPeriodEndTime` > block time)
// 2. didn't meet the minimum voter quorum
func (g Game) IsVotingExpired(blockTime time.Time, quorum int) bool {
	return blockTime.After(g.VotingEndTime) &&
		(quorum < DefaultParams().VoteQuorum)
}

// IsVotingFinished returns true if:
// 1. passed the voting period (`VotingPeriodEndTime` > block time)
// 2. met the minimum voter quorum
func (g Game) IsVotingFinished(blockTime time.Time, quorum int) bool {
	return blockTime.After(g.VotingEndTime) &&
		(quorum >= DefaultParams().VoteQuorum)
}

// Params holds default parameters for a game
type Params struct {
	ChallengeToBackingRatio sdk.Dec       // % backings at which game begins
	MinChallengeThreshold   sdk.Int       // min amount required to start a game
	MinChallengeStake       sdk.Int       // min amount required to join a challenge
	Expires                 time.Duration // time to expire if threshold not met
	VotingPeriod            time.Duration // length of challenge game / voting period
	VoteQuorum              int           // num voters (BCV) required
}

// DefaultParams creates a new MsgParams type with defaults
func DefaultParams() Params {
	return Params{
		ChallengeToBackingRatio: sdk.NewDecWithPrec(33, 2), // 33%
		MinChallengeThreshold:   sdk.NewInt(10000000000),   // 10 trustake
		MinChallengeStake:       sdk.NewInt(1000000000),    //  1 trustake
		// Expires:                 10 * 24 * time.Hour,
		// VotingPeriod:            1 * 24 * time.Hour,
		// VoteQuorum:              7,
		Expires:      1 * 24 * time.Hour,
		VotingPeriod: 3 * time.Hour,
		VoteQuorum:   3,
	}
}
