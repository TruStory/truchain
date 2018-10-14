package challenge

import (
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenge defines a challenge on a story
type Challenge struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	Amount       sdk.Coin       `json:"amount"`
	Arugment     string         `json:"arugment"`
	Creator      sdk.AccAddress `json:"creator"`
	Evidence     []url.URL      `json:"evidence,omitempty"`
	CreatedBlock int64          `json:"created_block"`
	CreatedTime  time.Time      `json:"created_time"`
}

// NewChallenge creates a new `Challenge` type with defaults
func NewChallenge(
	id int64,
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence []url.URL,
	createdBlock int64,
	createdTime time.Time) Challenge {

	return Challenge{
		ID:           id,
		StoryID:      storyID,
		Amount:       amount,
		Creator:      creator,
		Evidence:     evidence,
		CreatedBlock: createdBlock,
		CreatedTime:  time.Now(),
	}
}

// Params holds data for backing interest calculations
type Params struct {
	MinEvidenceCount int
	MaxEvidenceCount int
}

// NewParams creates a new BackingParams type with defaults
func NewParams() Params {
	return Params{
		MinEvidenceCount: 0,
		MaxEvidenceCount: 10,
	}
}
