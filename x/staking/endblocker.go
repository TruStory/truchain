package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := keeper.Logger(ctx)
	keeper.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time, func(stake Stake) bool {
		logger.Info(fmt.Sprintf("Processing expired stakeID %d argumentID %d", stake.ID, stake.ArgumentID))
		_, err := keeper.distributeReward(ctx, stake)
		if err != nil {
			panic(err)
		}
		stake.Expired = true
		keeper.setStake(ctx, stake)
		keeper.RemoveFromActiveStakeQueue(ctx, stake.ID, stake.EndTime)
		return false
	})
	return sdk.EmptyTags()
}
