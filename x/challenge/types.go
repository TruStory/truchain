package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
	MinEvidenceCount  int // min number of evidence URLs
	MaxEvidenceCount  int // max number of evidence URLs
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinArgumentLength: 10,
		MaxArgumentLength: 340,
		MinEvidenceCount:  0,
		MaxEvidenceCount:  10,
	}
}
