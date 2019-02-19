package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.handleExpiredStories(ctx)
	if err != nil {
		panic(err)
	}
	return sdk.EmptyTags()
}

// recursively process expired stories
func (k Keeper) handleExpiredStories(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", "distribution")

	expiredStoryQueue := k.expiredStoryQueue(ctx)

	if expiredStoryQueue.IsEmpty() {
		// done processing all expired stories
		// terminate
		return nil
	}

	var storyID int64
	if err := expiredStoryQueue.Peek(&storyID); err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("Handling expired story id: %d", storyID))

	expiredStoryQueue.Pop()

	err := k.distributeEarningsToBackers(ctx, storyID)
	if err != nil {
		return err
	}

	err = k.returnFundsToChallengers(ctx, storyID)
	if err != nil {
		return err
	}

	// handle next expired story
	return k.handleExpiredStories(ctx)
}

func (k Keeper) distributeEarningsToBackers(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", "distribution")

	backings, err := k.backingKeeper.BackingsByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	for _, backing := range backings {
		// give the principal back to the user (in trustake)
		_, _, err := k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Amount()})
		if err != nil {
			return err
		}

		// give the interest earned to the user (in cred)
		_, _, err = k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Interest})
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf(
			"Distributed earnings of %s with interest of %s to %s",
			backing.Amount().String(),
			backing.Interest.String(),
			backing.Creator().String()))
	}

	return nil
}

func (k Keeper) returnFundsToChallengers(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", "distribution")

	// get challenges
	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	// return funds
	for _, v := range challenges {
		msg := "Returned challenged amount %s back to %s for story %d."
		logger.Info(fmt.Sprintf(msg, v.Amount(), v.Creator(), storyID))

		_, _, err = k.bankKeeper.AddCoins(
			ctx, v.Creator(), sdk.Coins{v.Amount()})
		if err != nil {
			return err
		}
	}

	return nil
}
