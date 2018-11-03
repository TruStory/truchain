package types

import (
	"time"

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
