package backing

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var moduleCodec = codec.New()

// RegisterCodec registers messages into the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(BackStoryMsg{}, "backing/BackStoryMsg", nil)
	cdc.RegisterConcrete(LikeBackingArgumentMsg{}, "backing/LikeBackingArgumentMsg", nil)
}
