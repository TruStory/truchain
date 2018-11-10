package types

import (
	"net/url"
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

// Vote is a type that defines a vote on a story
type Vote struct {
	ID        int64          `json:"id"`
	Amount    sdk.Coin       `json:"amount"`
	Comment   string         `json:"comment,omitempty"`
	Creator   sdk.AccAddress `json:"creator"`
	Evidence  []url.URL      `json:"evidence,omitempty"`
	Vote      bool           `json:"vote"`
	Timestamp Timestamp      `json:"timestamp"`
}

// NewVote creates a new Vote type with defaults
func NewVote(
	id int64, amount sdk.Coin, creator sdk.AccAddress, vote bool, timestamp Timestamp) Vote {

	return Vote{id, amount, "", creator, nil, vote, timestamp}
}
