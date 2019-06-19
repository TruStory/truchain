package claim

import "github.com/cosmos/cosmos-sdk/codec"

var moduleCodec = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateClaim{}, "truchain/MsgCreateClaim", nil)
	c.RegisterConcrete(MsgDeleteClaim{}, "truchain/MsgDeleteClaim", nil)

	c.RegisterConcrete(Claim{}, "truchain/Claim", nil)
}
