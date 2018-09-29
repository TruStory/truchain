package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VoteMsg defines a message to vote on a story
type VoteMsg struct {
	StoryID int64          `json:"story_id"`
	Creator sdk.AccAddress `json:"creator"`
	Amount  sdk.Coin       `json:"amount"`
	Vote    bool           `json:"vote"`
}

// NewVoteMsg creates a new message to vote on a story
func NewVoteMsg(storyID int64, creator sdk.AccAddress, amount sdk.Coin, vote bool) VoteMsg {
	return VoteMsg{
		StoryID: storyID,
		Creator: creator,
		Amount:  amount,
		Vote:    vote,
	}
}

// Type implements Msg
func (msg VoteMsg) Type() string { return "Vote" }

// Name implements Msg
func (msg VoteMsg) Name() string { return msg.Type() }

// GetSignBytes implements Msg
func (msg VoteMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg VoteMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return ErrInvalidAmount("Invalid stake amount: " + msg.Amount.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg VoteMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
}
