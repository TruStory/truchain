package staking

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/staking/tags"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := keeper.Logger(ctx)
	results := make([]RewardResult, 0)
	keeper.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time, func(stake Stake) bool {
		logger.Info(fmt.Sprintf("Processing expired stakeID %d argumentID %d", stake.ID, stake.ArgumentID))
		result, err := keeper.distributeReward(ctx, stake)
		if err != nil {
			panic(err)
		}
		stake.Expired = true
		stake.Result = &result
		keeper.setStake(ctx, stake)
		keeper.RemoveFromActiveStakeQueue(ctx, stake.ID, stake.EndTime)
		results = append(results, result)
		return false
	})

	if len(results) == 0 {
		return sdk.EmptyTags()
	}
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return append(app.PushTag,
		sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Action, tags.ActionInterestRewardPaid,
			tags.RewardResults, b,
		)...,
	)
}
