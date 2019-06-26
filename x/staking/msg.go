package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgSubmitArgument{}
var _ sdk.Msg = &MsgSubmitUpvote{}
var _ sdk.Msg = &MsgDeleteArgument{}

const (
	TypeMsgSubmitArgument = "submit_argument"
	TypeMsgSubmitUpvote   = "submit_upvote"
	TypeMsgDeleteArgument = "delete_argument"
)

// MsgSubmitArgument msg for creating an argument.
type MsgSubmitArgument struct {
	ClaimID   uint64         `json:"claim_id"`
	Summary   string         `json:"summary"`
	Body      string         `json:"body"`
	StakeType StakeType      `json:"stake_type"`
	Creator   sdk.AccAddress `json:"creator"`
}

func (MsgSubmitArgument) Route() string {
	return RouterKey
}

func (MsgSubmitArgument) Type() string {
	return TypeMsgSubmitArgument
}

func (msg MsgSubmitArgument) ValidateBasic() sdk.Error {
	if !msg.StakeType.ValidForArgument() {
		return ErrCodeInvalidStakeType(msg.StakeType)
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address %s", msg.Creator.String()))
	}

	if len(msg.Body) == 0 {
		return ErrCodeInvalidBodyLength()
	}
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgSubmitArgument) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgSubmitArgument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgSubmitUpvote msg for upvoting an argument.
type MsgSubmitUpvote struct {
	ArgumentID uint64         `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
}

func (MsgSubmitUpvote) Route() string {
	return RouterKey
}

func (MsgSubmitUpvote) Type() string {
	return TypeMsgSubmitUpvote
}

func (MsgSubmitUpvote) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgSubmitUpvote) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgSubmitUpvote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgDeleteArgument msg for upvoting an argument.
type MsgDeleteArgument struct {
	ArgumentID uint64         `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
}

func (MsgDeleteArgument) Route() string {
	return RouterKey
}

func (MsgDeleteArgument) Type() string {
	return TypeMsgDeleteArgument
}

func (msg MsgDeleteArgument) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Must provide a valid address")
	}
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgDeleteArgument) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgDeleteArgument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
