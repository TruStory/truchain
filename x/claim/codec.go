package claim

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateClaim{}, "truchain/MsgCreateClaim", nil)
	c.RegisterConcrete(MsgDeleteClaim{}, "truchain/MsgDeleteClaim", nil)

	c.RegisterConcrete(Claim{}, "truchain/Claim", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
