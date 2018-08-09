package trustory

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Bond struct {
	ID           int64          `json:"id"`
	Amount       float64        `json:"amount"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	Period       int64          `json:"period"`
	StoryID      int64          `json:"story_id"`
}

// Story defines the basic properties of a votable story
type Story struct {
	ID          int64          `json:"id"`
	Body        string         `json:"body"`
	Creator     sdk.AccAddress `json:"creator"`
	SubmitBlock int64          `json:"submit_block`
	State       string         `json:"state"`
	YesVotes    int64          `json:"yes_votes`
	NoVotes     int64          `json:"no_votes"`
}

// NewStory creates a new story
func NewStory(
	id int64,
	body string,
	creator sdk.AccAddress,
	blockHeight int64) Story {
	return Story{
		ID:          id,
		Body:        body,
		Creator:     creator,
		SubmitBlock: blockHeight,
		State:       "Created",
		YesVotes:    0,
		NoVotes:     0,
	}
}

// updateVote updates the votes for each
func (s *Story) updateVote(option string, amount int64) sdk.Error {
	switch option {
	case "Yes":
		s.YesVotes += amount
		return nil
	case "No":
		s.NoVotes += amount
		return nil
	default:
		return ErrInvalidOption("Invalid option: " + option)
	}
}

//--------------------------------------------------------
//--------------------------------------------------------

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
