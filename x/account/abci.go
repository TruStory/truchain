package account

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	toUnjail, err := keeper.JailedAccountsAfter(ctx, ctx.BlockHeader().Time)
	if err != nil {
		panic(err)
	}
	for _, acct := range toUnjail {
		err = keeper.UnJail(ctx, acct.GetAddress())
		if err != nil {
			panic(err)
		}
		logger(ctx).Info(fmt.Sprintf("Unjailed %s", acct.String()))
	}

	return sdk.EmptyTags()
}