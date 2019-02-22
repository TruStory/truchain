package voting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processVotingStoryList(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// Iterate voting story list to see if a validation game has ended
func (k Keeper) processVotingStoryList(ctx sdk.Context) sdk.Error {
	// logger := ctx.Logger().With("module", "vote")

	var storyID int64
	k.votingStoryList(ctx).Iterate(&storyID, func(index uint64) bool {
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

		votingEndTime := story.VotingStartTime.Add(k.votingDuration(ctx))
		if ctx.BlockHeader().Time.Before(votingEndTime) {
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

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return
	}

	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, story.ID)
	if err != nil {
		return
	}

	tokenVotes, err := k.voteKeeper.TokenVotesByStoryID(ctx, story.ID)
	if err != nil {
		return
	}

	total = len(backings) + len(challenges) + len(tokenVotes)

	return total, nil
}

// func (k Keeper) returnFunds(ctx sdk.Context, gameID int64) sdk.Error {
// 	logger := ctx.Logger().With("module", "vote")

// 	// get challenges
// 	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, gameID)
// 	if err != nil {
// 		return err
// 	}

// 	// get token votes
// 	tokenVotes, err := k.TokenVotesByStoryID(ctx, gameID)
// 	if err != nil {
// 		return err
// 	}

// 	// collate votes
// 	var votes []app.Voter
// 	for _, v := range challenges {
// 		votes = append(votes, v)
// 	}
// 	for _, v := range tokenVotes {
// 		votes = append(votes, v)
// 	}

// 	// return funds
// 	for _, v := range votes {
// 		_, _, err = k.bankKeeper.AddCoins(
// 			ctx, v.Creator(), sdk.Coins{v.Amount()})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	logger.Info(fmt.Sprintf(
// 		"Returned funds for %d users for game %d", len(votes), gameID))

// 	return nil
// }
