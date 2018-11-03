package category

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Category is a type that defines the category for a story
type Category struct {
	ID          int64          `json:"id"`
	Creator     sdk.AccAddress `json:"creator"`
	Title       string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description,omitempty"`
	Timestamp   app.Timestamp  `json:"timestamp"`
}

// CoinName returns the name of the coin, alias for slug
func (c Category) CoinName() string {
	return c.Slug
}

// Params holds data for category parameters
type Params struct {
	MinTitleLen int
	MaxTitleLen int
	MinSlugLen  int
	MaxSlugLen  int
	MaxDescLen  int
}

// DefaultParams creates a new CategoryParams type with defaults
func DefaultParams() Params {
	return Params{
		MinTitleLen: 5,
		MaxTitleLen: 25,
		MinSlugLen:  3,
		MaxSlugLen:  15,
		MaxDescLen:  140,
	}
}
