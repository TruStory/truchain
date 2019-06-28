package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// staking tags
var (
	Creator       = "creator"
	RewardResults = "reward-results"

	ActionCreateArgument     = "create-argument"
	ActionCreateUpvote       = "create-upvote"
	ActionInterestRewardPaid = "interest-reward-paid"
	ActionUnhelpfulArgument  = "unhelpful-argument"

	TxCategory = "truchain-staking"
	Action     = sdk.TagAction
	Category   = sdk.TagCategory
)
