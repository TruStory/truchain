package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// RegisterAmino registers messages into the codec
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(BackStoryMsg{}, "truchain/BackStoryMsg", nil)
	cdc.RegisterConcrete(AddCommentMsg{}, "truchain/AddCommentMsg", nil)
	cdc.RegisterConcrete(SubmitEvidenceMsg{}, "truchain/SubmitEvidenceMsg", nil)
	cdc.RegisterConcrete(SubmitStoryMsg{}, "truchain/SubmitStoryMsg", nil)
	cdc.RegisterConcrete(VoteMsg{}, "truchain/VoteMsg", nil)
}

// ============================================================================

// getSignBytes is a helper function for `Msg` types that serializes
// the message into json bytes.
func getSignBytes(msg sdk.Msg) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// getSigners is a helper function  for `Msg` types that returns
// the signers of the message.
func getSigners(addr sdk.AccAddress) []sdk.AccAddress {
	return []sdk.AccAddress{addr}
}
