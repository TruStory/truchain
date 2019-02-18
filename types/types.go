package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Timestamp records the timestamp for a type
type Timestamp struct {
	CreatedBlock int64     `json:"created_block,omitempty"`
	CreatedTime  time.Time `json:"created_time,omitempty"`
	UpdatedBlock int64     `json:"updated_block,omitempty"`
	UpdatedTime  time.Time `json:"updated_time,omitempty"`
}

// NewTimestamp creates a new default Timestamp
func NewTimestamp(blockHeader abci.Header) Timestamp {
	return Timestamp{
		blockHeader.Height,
		blockHeader.Time,
		blockHeader.Height,
		blockHeader.Time,
	}
}

// Update updates an existing Timestamp and returns a new one
func (t Timestamp) Update(blockHeader abci.Header) Timestamp {
	t.UpdatedBlock = blockHeader.Height
	t.UpdatedTime = blockHeader.Time

	return t
}

// Voter defines an interface for any kind of voter. It should be implemented
// by any type that has voting capabilities, implicit or explicit.
type Voter interface {
	ID() int64
	StoryID() int64
	Amount() sdk.Coin
	Creator() sdk.AccAddress
	VoteChoice() bool
}

// Vote is a type that defines a vote on a story. It serves as an inner struct
// for `Backing`, `Challenge`, and `TokenVote`, containing common fields.
type Vote struct {
	ID        int64          `json:"id"`
	StoryID   int64          `json:"story_id"`
	Amount    sdk.Coin       `json:"amount"`
	Argument  string         `json:"argument,omitempty"`
	Creator   sdk.AccAddress `json:"creator"`
	Vote      bool           `json:"vote"`
	Timestamp Timestamp      `json:"timestamp"`
}

func (v Vote) String() string {
	return fmt.Sprintf("Vote<%v %t>", v.Amount, v.Vote)
}

// NewVote creates a new Vote type with defaults
func NewVote(
	id int64,
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	vote bool,
	timestamp Timestamp) Vote {

	return Vote{id, storyID, amount, "", creator, vote, timestamp}
}
