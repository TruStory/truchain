package account

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterKey{}, "truchain/MsgRegisterKey", nil)
	cdc.RegisterConcrete(AppAccount{}, "truchain/AppAccount", nil)
	cdc.RegisterConcrete(PrimaryAccount{}, "truchain/PrimaryAccount", nil)
	cdc.RegisterConcrete(MsgUpdateParams{}, "account/MsgUpdateParams", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
