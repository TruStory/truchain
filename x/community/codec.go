package community

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers messages into the codec
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgNewCommunity{}, "community/MsgNewCommunity", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "community/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "community/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "community/MsgUpdateParams", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
