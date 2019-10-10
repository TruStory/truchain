package distribution

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// BeginBlocker called every block, process expiring stakes
func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	supplyTotal := keeper.supplyKeeper.GetSupply(ctx)
	fmt.Println("supply " + supplyTotal.String())

	distAcc := keeper.supplyKeeper.GetModuleAccount(ctx, "distribution")
	fmt.Println(distAcc.GetName() + " " + distAcc.GetCoins().String())

	communityPool := keeper.cosmosDistKeeper.GetFeePoolCommunityCoins(ctx)
	fmt.Println("community pool " + communityPool.String())

	keeper.distributeInflation(ctx)

	fee := keeper.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	fmt.Println(fee.GetName() + " " + fee.GetCoins().String())

	acc := keeper.supplyKeeper.GetModuleAccount(ctx, UserGrowthPoolName)
	fmt.Println(acc.GetName() + " " + acc.GetCoins().String())
	acc1 := keeper.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName)
	fmt.Println(acc1.GetName() + " " + acc1.GetCoins().String())
	acc2 := keeper.supplyKeeper.GetModuleAccount(ctx, StakeholderPoolName)
	fmt.Println(acc2.GetName() + " " + acc2.GetCoins().String())
	//acc3 := keeper.supplyKeeper.GetModuleAccount(ctx, "user_stakes_tokens_pool")
	//fmt.Println(acc3.GetName() + " " + acc3.GetCoins().String())
}

func (k Keeper) distributeInflation(ctx sdk.Context) {
	// total inflation includes validator + user rewards
	totalInflation := k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(app.StakeDenom)
	// 50% of inflation goes to TruStory pools
	// the rest will go to validators + community when the Cosmos distribution begin blocker runs after this one
	userInflationDec := sdk.NewDecFromInt(totalInflation).QuoInt(sdk.NewInt(2))
	k.distributeInflationToUserGrowthPool(ctx, userInflationDec)
	k.distributeInflationToUserRewardPool(ctx, userInflationDec)
	k.distributeInflationToStakeholderPool(ctx, userInflationDec)
}

func (k Keeper) distributeInflationToUserGrowthPool(ctx sdk.Context, inflation sdk.Dec) {
	userGrowthAllocation := k.GetParams(ctx).UserGrowthAllocation
	userGrowthAmount := inflation.Mul(userGrowthAllocation)
	userGrowthCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userGrowthAmount.TruncateInt()))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, UserGrowthPoolName, userGrowthCoins)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) distributeInflationToUserRewardPool(ctx sdk.Context, inflation sdk.Dec) {
	userRewardAllocation := k.GetParams(ctx).UserRewardAllocation
	userRewardAmount := inflation.Mul(userRewardAllocation)
	userRewardCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userRewardAmount.TruncateInt()))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, UserRewardPoolName, userRewardCoins)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) distributeInflationToStakeholderPool(ctx sdk.Context, inflation sdk.Dec) {
	stakeholderAllocation := k.GetParams(ctx).StakeholderAllocation
	stakeholderAmount := inflation.Mul(stakeholderAllocation)
	stakeholderCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, stakeholderAmount.TruncateInt()))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, StakeholderPoolName, stakeholderCoins)
	if err != nil {
		panic(err)
	}
}
