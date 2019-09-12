package account

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.distributeInflation(ctx)
	keeper.unjailAccounts(ctx)
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
	// TODO: take this from params
	userGrowthAllocation := sdk.NewDecWithPrec(800, 3)

	acc := k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	userInflation := acc.GetCoins().AmountOf(app.StakeDenom)
	userInflationDec := sdk.NewDecFromIntWithPrec(userInflation, 3)
	userGrowthAmount := userInflationDec.Mul(userGrowthAllocation).RoundInt()
	userGrowthCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, userGrowthAmount))

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, UserGrowthPoolName, userGrowthCoins)
	if err != nil {
		panic(err)
	}
}
