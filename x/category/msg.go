package category

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateCategoryMsg defines a message to create a category
type CreateCategoryMsg struct {
	Title       string         `json:"title"`
	Creator     sdk.AccAddress `json:"creator"`
	Slug        string         `json:"slug"`
	Description string         `json:"description,omitempty"`
}

// NewCreateCategoryMsg creates a message to create a new category
func NewCreateCategoryMsg(
	title string,
	creator sdk.AccAddress,
	slug string,
	desc string) CreateCategoryMsg {

	return CreateCategoryMsg{
		Title:       title,
		Creator:     creator,
		Slug:        slug,
		Description: desc,
	}
}

// Route implements Msg
func (msg CreateCategoryMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg CreateCategoryMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg
func (msg CreateCategoryMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg CreateCategoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	params := DefaultMsgParams()
	if len(msg.Title) < params.MinTitleLen || len(msg.Title) > params.MaxTitleLen {
		return ErrInvalidCategoryMsg("Invalid title: " + msg.Title)
	}
	if len(msg.Slug) < params.MinSlugLen || len(msg.Slug) > params.MaxSlugLen {
		return ErrInvalidCategoryMsg("Invalid slug: " + msg.Slug)
	}
	if len(msg.Description) > params.MaxDescLen {
		return ErrInvalidCategoryMsg("Invalid description: " + msg.Description)
	}
	return nil
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg CreateCategoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
