package distribution

import "github.com/cosmos/cosmos-sdk/codec"

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	//RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
