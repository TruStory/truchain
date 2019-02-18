package vote

import (
	"fmt"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	params "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for vote
	StoreKey = "vote"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	Tally(ctx sdk.Context, gameID int64) (
		trueVotes []TokenVote, falseVotes []TokenVote, err sdk.Error)

	TokenVote(ctx sdk.Context, id int64) (vote TokenVote, err sdk.Error)

	TokenVotesByGameID(ctx sdk.Context, gameID int64) ([]TokenVote, sdk.Error)

	TokenVotesByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (vote TokenVote, err sdk.Error)

	TotalVoteAmountByGameID(ctx sdk.Context, gameID int64) (
		totalCoin sdk.Coin, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin,
		choice bool, argument string, creator sdk.AccAddress) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	votingStoryQueueKey sdk.StoreKey

	accountKeeper   auth.AccountKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	storyKeeper     story.WriteKeeper
	gameKeeper      game.WriteKeeper
	bankKeeper      bank.Keeper

	voterList app.UserList
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	votingStoryQueueKey sdk.StoreKey,
	accountKeeper auth.AccountKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	storyKeeper story.WriteKeeper,
	gameKeeper game.WriteKeeper,
	bankKeeper bank.Keeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		votingStoryQueueKey,
		accountKeeper,
		backingKeeper,
		challengeKeeper,
		storyKeeper,
		gameKeeper,
		bankKeeper,
		app.NewUserList(gameKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new vote on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin,
	choice bool, argument string, creator sdk.AccAddress) (int64, sdk.Error) {

	logger := ctx.Logger().With("module", "vote")

	if amount.Denom != params.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid voting token.")
	}

	// get the story
	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// make sure validation game has started
	// if story.GameID <= 0 {
	// 	return 0, ErrGameNotStarted(storyID)
	// }

	// check if this voter has already cast a vote
	if k.voterList.Includes(ctx, k, story.ID, creator) {
		return 0, ErrDuplicateVoteForGame(story.ID, creator)
	}

	// check if user has the funds
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds to vote on story.")
	}

	// deduct vote fee from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return 0, err
	}

	// create a new vote
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		Amount:    amount,
		Argument:  argument,
		Creator:   creator,
		Vote:      choice,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	tokenVote := TokenVote{vote}

	// persist vote
	k.set(ctx, tokenVote)

	// persist story <-> tokenVote association
	k.voterList.Append(ctx, k, story.ID, creator, vote.ID)

	logger.Info(fmt.Sprintf(
		"Voted on story %d with %s by %s", storyID, amount.String(), creator.String()))

	return vote.ID, nil
}

// TokenVote returns a `TokenVote` from the KVStore
func (k Keeper) TokenVote(ctx sdk.Context, id int64) (vote TokenVote, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	if bz == nil {
		return vote, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &vote)

	return vote, nil
}

// TokenVotesByGameID returns a list of votes for a given game
func (k Keeper) TokenVotesByGameID(
	ctx sdk.Context, gameID int64) (votes []TokenVote, err sdk.Error) {

	// iterate over voter list and get votes
	err = k.voterList.Map(ctx, k, gameID, func(voterID int64) sdk.Error {
		vote, err := k.TokenVote(ctx, voterID)
		if err != nil {
			return err
		}
		votes = append(votes, vote)

		return nil
	})

	return
}

// TokenVotesByStoryIDAndCreator returns a vote for a given story id and creator
func (k Keeper) TokenVotesByStoryIDAndCreator(
	ctx sdk.Context,
	storyID int64,
	creator sdk.AccAddress) (vote TokenVote, err sdk.Error) {

	// get the story
	s, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return vote, story.ErrInvalidStoryID(storyID)
	}

	// get the vote
	tokenVoteID := k.voterList.Get(ctx, k, s.ID, creator)
	vote, err = k.TokenVote(ctx, tokenVoteID)

	return
}

// Tally votes
func (k Keeper) Tally(
	ctx sdk.Context, gameID int64) (
	trueVotes []TokenVote, falseVotes []TokenVote, err sdk.Error) {

	err = k.voterList.Map(ctx, k, gameID, func(voteID int64) sdk.Error {
		tokenVote, err := k.TokenVote(ctx, voteID)
		if err != nil {
			return err
		}

		if tokenVote.VoteChoice() == true {
			trueVotes = append(trueVotes, tokenVote)
		} else {
			falseVotes = append(falseVotes, tokenVote)
		}

		return nil
	})

	return
}

// TotalVoteAmountByGameID returns the total of all votes for a game
func (k Keeper) TotalVoteAmountByGameID(ctx sdk.Context, gameID int64) (
	totalCoin sdk.Coin, err sdk.Error) {

	totalAmount := sdk.ZeroInt()

	votes, err := k.TokenVotesByGameID(ctx, gameID)
	if err != nil {
		return
	}

	for _, tokenVote := range votes {
		totalAmount = totalAmount.Add(tokenVote.Amount().Amount)
	}

	return sdk.NewCoin(params.StakeDenom, totalAmount), nil
}

// ============================================================================

func (k Keeper) gameQueue(ctx sdk.Context) queue.Queue {
	store := ctx.KVStore(k.votingStoryQueueKey)
	return queue.NewQueue(k.GetCodec(), store)
}

// saves a `Vote` in the KVStore
func (k Keeper) set(ctx sdk.Context, vote TokenVote) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(vote.ID()),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(vote))
}
