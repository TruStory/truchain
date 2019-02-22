package voting

import (
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

// Iterate voting story list to see if a validation game has ended
func (k Keeper) processVotingStoryQueue(ctx sdk.Context) sdk.Error {
	var storyID int64
	k.votingStoryQueue(ctx).List.Iterate(&storyID, func(index uint64) bool {
		quorum, err := k.quorum(ctx, storyID)
		if err != nil {
			panic(err)
		}

		if quorum < k.minQuorum(ctx) {
			// move to next story
			return false
		}

		story, err := k.storyKeeper.Story(ctx, storyID)
		if err != nil {
			panic(err)
		}

		if ctx.BlockHeader().Time.Before(story.VotingEndTime) {
			// move to next story
			return false
		}

		// only left with voting ended + met quorum stories
		err = k.verifyStory(ctx, storyID)
		if err != nil {
			panic(err)
		}

		return false
	})

	return nil
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
