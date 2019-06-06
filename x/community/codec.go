package community

import "github.com/cosmos/cosmos-sdk/codec"

var moduleCodec = codec.New()

// RegisterCodec registers messages into the codec
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgNewCommunity{}, "community/MsgNewCommunity", nil)
}
