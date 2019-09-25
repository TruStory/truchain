package account

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/supply"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.distributeInflation(ctx)
	keeper.unjailAccounts(ctx)

	totalSupply := keeper.calculateTotalSupply(ctx)
	if keeper.supplyKeeper.GetSupply(ctx).GetTotal().Empty() {
		fmt.Println("supply is empty..")
		keeper.supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))
	}
	fmt.Println(keeper.supplyKeeper.GetSupply(ctx).GetTotal().String())

	acc := keeper.supplyKeeper.GetModuleAccount(ctx, UserGrowthPoolName)
	fmt.Println(acc)

	acc1 := keeper.supplyKeeper.GetModuleAccount(ctx, StakeholderPoolName)
	fmt.Println(acc1)
}

func (k Keeper) calculateTotalSupply(ctx sdk.Context) sdk.Coins {
	var totalSupply sdk.Coins
	k.accountKeeper.IterateAccounts(ctx,
		func(acc auth.Account) (stop bool) {
			totalSupply = totalSupply.Add(acc.GetCoins())
			return false
		},
	)
	return totalSupply
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
	inflationAcc := k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	userInflation := inflationAcc.GetCoins().AmountOf(app.StakeDenom)
	userInflationDec := sdk.NewDecFromIntWithPrec(userInflation, 3)
	userGrowthAmount := userInflationDec.Mul(userGrowthAllocation).RoundInt()
	userGrowthCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userGrowthAmount))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, UserGrowthPoolName, userGrowthCoins)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) distributeInflationToStakeholderPool(ctx sdk.Context) {
	stakeholderAllocation := k.GetParams(ctx).StakeholderAllocation
	inflationAcc := k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	stakeholderInflation := inflationAcc.GetCoins().AmountOf(app.StakeDenom)
	stakeholderInflationDec := sdk.NewDecFromIntWithPrec(stakeholderInflation, 3)
	stakeholderAmount := stakeholderInflationDec.Mul(stakeholderAllocation).RoundInt()
	stakeholderCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, stakeholderAmount))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, StakeholderPoolName, stakeholderCoins)
	if err != nil {
		panic(err)
	}
}
