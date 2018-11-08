package challenge

import (
	"net/url"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenge defines a user's challenge on a story
type Challenge struct {
	ID        int64          `json:"id"`
	Amount    sdk.Coin       `json:"amount"`
	Argument  string         `json:"argument"`
	Creator   sdk.AccAddress `json:"creator"`
	Evidence  []url.URL      `json:"evidence,omitempty"`
	Timestamp app.Timestamp  `json:"timestamp"`
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
