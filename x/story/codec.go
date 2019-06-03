package story

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var moduleCodec = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(SubmitStoryMsg{}, "story/SubmitStoryMsg", nil)
}
