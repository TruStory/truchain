package vote

import (
	amino "github.com/tendermint/go-amino"
)

// RegisterAmino registers messages into the codec
func RegisterAmino(c *amino.Codec) {
	c.RegisterConcrete(CreateVoteMsg{}, "vote/CreateVoteMsg", nil)
}
