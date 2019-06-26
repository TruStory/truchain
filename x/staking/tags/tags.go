package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// staking tags
var (
	ArgumentCreator = "argument-creator"
	UpvoteCreator   = "upvote-creator"

	ActionCreateArgument    = "create-argument"
	ActionCreateUpvote      = "create-upvote"
	ActionUnhelpfulArgument = "unhelpful-argument"

	TxCategory = "truchain-staking"
	Action     = sdk.TagAction
	Category   = sdk.TagCategory
)
