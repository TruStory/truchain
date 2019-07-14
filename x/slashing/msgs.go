package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgSlashArgument represents the type of the message for creating new community
	TypeMsgSlashArgument = "slash_argument"
)

// MsgSlashArgument defines the message to slash an argument
type MsgSlashArgument struct {
	StakeID             uint64
	SlashType           SlashType
	SlashReason         SlashReason
	SlashDetailedReason string
	Creator             sdk.AccAddress
}

// NewMsgSlashArgument returns the messages to slash an argument
func NewMsgSlashArgument(stakeID uint64, slashType SlashType, slashReason SlashReason, slashDetailedReason string, creator sdk.AccAddress) MsgSlashArgument {
	return MsgSlashArgument{
		StakeID:             stakeID,
		SlashType:           slashType,
		SlashReason:         slashReason,
		SlashDetailedReason: slashDetailedReason,
		Creator:             creator,
	}
}

// ValidateBasic implements Msg
func (msg MsgSlashArgument) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Creator.String()))
	}

	if msg.SlashReason == SlashReasonOthers && len(msg.SlashDetailedReason) == 0 {
		return ErrInvalidSlashReason("Need to have detailed reason when chosen Others")
	}

	return nil
}

// Route implements Msg
func (msg MsgSlashArgument) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgSlashArgument) Type() string { return TypeMsgSlashArgument }

// GetSignBytes implements Msg
func (msg MsgSlashArgument) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg MsgSlashArgument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}
