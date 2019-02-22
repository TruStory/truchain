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

// VoteChoice implements `Voter`
func (v TokenVote) VoteChoice() bool {
	return v.Vote.Vote
}

// MsgParams holds default parameters for a vote
type MsgParams struct {
	MinArgumentLength int // min number of chars for argument
	MaxArgumentLength int // max number of chars for argument
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinArgumentLength: 10,
		MaxArgumentLength: 3000,
	}
}
