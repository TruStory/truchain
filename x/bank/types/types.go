package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

const (
	ModuleName        = "trubank2"
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName

	AttributeRecipient = "recipient"
)
