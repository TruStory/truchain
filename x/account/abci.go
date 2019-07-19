package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/staking/tags"
)

// EndBlocker called every block, process expiring stakes
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	toUnjail, err := keeper.JailedAccountsAfter(ctx, ctx.BlockHeader().Time)
	if err != nil {
		panic(err)
	}
	unjailed := make([]string, 0)
	for _, acct := range toUnjail {
		err = keeper.UnJail(ctx, acct.GetAddress())
		if err != nil {
			panic(err)
		}
		unjailed = append(unjailed, acct.GetAddress().String())
		logger(ctx).Info(fmt.Sprintf("Unjailed %s", acct.String()))
	}
	if len(unjailed) == 0 {
		return sdk.EmptyTags()
	}
	b, jsonErr := keeper.codec.MarshalJSON(unjailed)
	if jsonErr != nil {
		panic(jsonErr)
	}
	return append(app.PushTag,
		sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Action, tags.ActionInterestRewardPaid,
			tags.ExpiredStakes, b,
		)...,
	)
	return sdk.EmptyTags()
}
