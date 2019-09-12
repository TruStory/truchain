package staking

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.processExpiringStakes(ctx)
	keeper.distributeInflation(ctx)
}

func (k Keeper) processExpiringStakes(ctx sdk.Context) {
	logger := k.Logger(ctx)
	expiredStakes := make([]Stake, 0)
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

func (k Keeper) distributeInflation(ctx sdk.Context) {
	// TODO: take this from params
	// 20%
	userRewardAllocation := sdk.NewDecWithPrec(200, 3)

	acc := k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	fmt.Println(acc)

	acc2 := k.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName)
	fmt.Println(acc2)

	userInflation := acc.GetCoins().AmountOf(app.StakeDenom)
	userInflationDec := sdk.NewDecFromIntWithPrec(userInflation, 3)
	userRewardAmount := userInflationDec.Mul(userRewardAllocation).RoundInt()
	userRewardCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userRewardAmount))

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, UserRewardPoolName, userRewardCoins)
	if err != nil {
		panic(err)
	}
}
