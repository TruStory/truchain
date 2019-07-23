package claim

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateClaim{}
var _ sdk.Msg = &MsgEditClaim{}

// MsgCreateClaim defines a message to submit a story
type MsgCreateClaim struct {
	CommunityID string         `json:"community_id"`
	Body        string         `json:"body"`
	Creator     sdk.AccAddress `json:"creator"`
	Source      string         `json:"source,omitempty"`
}

// NewMsgCreateClaim creates a new message to create a claim
func NewMsgCreateClaim(communityID, body string, creator sdk.AccAddress, source string) MsgCreateClaim {
	return MsgCreateClaim{
		CommunityID: communityID,
		Body:        body,
		Creator:     creator,
		Source:      source,
	}
}

// Route is the name of the route for claim
func (msg MsgCreateClaim) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCreateClaim) Type() string {
	return "create_claim"
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreateClaim) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBodyTooShort(msg.Body)
	}
	if len(msg.CommunityID) == 0 {
		return ErrInvalidCommunityID(msg.CommunityID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCreateClaim) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCreateClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgDeleteClaim defines a message to submit a story
type MsgDeleteClaim struct {
	ID      uint64         `json:"id"`
	Creator sdk.AccAddress `json:"creator"`
}

// Route is the name of the route for claim
func (msg MsgDeleteClaim) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgDeleteClaim) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgDeleteClaim) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownClaim(msg.ID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgDeleteClaim) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgDeleteClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgEditClaim defines a message to submit a story
type MsgEditClaim struct {
	ID     uint64         `json:"id"`
	Body   string         `json:"body"`
	Editor sdk.AccAddress `json:"editor"`
}

// NewMsgEditClaim creates a new message to edit a claim
func NewMsgEditClaim(id uint64, body string, editor sdk.AccAddress) MsgEditClaim {
	return MsgEditClaim{
		ID:     id,
		Body:   body,
		Editor: editor,
	}
}

// Route is the name of the route for claim
func (msg MsgEditClaim) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgEditClaim) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgEditClaim) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownClaim(msg.ID)
	}
	if len(msg.Editor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Editor.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgEditClaim) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgEditClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Editor)}
}
