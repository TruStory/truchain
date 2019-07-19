package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// slashing tags
var (
	ArgumentCreator       = "argument-creator"
	ArgumentCreatorJailed = "argument-creator-jailed"
	SlashResults          = "slash-results"

	ActionCreateSlash = "create-slash"

	TxCategory = "truchain-slashing"
	Action     = sdk.TagAction
	Category   = sdk.TagCategory
)
