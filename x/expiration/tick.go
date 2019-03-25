package expiration

import (
	"fmt"

	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processStoryQueue(ctx)
	if err != nil {
		panic(err)
	}
	return sdk.EmptyTags()
}

// processStoryQueue recursively process expired stories.
// If a story id gets in this queue, it means that it never went
// through a voting period. Therefore, rewards are distributed to
// backers, and funds are returned to challengers. There are no
// voters because the voting period never started.
func (k Keeper) processStoryQueue(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", StoreKey)

	storyQueue := k.storyQueue(ctx)

	if storyQueue.IsEmpty() {
		// done processing all expired stories
		// terminatex
		return nil
	}

	var storyID int64
	if err := storyQueue.Peek(&storyID); err != nil {
		panic(err)
	}

	currentStory, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Checking %s", currentStory))

	if ctx.BlockHeader().Time.Before(currentStory.ExpireTime) {
		// return and wait until next block to check if story has expired
		return nil
	}

	logger.Info(fmt.Sprintf("Handling expired story id: %d", storyID))

	storyQueue.Pop()

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

	err = k.stakeKeeper.RedistributeStake(ctx, votes)
	if err != nil {
		return err
	}

	err = k.stakeKeeper.DistributeInterest(ctx, votes)
	if err != nil {
		return err
	}

	currentStory.Status = story.Expired
	k.storyKeeper.UpdateStory(ctx, currentStory)

	// handle next expired story
	return k.processStoryQueue(ctx)
}
