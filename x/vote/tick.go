package vote

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := k.filterGameQueue(ctx, k.gameQueue(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// filterGameQueue checks to see if a validation game has ended, then processes
// that game. It calls itself recursively until all games have been processed.
func (k Keeper) filterGameQueue(ctx sdk.Context, gameQueue queue.Queue) sdk.Error {
	logger := ctx.Logger().With("module", "vote")

	if gameQueue.IsEmpty() {
		return nil
	}

	var gameID int64
	if err := gameQueue.Peek(&gameID); err != nil {
		panic(err)
	}

	// retrieve the game
	game, err := k.gameKeeper.Game(ctx, gameID)
	if err != nil {
		return err
	}

	blockTime := ctx.BlockHeader().Time

	quorum, err := k.quorum(ctx, game.StoryID)
	if err != nil {
		return err
	}

	// handle expired voting periods
	if game.IsVotingExpired(blockTime, quorum) {

		logger.Info(
			fmt.Sprintf(
				"Voting period expired for story: %d", game.StoryID))

		// remove from queue
		gameQueue.Pop()

		// return funds
		err = k.returnFunds(ctx, gameID)
		if err != nil {
			return err
		}

		// update story
		err = k.storyKeeper.ExpireGame(ctx, game.StoryID)
		if err != nil {
			return err
		}

		// process next game
		return k.filterGameQueue(ctx, gameQueue)
	}

	// Terminate recursion on finding the first unfinished game,
	// because it means all the ones behind it in the queue
	// are also unfinished.
	if !game.IsVotingFinished(blockTime, quorum) {
		return nil
	}

	// only left with finished games at this point...
	// remove finished game from queue
	gameQueue.Pop()

	// process game
	err = processGame(ctx, k, game)
	if err != nil {
		return err
	}

	// check next game
	return k.filterGameQueue(ctx, gameQueue)
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

	challenges, err := k.challengeKeeper.ChallengesByGameID(ctx, story.ID)
	if err != nil {
		return
	}

	tokenVotes, err := k.TokenVotesByGameID(ctx, story.ID)
	if err != nil {
		return
	}

	total = len(backings) + len(challenges) + len(tokenVotes)

	return total, nil
}

func (k Keeper) returnFunds(ctx sdk.Context, gameID int64) sdk.Error {
	logger := ctx.Logger().With("module", "vote")

	// get challenges
	challenges, err := k.challengeKeeper.ChallengesByGameID(ctx, gameID)
	if err != nil {
		return err
	}

	// get token votes
	tokenVotes, err := k.TokenVotesByGameID(ctx, gameID)
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
		"Returned funds for %d users for game %d", len(votes), gameID))

	return nil
}
