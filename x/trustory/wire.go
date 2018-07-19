package trustory

import "github.com/cosmos/cosmos-sdk/wire"

// RegisterWire registers messages into the wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(SubmitStoryMsg{}, "trustory/SubmitStoryMsg", nil)
	cdc.RegisterConcrete(VoteMsg{}, "trustory/VoteMsg", nil)
}
