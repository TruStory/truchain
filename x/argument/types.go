package argument

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// getArgument for a story
type Argument struct {
	ID      int64 `json:"id"`
	StoryID int64 `json:"story_id"`

	// association with backing or challenge
	StakeID int64 `json:"stake_id"`

	Body      string         `json:"body"`
	Creator   sdk.AccAddress `json:"creator"`
	Timestamp app.Timestamp  `json:"timestamp"`
}

// Like for an argument
type Like struct {
	ArgumentID int64          `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
	Timestamp  app.Timestamp  `json:"timestamp"`
}
