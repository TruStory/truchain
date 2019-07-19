package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// account tags
var (
	ActionUnjailAccounts = "unjail-accounts"

	UnjailedAccounts = "unjailed-accounts"

	TxCategory = "truchain-account"
	Action     = sdk.TagAction
	Category   = sdk.TagCategory
)
