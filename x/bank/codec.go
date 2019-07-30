package bank

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgPayReward{}, "truchain/MsgPayReward", nil)
	c.RegisterConcrete(MsgSendGift{}, "truchain/MsgSendGift", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "bank/MsgUpdateParams", nil)

	c.RegisterConcrete(Transaction{}, "truchain/Transaction", nil)
}

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
