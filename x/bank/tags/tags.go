package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// staking tags
var (
	ActionPayReward = "pay-reward"
	TxCategory      = "truchain-bank"
	Action          = sdk.TagAction
	Category        = sdk.TagCategory
	Sender          = sdk.TagSender
	Recipient       = "recipient"
)
