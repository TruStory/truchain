package account

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.distributeInflation(ctx)
	keeper.unjailAccounts(ctx)

	acc := keeper.supplyKeeper.GetModuleAccount(ctx, UserGrowthPoolName)
	fmt.Println(acc.GetName() + " " + acc.GetCoins().String())

	acc1 := keeper.supplyKeeper.GetModuleAccount(ctx, StakeholderPoolName)
	fmt.Println(acc1.GetName() + " " + acc1.GetCoins().String())
}

func (k Keeper) unjailAccounts(ctx sdk.Context) {
	toUnjail, err := k.JailedAccountsBefore(ctx, ctx.BlockHeader().Time)
	if err != nil {
		panic(err)
	}

	for _, acct := range toUnjail {
		err = k.UnJail(ctx, acct.PrimaryAddress())
		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeUnjailedAccount,
				sdk.NewAttribute(AttributeKeyUser, acct.PrimaryAddress().String()),
			),
		)

		k.Logger(ctx).Info(fmt.Sprintf("Unjailed %s", acct.String()))
	}
}

func (k Keeper) distributeInflation(ctx sdk.Context) {
	k.distributeInflationToUserGrowthPool(ctx)
	k.distributeInflationToStakeholderPool(ctx)
}

func (k Keeper) distributeInflationToUserGrowthPool(ctx sdk.Context) {
	userGrowthAllocation := k.GetParams(ctx).UserGrowthAllocation
	//fmt.Println("user growth allocation " + userGrowthAllocation.String())

	inflationAcc := k.supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName)
	inflation := inflationAcc.GetCoins()
	//fmt.Println("inflation " + inflation.String())

	inflationDec := sdk.NewDecFromInt(inflation.AmountOf("tru"))
	//fmt.Println("inflation dec " + inflationDec.String())

	userGrowthAmount := inflationDec.Mul(userGrowthAllocation).RoundInt()
	//fmt.Println("user growth amount rounded int " + userGrowthAmount.String())

	userGrowthCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userGrowthAmount))
	//fmt.Println("user growth coins " + userGrowthCoins.String())

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, distribution.ModuleName, UserGrowthPoolName, userGrowthCoins)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) distributeInflationToStakeholderPool(ctx sdk.Context) {
	stakeholderAllocation := k.GetParams(ctx).StakeholderAllocation
	inflationAcc := k.supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName)
	stakeholderInflation := inflationAcc.GetCoins().AmountOf(app.StakeDenom)
	stakeholderInflationDec := sdk.NewDecFromInt(stakeholderInflation)
	stakeholderAmount := stakeholderInflationDec.Mul(stakeholderAllocation).RoundInt()
	stakeholderCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, stakeholderAmount))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, distribution.ModuleName, StakeholderPoolName, stakeholderCoins)
	if err != nil {
		panic(err)
	}
}
