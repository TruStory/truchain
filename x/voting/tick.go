package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processVotingStoryQueue(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// Recursively process voting story queue to see if voting has ended
func (k Keeper) processVotingStoryQueue(ctx sdk.Context) sdk.Error {
	votingStoryQueue := k.votingStoryQueue(ctx)

	if votingStoryQueue.IsEmpty() {
		return nil
	}

	var storyID int64
	peekErr := votingStoryQueue.Peek(&storyID)
	if peekErr != nil {
		panic(peekErr)
	}

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	if ctx.BlockHeader().Time.Before(story.VotingEndTime) {
		// no stories to process
		// check again after next block
		return nil
	}

	// only left with voting ended stories that may or may not
	// have met the quorum...

	quorum, err := k.quorum(ctx, storyID)
	if err != nil {
		return err
	}

	if quorum < k.minQuorum(ctx) {
		votingStoryQueue.Pop()

		err = k.storyKeeper.EndVotingPeriod(ctx, storyID, false, false)
		if err != nil {
			return err
		}

		err = k.returnFunds(ctx, storyID)
		if err != nil {
			return err
		}

		// process next story
		return k.processVotingStoryQueue(ctx)
	}

	// only left with voting ended + met quorum stories
	votingStoryQueue.Pop()

	err = k.verifyStory(ctx, storyID)
	if err != nil {
		return err
	}

	// process next story
	return k.processVotingStoryQueue(ctx)
}

// quorum returns the total count of backings, challenges, votes
func (k Keeper) quorum(ctx sdk.Context, storyID int64) (total int, err sdk.Error) {
	backings, err := k.backingKeeper.BackingsByStoryID(ctx, storyID)
	if err != nil {
		return
	}

	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, storyID)
	if err != nil {
		return
	}

	tokenVotes, err := k.voteKeeper.TokenVotesByStoryID(ctx, storyID)
	if err != nil {
		return
	}

	total = len(backings) + len(challenges) + len(tokenVotes)

	return total, nil
}

func (k Keeper) returnFunds(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", "voting")

	// TODO: do backers get back their principle too??
	// check expiration time...
	// basically, don't return funds twice...

	// get challenges
	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	// get token votes
	tokenVotes, err := k.voteKeeper.TokenVotesByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	// collate votes
	var votes []app.Voter
	for _, v := range challenges {
		votes = append(votes, v)
	}
	for _, v := range tokenVotes {
		votes = append(votes, v)
	}

	// return funds
	for _, v := range votes {
		_, _, err = k.bankKeeper.AddCoins(
			ctx, v.Creator(), sdk.Coins{v.Amount()})
		if err != nil {
			return err
		}
	}

	logger.Info(fmt.Sprintf(
		"Returned funds for %d users for story %d", len(votes), storyID))

	return nil
}
