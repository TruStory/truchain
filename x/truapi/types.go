package truapi

import (
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

	Argument argument.Argument `graphql:"-"`
}
