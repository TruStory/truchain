package stake

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg defines data common to backing, challenge, and
// token vote messages.
type Msg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Argument string         `json:"argument,omitempty"`
	Creator  sdk.AccAddress `json:"creator"`
}

// Voter defines an interface for any kind of voter. It should be implemented
// by any type that has voting capabilities, implicit or explicit.
type Voter interface {
	ID() int64
	StoryID() int64
	Amount() sdk.Coin
	Weight() sdk.Int
	Creator() sdk.AccAddress
	VoteChoice() bool
	UpdateWeight(sdk.Int)
	Timestamp() app.Timestamp
}

// Vote is a type that defines a vote on a story. It serves as an inner struct
// for `Backing`, `Challenge`, and `TokenVote`, containing common fields.
type Vote struct {
	ID         int64          `json:"id"`
	StoryID    int64          `json:"story_id"`
	Amount     sdk.Coin       `json:"amount"`
	ArgumentID int64          `json:"argument_id,omitempty"`
	Weight     sdk.Int        `json:"weight,omitempty"`
	Creator    sdk.AccAddress `json:"creator"`
	Vote       bool           `json:"vote"`
	Timestamp  app.Timestamp  `json:"timestamp"`
}

// UpdateWeight mutates the vote weight as a result of weighted voting
func (v *Vote) UpdateWeight(weight sdk.Int) {
	v.Weight = weight
}

func (v Vote) String() string {
	return fmt.Sprintf("Vote<%v %t>", v.Amount, v.Vote)
}

// NewVote creates a new Vote type with defaults
func NewVote(
	id int64,
	storyID int64,
	amount sdk.Coin,
	weight sdk.Int,
	creator sdk.AccAddress,
	vote bool,
	timestamp app.Timestamp) Vote {

	return Vote{id, storyID, amount, 0, sdk.NewInt(0), creator, vote, timestamp}
}
