package backing

import (
	"github.com/TruStory/truchain/x/argument"
	amino "github.com/tendermint/go-amino"
)

// RegisterAmino registers messages into the codec
func RegisterAmino(c *amino.Codec) {
	c.RegisterConcrete(BackStoryMsg{}, "backing/BackStoryMsg", nil)
	c.RegisterConcrete(argument.LikeArgumentMsg{}, "argument/LikeArgumentMsg", nil)
}
