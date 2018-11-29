package story

import (
	amino "github.com/tendermint/go-amino"
)

// RegisterAmino registers messages into the codec
func RegisterAmino(c *amino.Codec) {
	c.RegisterConcrete(AddEvidenceMsg{}, "story/SubmitEvidenceMsg", nil)
	c.RegisterConcrete(SubmitStoryMsg{}, "story/SubmitStoryMsg", nil)
	c.RegisterConcrete(FlagStoryMsg{}, "story/FlagStoryMsg", nil)
}
