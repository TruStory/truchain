package expiration

import (
	"fmt"

	"github.com/TruStory/truchain/x/stake"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processExpiredStoryQueue(ctx)
	if err != nil {
		panic(err)
	}
	return sdk.EmptyTags()
}

// processExpiredStoryQueue recursively process expired stories.
// If a story id gets in this queue, it means that it never went
// through a voting period. Therefore, rewards are distributed to
// backers, and funds are returned to challengers. There are no
// voters because the voting period never started.
func (k Keeper) processExpiredStoryQueue(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", "expiration")

	expiringStoryQueue := k.expiringStoryQueue(ctx)

	if expiringStoryQueue.IsEmpty() {
		// done processing all expired stories
		// terminate
		return nil
	}

	var storyID int64
	if err := expiringStoryQueue.Peek(&storyID); err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Handling expired story id: %d", storyID))

	expiringStoryQueue.Pop()

	var votes []stake.Voter

	backings, err := k.backingKeeper.BackingsByStoryID(ctx, storyID)
	if err != nil {
		return err
	}
	for _, backing := range backings {
		votes = append(votes, backing)
	}

	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, storyID)
	if err != nil {
		return err
	}
	for _, challenge := range challenges {
		votes = append(votes, challenge)
	}

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	err = k.stakeKeeper.DistributePrincipalAndInterest(ctx, votes, story.CategoryID)
	if err != nil {
		return err
	}

	// handle next expired story
	return k.processExpiredStoryQueue(ctx)
}
