package challenge

import (
	app "github.com/TruStory/truchain/types"
)

// Challenge defines a user's challenge on a story
type Challenge struct {
	app.Vote
}

// AmountDenom implements `Voter`
func (c Challenge) AmountDenom() string {
	return c.Amount.Denom
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
