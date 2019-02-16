package category

import (
	app "github.com/TruStory/truchain/types"
)

// Category is a type that defines the category for a story
type Category struct {
	ID          int64         `json:"id"`
	Title       string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description,omitempty"`
	Timestamp   app.Timestamp `json:"timestamp"`
}

// Denom returns the name of the coin, alias for slug
func (c Category) Denom() string {
	return c.Slug
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
