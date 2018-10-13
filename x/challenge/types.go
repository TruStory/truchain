package challenge

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenge defines a challenge on a story
type Challenge struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	Amount       sdk.Coin       `json:"amount"`
	Arugment     string         `json:"arugment,omitempty"`
	Creator      sdk.AccAddress `json:"creator"`
	Evidence     string         `json:"evidence,omitempty"`
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
	evidence string,
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
