package claim

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateClaim{}

// MsgCreateClaim defines a message to submit a story
type MsgCreateClaim struct {
	CommunityID uint64         `json:"community_id"`
	Body        string         `json:"body"`
	Creator     sdk.AccAddress `json:"creator"`
	Source      string         `json:"source,omitempty"`
}

// NewMsgCreateClaim creates a new message to create a claim
func NewMsgCreateClaim(communityID uint64, body string, creator sdk.AccAddress, source string) MsgCreateClaim {
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
	if msg.CommunityID == 0 {
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
		return ErrInvalidID()
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
