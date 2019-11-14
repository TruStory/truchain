package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.processExpiringStakes(ctx)
}

func (k Keeper) processExpiringStakes(ctx sdk.Context) {
	logger := k.Logger(ctx)
	expiredStakes := make([]Stake, 0)
	fmt.Println("processing expired stakes")
	k.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time, func(stake Stake) bool {
		logger.Info(fmt.Sprintf("Processing expired stakeID %d argumentID %d", stake.ID, stake.ArgumentID))
		result, err := k.distributeReward(ctx, stake)
		if err != nil {
			panic(err)
		}
		stake.Expired = true
		stake.Result = &result
		k.setStake(ctx, stake)
		k.RemoveFromActiveStakeQueue(ctx, stake.ID, stake.EndTime)
		expiredStakes = append(expiredStakes, stake)
		return false
	})

	if len(expiredStakes) == 0 {
		return
	}

	b, err := k.codec.MarshalJSON(expiredStakes)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeInterestRewardPaid,
			sdk.NewAttribute(AttributeKeyExpiredStakes, string(b)),
		),
	)
}
