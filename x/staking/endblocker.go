package staking

import (
	"fmt"
	"math"

	chain "github.com/TruStory/truchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)
	expiredStakes := make([]Stake, 0)
	initialMultiples := make(map[string]int)
	keeper.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time, func(stake Stake) bool {
		logger.Info(fmt.Sprintf("Processing expired stakeID %d argumentID %d", stake.ID, stake.ArgumentID))

		// to evaluate if staking limit is increased later on
		if _, exists := initialMultiples[stake.Creator.String()]; !exists {
			initialMultiples[stake.Creator.String()] = getStakeLimitUpgradeMultiple(keeper.TotalEarnedCoins(ctx, stake.Creator))
		}
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
			EventTypeInterestRewardPaid,
			sdk.NewAttribute(AttributeKeyExpiredStakes, string(b)),
		),
	)

	// evaluating whether anybody's staking limit should be increased
	for address, initial := range initialMultiples {
		creator, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			panic(err)
		}
		earned := keeper.TotalEarnedCoins(ctx, creator)
		final := getStakeLimitUpgradeMultiple(earned)
		if final-initial == 0 {
			// no change in limits
			continue
		}

		limit := getUpgradedStakeLimit(final)
		logger.Info(fmt.Sprintf("Upgraded the staking limit for %s to %d", address, limit))
		upgradeBz, err := keeper.codec.MarshalJSON(StakeLimitUpgrade{
			Address:     creator,
			NewLimit:    limit,
			EarnedStake: sdk.NewCoin(chain.StakeDenom, earned),
		})
		if err != nil {
			panic(err)
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeStakeLimitIncreased,
				sdk.NewAttribute(AttributeKeyStakeLimitUpgrade, string(upgradeBz)),
			),
		)
	}
}

func getStakeLimitUpgradeMultiple(earned sdk.Int) int {
	/**
		if earned stake is n,
		we'll compute the multiple, x, using:
		n/10 = x

		eg.
		10/10 = 1
		20/10 = 2

		so that,
		limit = 1000 + (500 * (x-1))

		eg.
		for n = 10, limit = 1000 + (500 * 0) = 1000
		for n = 20, limit = 1000 + (500 * 1) = 1500
		for n = 30, limit = 1000 + (500 * 2) = 2000
	**/

	if earned.IsNegative() {
		return 0
	}

	divisor := int64(math.Pow10(10)) // 10^10 (10^9 for removing the decimal part; 10^1 for finding the multiple)
	return int(earned.QuoRaw(divisor).Int64())
}

func getUpgradedStakeLimit(multiple int) int {
	return 1000 + (500 * (multiple - 1))
}
