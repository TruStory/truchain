package category

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Category is a type that defines the category for a story
type Category struct {
	ID          int64         `json:"id,omitempty"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Description string        `json:"description,omitempty"`
	TotalCred   sdk.Coin      `json:"total_cred"`
	Timestamp   app.Timestamp `json:"timestamp,omitempty"`
}

// Denom returns the name of the coin, alias for slug
func (c Category) Denom() string {
	return c.Slug
}

func (c Category) String() string {
	return fmt.Sprintf(
		"Category<%s %s %s>", c.Title, c.Slug, c.Description)
}

// MsgParams holds data for category parameters
type MsgParams struct {
	MinTitleLen int
	MaxTitleLen int
	MinSlugLen  int
	MaxSlugLen  int
	MaxDescLen  int
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinTitleLen: 5,
		MaxTitleLen: 25,
		MinSlugLen:  3,
		MaxSlugLen:  15,
		MaxDescLen:  140,
	}
}
