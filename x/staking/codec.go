package staking

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgSubmitArgument{}, "truchain/MsgSubmitArgument", nil)
	c.RegisterConcrete(MsgSubmitUpvote{}, "truchain/MsgUpvoteArgument", nil)
	c.RegisterConcrete(MsgEditArgument{}, "truchain/MsgEditArgument", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "staking/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "staking/MsgRemoveAdmin", nil)

	c.RegisterConcrete(Stake{}, "truchain/Stake", nil)
	c.RegisterConcrete(Argument{}, "truchain/Argument", nil)

}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
