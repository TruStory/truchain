package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ChallengesByID holds a Challenge struct and its associated Challenge ID
type ChallengesByID struct {
	ChallengeID int64       `json:"challenge_id"`
	Challenge   []Challenge `json:"challenge"`
}

// Challengers struct contains a creator address asosociated with its challengeID
type Challengers struct {
	ChallengeID int64          `json:"challenge_id"`
	AccAddress  sdk.AccAddress `json:"creator"`
}

// ChallengersPerGame contains an array of challengers per GameID
type ChallengersPerGame struct {
	GameID      int64         `json:"game_id"`
	Challengers []Challengers `json:"challengers"`
}

// Challenge defines a user's challenge on a story
type Challenge struct {
	app.Vote `json:"vote"`
}

// ID implements `Voter`
func (c Challenge) ID() int64 {
	return c.Vote.ID
}

// Amount implements `Voter`
func (c Challenge) Amount() sdk.Coin {
	return c.Vote.Amount
}

// Creator implements `Voter`
func (c Challenge) Creator() sdk.AccAddress {
	return c.Vote.Creator
}

// VoteChoice implements `Voter`
func (c Challenge) VoteChoice() bool {
	return c.Vote.Vote
}

// MsgParams holds default parameters for a challenge
type MsgParams struct {
	MinArgumentLength int // min number of chars for argument
	MaxArgumentLength int // max number of chars for argument
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinArgumentLength: 10,
		MaxArgumentLength: 3000,
	}
}
