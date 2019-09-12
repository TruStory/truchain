package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	toUnjail, err := keeper.JailedAccountsBefore(ctx, ctx.BlockHeader().Time)
	if err != nil {
		panic(err)
	}

	for _, acct := range toUnjail {
		err = keeper.UnJail(ctx, acct.PrimaryAddress())
		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeUnjailedAccount,
				sdk.NewAttribute(AttributeKeyUser, acct.PrimaryAddress().String()),
			),
		)

		keeper.Logger(ctx).Info(fmt.Sprintf("Unjailed %s", acct.String()))
	}
}

func distributeInflation(ctx sdk.Context, keeper Keeper) sdk.Error {
	//userGrowthAllocation := sdk.NewDecWithPrec(800, 3)
	//
	//amount := sdk.NewCoins(sdk.NewInt64Coin("tru", 23423))
	//err := keeper.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, UserGrowthPoolName, amount)
	//if err != nil {
	//	return err
	//}

	return nil
}
