package vote

import (
	app "github.com/TruStory/truchain/types"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	store := ctx.KVStore(k.activeGamesQueueKey)
	q := queue.NewQueue(k.GetCodec(), store)

	err := k.checkGames(ctx, q)
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// checkGames checks to see if a validation game has ended.
// It calls itself recursively until all games have been processed.
func (k Keeper) checkGames(ctx sdk.Context, gameQueue queue.Queue) sdk.Error {
	// check the head of the queue
	var gameID int64
	if err := gameQueue.Peek(&gameID); err != nil {
		return nil
	}

	// retrieve the game
	game, err := k.gameKeeper.Game(ctx, gameID)
	if err != nil {
		return err
	}

	blockTime := ctx.BlockHeader().Time

	// handle expired games
	// an expired game meets the following criteria:
	// 1. passed the voting period (`EndTime` > block time)
	// 2. didn't meet the minimum voter quorum
	if game.Expired(blockTime) {
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
		return k.checkGames(ctx, gameQueue)
	}

	// terminate recursion on finding the first non-ended game
	// an ended game meets the following criteria:
	// 1. passed the voting period (`EndTime` > block time)
	// 2. met the minimum voter quorum
	if !game.Ended(blockTime) {
		return nil
	}

	// remove ended game from queue
	gameQueue.Pop()

	// process ended game
	err = processGame(ctx, k, game)
	if err != nil {
		return err
	}

	// check next game
	return k.checkGames(ctx, gameQueue)
}

func (k Keeper) returnFunds(ctx sdk.Context, gameID int64) sdk.Error {
	// get challenges
	challenges, err := k.challengeKeeper.ChallengesByGameID(ctx, gameID)
	if err != nil {
		return err
	}

	// get token votes
	tokenVotes, err := k.TokenVotesByGame(ctx, gameID)
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

	return nil
}
