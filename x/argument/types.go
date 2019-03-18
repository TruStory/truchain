package argument

import (
	app "github.com/TruStory/truchain/types"
)

// Argument for a story
type Argument struct {
	ID      int64 `json:"id"`
	StoryID int64 `json:"story_id"`

	// association with backing or challenge
	StakeID   int64       `json:"stake_id"`
	StakeType interface{} `json:"stake_type"`

	Body      string        `json:"body"`
	Timestamp app.Timestamp `json:"timestamp"`
}

// Like for an argument
type Like struct {
	ArgumentID int64
	Timestamp  app.Timestamp
}
