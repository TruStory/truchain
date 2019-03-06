package vote

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokenVote defines a simple token vote on a story
type TokenVote struct {
	app.Vote `json:"vote"`
}

// ID implements `Voter`
func (v TokenVote) ID() int64 {
	return v.Vote.ID
}

// StoryID implements `Voter`
func (v TokenVote) StoryID() int64 {
	return v.Vote.ID
}

// Amount implements `Voter`
func (v TokenVote) Amount() sdk.Coin {
	return v.Vote.Amount
}

// Creator implements `Voter`
func (v TokenVote) Creator() sdk.AccAddress {
	return v.Vote.Creator
}

// Weight implements `Voter`
func (v TokenVote) Weight() sdk.Int {
	return v.Vote.Weight
}

// FullVote returns the vote for setter purposes
func (v TokenVote) FullVote() app.Vote {
	return v.Vote
}

// VoteChoice implements `Voter`
func (v TokenVote) VoteChoice() bool {
	return v.Vote.Vote
}

// Timestamp implements `Voter.Timestamp`
func (v TokenVote) Timestamp() app.Timestamp {
	return v.Vote.Timestamp
}
