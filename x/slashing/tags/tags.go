package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// staking tags
var (
	ArgumentCreator       = "argument-creator"
	ArgumentCreatorJailed = "argument-creator-jailed"
	SlashResults          = "slash-results"

	ActionCreateSlash = "create-slash"

	TxCategory = "truchain-slashing"
	Action     = sdk.TagAction
	Category   = sdk.TagCategory
)
