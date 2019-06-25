package account

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	// cdc.RegisterInterface((*TruAccount)(nil), nil)
	cdc.RegisterConcrete(MsgRegisterKey{}, "truchain/MsgRegisterKey", nil)
	cdc.RegisterConcrete(AppAccount{}, "truchain/AppAccount", nil)
	// cdc.RegisterConcrete(AppAccount{}, "sdkAuth/AppAccount", nil)

	// cdc.RegisterConcrete(&types.AppAccount{}, "types/AppAccount", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
