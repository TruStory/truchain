package expiration

import (
	"encoding/json"
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	completed, err := k.processStoryQueue(ctx, make([]app.CompletedStory, 0))
	if err != nil {
		panic(err)
	}
	result := &app.CompletedStoriesNotificationResult{
		Stories: completed,
	}
	b, mErr := json.Marshal(result)
	if mErr != nil {
		panic(mErr)
	}
	if len(completed) == 0 {
		return sdk.EmptyTags()
	}
	return append(app.PushTag, sdk.NewTags(app.KeyCompletedStoriesTag, b)...)
}

func (k Keeper) processStoryQueue(ctx sdk.Context, completed []app.CompletedStory) ([]app.CompletedStory, sdk.Error) {
	logger := ctx.Logger().With("module", StoreKey)

	storyQueue := k.storyQueue(ctx)

	if storyQueue.IsEmpty() {
		// done processing all expired stories
		// terminate
		return completed, nil
	}

	var storyID int64
	if err := storyQueue.Peek(&storyID); err != nil {
		panic(err)
	}

	currentStory, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return completed, err
	}

	logger.Info(fmt.Sprintf("Checking %s", currentStory))

	if ctx.BlockHeader().Time.Before(currentStory.ExpireTime) {
		// return and wait until next block to check if story has expired
		return completed, nil
	}

	logger.Info(fmt.Sprintf("Handling expired story id: %d", storyID))

	storyQueue.Pop()

	var votes []stake.Voter

	backings, err := k.backingKeeper.BackingsByStoryID(ctx, storyID)
	if err != nil {
		return completed, err
	}
	for _, backing := range backings {
		votes = append(votes, backing)
	}

	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, storyID)
	if err != nil {
		return nil, err
	}
	for _, challenge := range challenges {
		votes = append(votes, challenge)
	}

	if len(votes) > 0 {
		err = k.stakeKeeper.RedistributeStake(ctx, votes)
		if err != nil {
			return completed, err
		}

		err = k.stakeKeeper.DistributeInterest(ctx, votes)
		if err != nil {
			return completed, err
		}
	}

	currentStory.Status = story.Expired
	k.storyKeeper.UpdateStory(ctx, currentStory)

	backers, err := k.backingKeeper.BackersByStoryID(ctx, currentStory.ID)
	if err != nil {
		return completed, err
	}

	challengers, err := k.challengeKeeper.ChallengersByStoryID(ctx, currentStory.ID)
	if err != nil {
		return completed, err
	}
	completed = append(completed,
		app.CompletedStory{
			ID:          currentStory.ID,
			Creator:     currentStory.Creator,
			Backers:     backers,
			Challengers: challengers,
		})

	// handle next expired story
	return k.processStoryQueue(ctx, completed)
}
