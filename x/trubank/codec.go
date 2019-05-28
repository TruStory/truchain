package trubank

import (
	amino "github.com/tendermint/go-amino"
)

// RegisterAmino registers messages into the codec
func RegisterAmino(c *amino.Codec) {
	c.RegisterConcrete(PayRewardMsg{}, "trubank/PayRewardMsg", nil)
}
