package trustory

import (
	"encoding/json"
	"time"

	"github.com/TruStory/cosmos-sdk/types"
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

// Type implements Msg
func (msg PlaceBondMsg) Type() string {
	return "truStory"
}

// GetSignBytes implements Msg
func (msg PlaceBondMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic implements Msg
func (msg PlaceBondMsg) ValidateBasic() types.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	// if len(strings.TrimSpace(msg.Body)) <= 0 {
	// 	return ErrInvalidBody("Cannot submit a story with an empty body")
	// }
	return nil
}

// ValidateBasic implements Msg
// func (msg VoteMsg) ValidateBasic() sdk.Error {
// 	if len(msg.Voter) == 0 {
// 		return sdk.ErrInvalidAddress("Invalid address: " + msg.Voter.String())
// 	}
// 	if msg.StoryID <= 0 {
// 		return ErrInvalidStoryID("StoryID cannot be negative")
// 	}

// 	if len(strings.TrimSpace(msg.Option)) <= 0 {
// 		return ErrInvalidOption("Option can't be blank")
// 	}

// 	return nil
// }
