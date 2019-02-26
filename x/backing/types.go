package backing

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Backing type
type Backing struct {
	app.Vote `json:"vote"`

	Interest sdk.Coin `json:"interest"`
}

// ID implements `Voter`
func (b Backing) ID() int64 {
	return b.Vote.ID
}

// StoryID implements `Voter`
func (b Backing) StoryID() int64 {
	return b.Vote.StoryID
}

// Amount implements `Voter`
func (b Backing) Amount() sdk.Coin {
	return b.Vote.Amount
}

// Creator implements `Voter`
func (b Backing) Creator() sdk.AccAddress {
	return b.Vote.Creator
}

// VoteChoice implements `Voter`
func (b Backing) VoteChoice() bool {
	return b.Vote.Vote
}
