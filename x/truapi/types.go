package truapi

import (
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CredArgument represents an argument that earned cred based on likes.
type CredArgument struct {
	ID        int64          `json:"id" graphql:"id" `
	StoryID   int64          `json:"storyId" graphql:"storyId"`
	Body      string         `json:"body"`
	Creator   sdk.AccAddress `json:"creator" `
	Timestamp app.Timestamp  `json:"timestamp"`
	Vote      bool           `json:"vote"`
	Amount    sdk.Coin       `json:"coin"`
}

// CommentNotificationRequest is the payload sent to pushd for sending notifications.
type CommentNotificationRequest struct {
	// ID is the comment id.
	ID              int64     `json:"id"`
	ArgumentCreator string    `json:"argument_creator"`
	ArgumentID      int64     `json:"argumentId"`
	StoryID         int64     `json:"storyId"`
	Creator         string    `json:"creator"`
	Timestamp       time.Time `json:"timestamp"`
}

// StakeArgument is the argument that a stake references to.
type StakeArgument struct {
	argument.Argument
}
