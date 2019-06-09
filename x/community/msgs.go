package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgNewCommunity represents the type of the message for creating new community
	TypeMsgNewCommunity = "new_community"
)

// MsgNewCommunity defines the message to create new community
type MsgNewCommunity struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Creator     sdk.AccAddress `json:"creator"`
}

// NewMsgNewCommunity returns the messages to create a new community
func NewMsgNewCommunity(name, slug, description string, creator sdk.AccAddress) MsgNewCommunity {
	return MsgNewCommunity{
		Name:        name,
		Slug:        slug,
		Description: description,
		Creator:     creator,
	}
}

// ValidateBasic implements Msg
func (msg MsgNewCommunity) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Creator.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgNewCommunity) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgNewCommunity) Type() string { return TypeMsgNewCommunity }

// GetSignBytes implements Msg
func (msg MsgNewCommunity) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg MsgNewCommunity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}
