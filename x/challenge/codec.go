package challenge

import "github.com/cosmos/cosmos-sdk/codec"

var moduleCodec = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(CreateChallengeMsg{}, "challenge/SubmitChallengeMsg", nil)
	c.RegisterConcrete(LikeChallengeArgumentMsg{}, "challenge/LikeChallengeArgumentMsg", nil)
}
