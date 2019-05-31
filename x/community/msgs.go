package community

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	params := DefaultMsgParams()
	if len(msg.Name) < params.MinNameLen || len(msg.Name) > params.MaxNameLen {
		return ErrInvalidCommunityMsg(
			fmt.Sprintf("Name must be between %d-%d chars in length", params.MinNameLen, params.MaxNameLen),
		)
	}
	if len(msg.Slug) < params.MinSlugLen || len(msg.Slug) > params.MaxSlugLen {
		return ErrInvalidCommunityMsg(
			fmt.Sprintf("Slug must be between %d-%d chars in length", params.MinSlugLen, params.MaxSlugLen),
		)
	}
	if len(msg.Description) > params.MaxDescriptionLen {
		return ErrInvalidCommunityMsg(
			fmt.Sprintf("Description must be less than %d chars in length", params.MaxDescriptionLen),
		)
	}
	return nil
}

// Route implements Msg
func (msg MsgNewCommunity) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg MsgNewCommunity) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg
func (msg MsgNewCommunity) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg MsgNewCommunity) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
