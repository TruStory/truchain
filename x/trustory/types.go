package trustory

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story defines the basic properties of a votable story
type Story struct {
	Body        string      `json:"body"`
	Creator     sdk.Address `json:"creator"`
	SubmitBlock int64       `json:"submit_block`
	State       string      `json:"state"`
	YesVotes    int64       `json:"yes_votes`
	NoVotes     int64       `json:"no_votes"`
}

// NewStory creates a new story
func NewStory(body string, creator sdk.Address, blockHeight int64) Story {
	return Story{
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
	Creator sdk.Address
}

// NewSubmitStoryMsg submits a message with a new story
func NewSubmitStoryMsg(body string, creator sdk.Address) SubmitStoryMsg {
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
func (msg SubmitStoryMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Submitter}
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len(strings.TimeSpace(msg.Body)) <= 0 {
		return ErrInvalidBody("Cannot submit a story with an empty body")
	}

	return nil
}

func (msg SubmitStoryMsg) String() string {
	return fmt.Sprintf("SubmitStoryMsg{%v}", msg.Body)
}

// VoteMsg defines a message to vote on a story
type VoteMsg struct {
	StoryID int64
	Option  string
	Voter   sdk.Address
}
