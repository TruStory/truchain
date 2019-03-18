package argument

import (
	amino "github.com/tendermint/go-amino"
)

// RegisterAmino registers messages into the codec
func RegisterAmino(c *amino.Codec) {
	c.RegisterConcrete(LikeArgumentMsg{}, "argument/LikeArgumentMsg", nil)
}
