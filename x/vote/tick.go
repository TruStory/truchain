package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/game"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	store := ctx.KVStore(k.activeGamesQueueKey)
	q := queue.NewQueue(k.GetCodec(), store)

	err := checkGames(ctx, k, q)
	if err != nil {
		panic(err)
	}

	// TODO: maybe tags should return err?

	return sdk.NewTags()
}

// ============================================================================

// checkGames checks to see if a validation game has ended.
// It calls itself recursively until all games have been processed.
func checkGames(ctx sdk.Context, k Keeper, q queue.Queue) (err sdk.Error) {
	// check the head of the queue
	var gameID int64
	if err := q.Peek(&gameID); err != nil {
		return nil
	}

	// retrieve the game
	game, err := k.gameKeeper.Get(ctx, gameID)
	if err != nil {
		return err
	}

	// terminate recursion on finding the first non-ended game
	if game.Ended(ctx.BlockHeader().Time) {
		return nil
	}

	// remove ended game from queue
	q.Pop()

	// tally backings, challenges, and votes
	trueVotes, falseVotes, err := tally(ctx, k, game)
	if err != nil {
		return err
	}

	// update reward pool, return funds to losers, reward winners
	if confirmStory(trueVotes, falseVotes) {
		err = confirmedStoryRewardPool(ctx, k, game, falseVotes)
		if err != nil {
			return err
		}

		err = returnFunds(ctx, k, game, falseVotes)
		if err != nil {
			return err
		}

		err = rewardWinners(ctx, k, game, trueVotes)
		if err != nil {
			return err
		}
	} else {
		err = rejectedStoryRewardPool(ctx, k, game, trueVotes, falseVotes)
		if err != nil {
			return err
		}

		err = returnFunds(ctx, k, game, trueVotes)
		if err != nil {
			return err
		}

		err = rewardWinners(ctx, k, game, falseVotes)
		if err != nil {
			return err
		}
	}

	return checkGames(ctx, k, q)
}

func tally(
	ctx sdk.Context,
	k Keeper,
	game game.Game) (trueVotes []interface{}, falseVotes []interface{}, err sdk.Error) {

	// tally backings
	yesBackings, noBackings, err := k.backingKeeper.Tally(ctx, game.StoryID)
	if err != nil {
		return
	}
	trueVotes = append(trueVotes, yesBackings)
	falseVotes = append(falseVotes, noBackings)

	// tally challenges
	yesChallenges, noChallenges, err := k.challengeKeeper.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	trueVotes = append(trueVotes, yesChallenges)
	falseVotes = append(falseVotes, noChallenges)

	// tally votes
	trueTokenVotes, falseTokenVotes, err := k.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	trueVotes = append(trueVotes, trueTokenVotes)
	falseVotes = append(falseVotes, falseTokenVotes)

	return
}

// determine if a story is confirmed or rejected
func confirmStory(trueVotes []interface{}, falseVotes []interface{}) (confirmed bool) {
	// calculate weighted votes
	trueWeight := calculateWeightedVote(trueVotes)
	falseWeight := calculateWeightedVote(falseVotes)

	if trueWeight.GT(falseWeight) {
		// story confirmed
		return true
	}

	// story rejected
	return false
}

// calculate weighted vote based on user's total category coin balance
func calculateWeightedVote(poll []interface{}) sdk.Int {
	weightedAmount := sdk.ZeroInt()
	for _, vote := range poll {
		v := vote.(app.Vote)
		user := auth.NewBaseAccountWithAddress(v.Creator)
		categoryCoins := user.Coins.AmountOf(v.Amount.Denom)
		weightedAmount = weightedAmount.Add(categoryCoins)
	}

	return weightedAmount
}

// people who voted no on a confirmed story
func confirmedStoryRewardPool(
	ctx sdk.Context,
	k Keeper,
	game game.Game,
	no []interface{}) (err sdk.Error) {

	for _, vote := range no {
		switch v := vote.(type) {
		case backing.Backing:
			err = handleConfirmedStoryNoVoteBacker(ctx, k, v, game)
		case challenge.Challenge:
			// skip
			// already added amount to reward pool, lost funds
		case app.Vote:
			// skip
			// already added amount to reward pool, lost funds
		default:
			return ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func rejectedStoryRewardPool(
	ctx sdk.Context,
	k Keeper,
	game game.Game,
	yes []interface{},
	no []interface{}) (err sdk.Error) {

	for _, vote := range yes {
		switch v := vote.(type) {
		case backing.Backing:
			err = handleRejectedStoryYesVoteBacker(ctx, k, v, game)
		case app.Vote:
			err = handleRejectedStoryYesVoteVoter(ctx, k, v, game)
		default:
			err = ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

	for _, vote := range no {
		switch v := vote.(type) {
		case backing.Backing:
			err = handleRejectedStoryNoVoteBacker(ctx, k, v, game)
		case challenge.Challenge:
			// skip
			// already added amount to reward pool
		case app.Vote:
			// skip
			// already added vote fee to reward pool
		default:
			return ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

	return
}

// return funds to losers
func returnFunds(
	ctx sdk.Context, k Keeper, game game.Game, losers []interface{}) (err sdk.Error) {

	for _, vote := range losers {
		switch v := vote.(type) {
		case backing.Backing:
			// return backing amount to backer
			_, _, err = k.bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{v.Amount})
		case challenge.Challenge:
			// skip
		case app.Vote:
			// skip
		default:
			return ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// reward winners
func rewardWinners(
	ctx sdk.Context, k Keeper, game game.Game, winners []interface{}) (err sdk.Error) {

	// divide reward pool equally
	numWinners := int64(len(winners))
	rewardAmount := game.Pool.Amount.Div(sdk.NewInt(numWinners))
	rewardCoin := sdk.NewCoin(game.Pool.Denom, rewardAmount)

	// distribute reward
	for _, vote := range winners {
		v := vote.(app.Vote)
		_, _, err = k.bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{rewardCoin})
		if err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================

// backer who changed their implicit TRUE vote to FALSE and lost
func handleConfirmedStoryNoVoteBacker(
	ctx sdk.Context, k Keeper, b backing.Backing, game game.Game) sdk.Error {

	// TODO: shouldn't this be from a "backing pool"?
	// remove backing amount from reward pool
	// game.Pool = game.Pool.Minus(b.Amount)

	// return backing amount to backer
	_, _, err := k.bankKeeper.AddCoins(ctx, b.Creator, sdk.Coins{b.Amount})
	if err != nil {
		return err
	}

	// slash inflationary rewards and add to reward pool
	game.Pool = game.Pool.Plus(b.Interest)

	// persist changes to reward pool
	k.gameKeeper.Set(ctx, game)

	return nil
}

// ============================================================================

func handleRejectedStoryYesVoteBacker(
	ctx sdk.Context, k Keeper, b backing.Backing, game game.Game) sdk.Error {

	// forfeit backing and inflationary rewards, add to reward pool
	game.Pool = game.Pool.Plus(b.Amount).Plus(b.Interest)

	// persist changes to reward pool
	k.gameKeeper.Set(ctx, game)

	return nil
}

// token holders who voted TRUE
func handleRejectedStoryYesVoteVoter(
	ctx sdk.Context, k Keeper, v app.Vote, game game.Game) sdk.Error {

	// forfeit vote fee and add to reward pool
	game.Pool = game.Pool.Plus(v.Amount)

	// persist changes to reward pool
	k.gameKeeper.Set(ctx, game)

	return nil
}

// backer who changed their implicit TRUE vote to FALSE, and lost
func handleRejectedStoryNoVoteBacker(
	ctx sdk.Context, k Keeper, b backing.Backing, game game.Game) sdk.Error {

	// return backing
	_, _, err := k.bankKeeper.AddCoins(ctx, b.Creator, sdk.Coins{b.Amount})
	if err != nil {
		return err
	}

	// slash inflationary reward and add to reward pool
	game.Pool = game.Pool.Plus(b.Interest)

	// persist changes to reward pool
	k.gameKeeper.Set(ctx, game)

	return nil
}
