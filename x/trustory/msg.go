package trustory

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PlaceBondMsg defines a message to bond to a story
type PlaceBondMsg struct {
	StoryID int64
	Stake   sdk.Coins
	Creator sdk.AccAddress
	Period  time.Time
}

// NewPlaceBondMsg creates a message to place a new bond
func NewPlaceBondMsg(
	storyID int64,
	stake sdk.Coins,
	creator sdk.AccAddress,
	period time.Time) PlaceBondMsg {
	return PlaceBondMsg{
		StoryID: storyID,
		Stake:   stake,
		Creator: creator,
		Period:  period,
	}
}

// SubmitStoryMsg defines a message to create a story
type SubmitStoryMsg struct {
	Body    string
	Creator sdk.AccAddress
}

// NewSubmitStoryMsg submits a message with a new story
func NewSubmitStoryMsg(body string, creator sdk.AccAddress) SubmitStoryMsg {
	return SubmitStoryMsg{
		Body:    body,
		Creator: creator,
	}
}

// Type implements Msg
func (msg SubmitStoryMsg) Type() string {
	return "truStory"
}

// Get implement Msg
func (msg SubmitStoryMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes implements Msg
func (msg SubmitStoryMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements Msg
func (msg SubmitStoryMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len(strings.TrimSpace(msg.Body)) <= 0 {
		return ErrInvalidBody("Cannot submit a story with an empty body")
	}

	return nil
}

func (msg SubmitStoryMsg) String() string {
	return fmt.Sprintf("SubmitStoryMsg{%v}", msg.Body)
}

// VoteMsg defines a message to vote on a specific story
type VoteMsg struct {
	StoryID int64
	Option  string
	Voter   sdk.AccAddress
}

// NewVoteMsg creates a VoteMsg instance
func NewVoteMsg(storyID int64, option string, voter sdk.AccAddress) VoteMsg {
	return VoteMsg{
		StoryID: storyID,
		Option:  option,
		Voter:   voter,
	}
}

// Type implements Msg
func (msg VoteMsg) Type() string {
	return "truStory"
}

// Get implements Msg
func (msg VoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes implements Msg
func (msg VoteMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements Msg
func (msg VoteMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

// ValidateBasic implements Msg
func (msg VoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Voter.String())
	}
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}

	if len(strings.TrimSpace(msg.Option)) <= 0 {
		return ErrInvalidOption("Option can't be blank")
	}

	return nil
}

// String implements Msg
func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoteMsg{%v, %v, %v}", msg.StoryID, msg.Option, msg.Voter)
}
