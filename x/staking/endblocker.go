package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)
	expiredStakes := make([]Stake, 0)
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
		expiredStakes = append(expiredStakes, stake)
		return false
	})

	if len(expiredStakes) == 0 {
		return
	}

	b, err := keeper.codec.MarshalJSON(expiredStakes)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ModuleName,
			sdk.NewAttribute(sdk.AttributeKeyAction, AttributeValueInterestRewardPaid),
			sdk.NewAttribute(AttributeKeyExpiredStakes, string(b)),
		),
	)
}
