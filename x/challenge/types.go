package challenge

import (
	"fmt"

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

// StoryID implements `Voter`
func (c Challenge) StoryID() int64 {
	return c.Vote.StoryID
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

// Timestamp implements `Voter.Timestamp`
func (c Challenge) Timestamp() app.Timestamp {
	return c.Vote.Timestamp
}

func (c Challenge) String() string {
	return fmt.Sprintf("Challenge<%s>", c.Vote.String())
}
