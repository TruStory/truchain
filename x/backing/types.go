package backing

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Backing type
type Backing struct {
	*app.Vote `json:"vote"`
}

// ID implements `Voter.ID`
func (b Backing) ID() int64 {
	return b.Vote.ID
}

// StoryID implements `Voter.StoryID`
func (b Backing) StoryID() int64 {
	return b.Vote.StoryID
}

// Amount implements `Voter.Amount`
func (b Backing) Amount() sdk.Coin {
	return b.Vote.Amount
}

// Creator implements `Voter.Creator`
func (b Backing) Creator() sdk.AccAddress {
	return b.Vote.Creator
}

// Weight implements `Voter.Weight`
func (b Backing) Weight() sdk.Int {
	return b.Vote.Weight
}

// VoteChoice implements `Voter.VoteChoice`
func (b Backing) VoteChoice() bool {
	return b.Vote.Vote
}

// Timestamp implements `Voter.Timestamp`
func (b Backing) Timestamp() app.Timestamp {
	return b.Vote.Timestamp
}
