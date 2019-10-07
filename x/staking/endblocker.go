package staking

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.processExpiringStakes(ctx)
	keeper.distributeInflation(ctx)

	supplyTotal := keeper.supplyKeeper.GetSupply(ctx)
	fmt.Println("supply " + supplyTotal.String())

	//distAcc := keeper.supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName)
	//fmt.Println(distAcc.GetName() + "            " + distAcc.GetCoins().String())

	rewardAcc := keeper.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName)
	fmt.Println(rewardAcc.GetName() + " " + rewardAcc.GetCoins().String())
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
	userRewardAllocation := k.GetParams(ctx).UserRewardAllocation

	acc := k.supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName)

	userInflation := acc.GetCoins().AmountOf(app.StakeDenom)
	userInflationDec := sdk.NewDecFromInt(userInflation)
	userRewardAmount := userInflationDec.Mul(userRewardAllocation).RoundInt()
	userRewardCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userRewardAmount))
	//fmt.Println("user reward coins " + userRewardCoins.String())

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, distribution.ModuleName, UserRewardPoolName, userRewardCoins)
	if err != nil {
		panic(err)
	}
}
