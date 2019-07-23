package slashing

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSlashArgument{}, "truchain/MsgSlashArgument", nil)
	cdc.RegisterConcrete(MsgUpdateParams{}, "slashing/MsgUpdateParams", nil)

	cdc.RegisterConcrete(Slash{}, "truchain/Slash", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
