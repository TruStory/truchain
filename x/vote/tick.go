package vote

import (
	"github.com/TruStory/truchain/x/game"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	confirmed, yes, no, err := tally(ctx, k, game)
	if err != nil {
		return err
	}

	if confirmed {
		err = distributeConfirmedCase(ctx, k, yes, no)
	} else {
		err = distributeRejectedCase(ctx, k, yes, no)
	}
	if err != nil {
		return err
	}

	return checkGames(ctx, k, q)
}

func distributeConfirmedCase(
	ctx sdk.Context, k Keeper, yes []interface{}, no []interface{}) (err sdk.Error) {

	// for _, vote := range votes {
	// 	switch vote.(type) {
	// 	case backing.Backing:
	// 		//
	// 	case challenge.Challenge:
	// 		//
	// 	case app.Vote:
	// 		//
	// 	default:
	// 		return sdk.ErrInternal("invalid type")
	// 	}
	// }

	return
}

func distributeRejectedCase(
	ctx sdk.Context, k Keeper, yes []interface{}, no []interface{}) (err sdk.Error) {

	// for _, vote := range votes {
	// 	switch vote.(type) {
	// 	case backing.Backing:
	// 		//
	// 	case challenge.Challenge:
	// 		//
	// 	case app.Vote:
	// 		//
	// 	default:
	// 		return sdk.ErrInternal("invalid type")
	// 	}
	// }

	return
}

func tally(
	ctx sdk.Context,
	k Keeper,
	game game.Game) (confirmed bool, yes []interface{}, no []interface{}, err sdk.Error) {

	// TODO: win is votes weighted by amount

	// tally backings
	yesBackings, noBackings, err := k.backingKeeper.Tally(ctx, game.StoryID)
	if err != nil {
		return
	}
	yes = append(yes, yesBackings)
	no = append(no, noBackings)

	// tally challenges
	yesChallenges, noChallenges, err := k.challengeKeeper.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	yes = append(yes, yesChallenges)
	no = append(no, noChallenges)

	// tally votes
	yesVotes, noVotes, err := k.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	yes = append(yes, yesVotes)
	no = append(no, noVotes)

	numYesVotes := len(yes)
	numNoVotes := len(no)

	if numYesVotes > numNoVotes {
		// story confirmed
		return true, yes, no, nil
	}

	// story rejected
	return false, yes, no, nil
}
