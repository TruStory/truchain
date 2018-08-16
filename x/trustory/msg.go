package trustory

import (
	"encoding/json"
	"time"

	"github.com/TruStory/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PlaceBondMsg defines a message to bond to a story
type PlaceBondMsg struct {
	StoryID int64          `json:"story_id"`
	Stake   sdk.Coin       `json:"stake"`
	Creator sdk.AccAddress `json:"creator"`
	Period  time.Time      `json:"period"`
}

// NewPlaceBondMsg creates a message to place a new bond
func NewPlaceBondMsg(
	storyID int64,
	stake sdk.Coin,
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
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Stake.IsValid == false {
		return sdk.ErrInvalidBondAmount("Invalid bond amount: " + msg.Stake.String())
	}
	if msg.Period.IsZero == true {
		return sdk.ErrInvalidBondPeriod("Invalid bond period: " + msg.Period.String())
	}
	return nil
}
