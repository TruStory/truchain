package challenge

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Challenge defines a challenge on a story
type Challenge struct {
	ID           int64
	StoryID      int64
	Amount       sdk.Coin
	Arugment     string
	Creator      sdk.AccAddress
	Evidence     string // TODO: in here or story?
	CreatedBlock int64
	CreatedTime  time.Time
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

// StoryList defines a list of challenged story IDs
type StoryList []int64

// IsEmpty checks if the story list is empty
func (sl StoryList) IsEmpty() bool {
	if len(sl) == 0 {
		return true
	}
	return false
}
