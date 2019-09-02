package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// bank tags
var (
	ActionPayReward = "pay-reward"
	ActionPayGift   = "pay-gift"
	TxCategory      = "truchain-bank"
	Action          = sdk.TagAction
	Category        = sdk.TagCategory
	Sender          = sdk.TagSender
	Recipient       = "recipient"
)