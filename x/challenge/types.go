package challenge

import (
	app "github.com/TruStory/truchain/types"
)

// Challenge defines a user's challenge on a story
type Challenge struct {
	app.Vote
}

// Params holds default parameters for a challenge
type Params struct {
	MinArgumentLength int // min number of chars for argument
	MaxArgumentLength int // max number of chars for argument
	MinEvidenceCount  int // min number of evidence URLs
	MaxEvidenceCount  int // max number of evidence URLs
}

// DefaultParams creates a new Params type with defaults
func DefaultParams() Params {
	return Params{
		MinArgumentLength: 10,
		MaxArgumentLength: 340,
		MinEvidenceCount:  0,
		MaxEvidenceCount:  10,
	}
}
