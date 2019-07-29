package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// slashing tags
var (
	SlashResults      = "slash-results"
	ActionCreateSlash = "create-slash"
	MinSlashCount     = "min-slash-count"

	TxCategory = "truchain-slashing"
	Action     = sdk.TagAction
	Category   = sdk.TagCategory
)
