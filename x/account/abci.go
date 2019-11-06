package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) {
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
