package distribution

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker called every block, process expiring stakes
func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	supplyTotal := keeper.supplyKeeper.GetSupply(ctx)
	fmt.Println("supply " + supplyTotal.String())

	//communityPool := keeper.cosmosDistKeeper.GetFeePoolCommunityCoins(ctx)
	//fmt.Println("community pool " + communityPool.String())

	keeper.distributeInflation(ctx)

	acc := keeper.supplyKeeper.GetModuleAccount(ctx, UserGrowthPoolName)
	fmt.Println(acc.GetName() + " " + acc.GetCoins().String())
	acc1 := keeper.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName)
	fmt.Println(acc1.GetName() + " " + acc1.GetCoins().String())
	acc2 := keeper.supplyKeeper.GetModuleAccount(ctx, StakeholderPoolName)
	fmt.Println(acc2.GetName() + " " + acc2.GetCoins().String())
}

func (k Keeper) distributeInflation(ctx sdk.Context) {
	communityPoolAmt := k.cosmosDistKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(app.StakeDenom)
	fmt.Println("community pool before " + communityPoolAmt.String())
	k.distributeInflationToUserGrowthPool(ctx, communityPoolAmt)
	k.distributeInflationToUserRewardPool(ctx, communityPoolAmt)
	k.distributeInflationToStakeholderPool(ctx, communityPoolAmt)
	communityPoolAmt = k.cosmosDistKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(app.StakeDenom)
	fmt.Println("community pool after  " + communityPoolAmt.String())
}

func (k Keeper) distributeInflationToUserGrowthPool(ctx sdk.Context, inflation sdk.Dec) {
	userGrowthAllocation := k.GetParams(ctx).UserGrowthAllocation
	fmt.Println("user growth allocation " + userGrowthAllocation.String())
	userGrowthAmount := inflation.Mul(userGrowthAllocation).TruncateInt()
	userGrowthCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userGrowthAmount))
	userGrowthAcc := k.supplyKeeper.GetModuleAccount(ctx, UserGrowthPoolName)
	err := k.cosmosDistKeeper.DistributeFromFeePool(ctx, userGrowthCoins, userGrowthAcc.GetAddress())
	if err != nil {
		panic(err)
	}
}

func (k Keeper) distributeInflationToUserRewardPool(ctx sdk.Context, inflation sdk.Dec) {
	userRewardAllocation := k.GetParams(ctx).UserRewardAllocation
	fmt.Println("user reward allocation " + userRewardAllocation.String())
	userRewardAmount := inflation.Mul(userRewardAllocation).TruncateInt()
	userRewardCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userRewardAmount))
	userRewardAcc := k.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName)
	err := k.cosmosDistKeeper.DistributeFromFeePool(ctx, userRewardCoins, userRewardAcc.GetAddress())
	if err != nil {
		panic(err)
	}
}

func (k Keeper) distributeInflationToStakeholderPool(ctx sdk.Context, inflation sdk.Dec) {
	stakeholderAllocation := k.GetParams(ctx).StakeholderAllocation
	fmt.Println("stakeholder allocation " + stakeholderAllocation.String())
	stakeholderAmount := inflation.Mul(stakeholderAllocation).TruncateInt()
	stakeholderCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, stakeholderAmount))
	stakeholderAcc := k.supplyKeeper.GetModuleAccount(ctx, StakeholderPoolName)
	err := k.cosmosDistKeeper.DistributeFromFeePool(ctx, stakeholderCoins, stakeholderAcc.GetAddress())
	if err != nil {
		panic(err)
	}
}
