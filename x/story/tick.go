package story

import (
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processStoryQueue(ctx, k.storyQueue(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// processStoryQueue checks to see if a story has expired
func (k Keeper) processStoryQueue(ctx sdk.Context, storyQueue queue.Queue) sdk.Error {
	logger := ctx.Logger().With("module", "story")

	if storyQueue.IsEmpty() {
		return nil
	}

	var storyID int64
	if err := storyQueue.Peek(&storyID); err != nil {
		panic(err)
	}

	story, err := k.Story(ctx, storyID)
	if err != nil {
		return err
	}

	logger.Info("Processing " + story.String())

	if ctx.BlockHeader().Time.After(story.ExpireTime) {

	}

	// 	logger.Info(
	// 		fmt.Sprintf(
	// 			"Voting period expired for story: %d", game.StoryID))

	// 	// remove from queue
	// 	gameQueue.Pop()

	// 	// return funds
	// 	err = k.returnFunds(ctx, gameID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// update story
	// 	err = k.storyKeeper.ExpireGame(ctx, game.StoryID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// process next game
	// 	return k.filterGameQueue(ctx, gameQueue)
	// }

	// // Terminate recursion on finding the first unfinished game,
	// // because it means all the ones behind it in the queue
	// // are also unfinished.
	// if !game.IsVotingFinished(blockTime, quorum) {
	// 	return nil
	// }

	// // only left with finished games at this point...
	// // remove finished game from queue
	// gameQueue.Pop()

	// // process game
	// err = processGame(ctx, k, game)
	// if err != nil {
	// 	return err
	// }

	// check next story
	return k.processStoryQueue(ctx, storyQueue)
}
