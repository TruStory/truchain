package vote

import (
	"fmt"

	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/trubank"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

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

	GetParams(ctx sdk.Context) Params

	Tally(ctx sdk.Context, storyID int64) (
		trueVotes []TokenVote, falseVotes []TokenVote, err sdk.Error)

	TokenVote(ctx sdk.Context, id int64) (vote TokenVote, err sdk.Error)
	TokenVotesByStoryID(ctx sdk.Context, storyID int64) ([]TokenVote, sdk.Error)

	TokenVotesByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (vote TokenVote, err sdk.Error)

	TotalVoteAmountByStoryID(ctx sdk.Context, storyID int64) (
		totalCoin sdk.Coin, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin,
		choice bool, argument string, creator sdk.AccAddress) (int64, sdk.Error)
	Update(ctx sdk.Context, vote TokenVote)
	ToggleVote(ctx sdk.Context, storyID int64, amount sdk.Coin, argument string, creator sdk.AccAddress) (int64, sdk.Error)
	SetParams(ctx sdk.Context, params Params)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	votingStoryQueueKey sdk.StoreKey

	stakeKeeper     stake.Keeper
	accountKeeper   auth.AccountKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	storyKeeper     story.WriteKeeper
	bankKeeper      bank.Keeper
	trubankKeeper   trubank.Keeper
	paramStore      params.Subspace

	voterList app.UserList
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	votingStoryQueueKey sdk.StoreKey,
	stakeKeeper stake.Keeper,
	accountKeeper auth.AccountKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	storyKeeper story.WriteKeeper,
	bankKeeper bank.Keeper,
	trubankKeeper trubank.Keeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		votingStoryQueueKey,
		stakeKeeper,
		accountKeeper,
		backingKeeper,
		challengeKeeper,
		storyKeeper,
		bankKeeper,
		trubankKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
		app.NewUserList(storyKeeper.GetStoreKey()),
	}
}

// ToggleVote toggles a vote for a given story id and account address
func (k Keeper) ToggleVote(ctx sdk.Context, storyID int64, amount sdk.Coin, argument string, creator sdk.AccAddress) (int64, sdk.Error) {
	logger := ctx.Logger().With("module", StoreKey).With("storyID", storyID)
	s, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return 0, err
	}

	if s.Status != story.Challenged {
		return 0, ErrInvalidStoryState(s.Status.String())
	}

	// Check if user backed
	b, err := k.backingKeeper.BackingByStoryIDAndCreator(ctx, storyID, creator)
	if err != nil && err.Code() != backing.CodeNotFound {
		return 0, err
	}
	if b.Vote != nil && err == nil {
		logger.Info("Toggling backing vote to challenge vote")
		err = k.backingKeeper.Delete(ctx, b)
		if err != nil {
			return 0, err
		}
		id, err := k.challengeKeeper.Create(ctx, storyID, b.Amount(), b.Argument, creator, true)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	// Check if user challenged
	c, err := k.challengeKeeper.ChallengeByStoryIDAndCreator(ctx, storyID, creator)
	// if err is different than not found
	if err != nil && err.Code() != challenge.CodeNotFound {
		return 0, err
	}
	if c.Vote != nil && err == nil {
		logger.Info("Toggling challenge vote to backing vote")
		err = k.challengeKeeper.Delete(ctx, c)
		if err != nil {
			return 0, err
		}
		id, err := k.backingKeeper.Create(ctx, storyID, c.Amount(), c.Argument, creator, true)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	// Check if user has a token vote and toggle the vote value
	tv, err := k.TokenVotesByStoryIDAndCreator(ctx, storyID, creator)
	if err != nil {
		return 0, err
	}
	if tv.Vote != nil {
		choice := tv.VoteChoice()
		logger.Info(fmt.Sprintf("toggling token vote from %T fo %T", choice, !choice))
		tv.Vote.Vote = !choice
		// Update the time when the vote was toggled
		tv.Vote.Timestamp = app.NewTimestamp(ctx.BlockHeader())
		k.Update(ctx, tv)
	}
	return tv.ID(), nil

}

// ============================================================================

// Create adds a new vote on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin,
	choice bool, argument string, creator sdk.AccAddress) (int64, sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	// get the story
	currentStory, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// make sure voting has started
	if currentStory.Status != story.Challenged {
		return 0, ErrVotingNotStarted(storyID)
	}

	err = k.stakeKeeper.ValidateArgument(ctx, argument)
	if err != nil {
		return 0, err
	}

	minStake := k.GetParams(ctx).StakeAmount
	if amount.IsLT(minStake) {
		return 0, sdk.ErrInsufficientFunds("Below minimum stake.")
	}

	if amount.Denom != app.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid voting token.")
	}

	// check if this voter has already cast a vote
	if k.voterList.Includes(ctx, k, currentStory.ID, creator) {
		return 0, ErrDuplicateVote(currentStory.ID, creator)
	}

	// check if user has the funds
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds to vote on story.")
	}

	// create a new vote
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		Amount:    amount,
		Argument:  argument,
		Weight:    sdk.NewInt(0),
		Creator:   creator,
		Vote:      choice,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	tokenVote := TokenVote{&vote}

	typeOfVote := trubank.Challenge
	if choice {
		typeOfVote = trubank.Backing
	}

	// deduct challenge amount from user
	_, err = k.trubankKeeper.SubtractCoin(ctx, creator, amount, storyID, typeOfVote, tokenVote.ID())
	if err != nil {
		return 0, err
	}

	// persist vote
	k.set(ctx, tokenVote)

	// persist story <-> tokenVote association
	k.voterList.Append(ctx, k, currentStory.ID, creator, vote.ID)

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

// TokenVotesByStoryID returns a list of votes for a given game
func (k Keeper) TokenVotesByStoryID(
	ctx sdk.Context, storyID int64) (votes []TokenVote, err sdk.Error) {

	// iterate over voter list and get votes
	err = k.voterList.Map(ctx, k, storyID, func(voterID int64) sdk.Error {
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
	ctx sdk.Context, storyID int64) (
	trueVotes []TokenVote, falseVotes []TokenVote, err sdk.Error) {

	err = k.voterList.Map(ctx, k, storyID, func(voteID int64) sdk.Error {
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

// TotalVoteAmountByStoryID returns the total of all votes for a game
func (k Keeper) TotalVoteAmountByStoryID(ctx sdk.Context, storyID int64) (
	totalCoin sdk.Coin, err sdk.Error) {

	totalAmount := sdk.ZeroInt()

	votes, err := k.TokenVotesByStoryID(ctx, storyID)
	if err != nil {
		return
	}

	for _, tokenVote := range votes {
		totalAmount = totalAmount.Add(tokenVote.Amount().Amount)
	}

	return sdk.NewCoin(app.StakeDenom, totalAmount), nil
}

// Update updates the vote
func (k Keeper) Update(ctx sdk.Context, vote TokenVote) {

	newVote := TokenVote{
		Vote: vote.Vote,
	}

	k.set(ctx, newVote)
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

func (k Keeper) validateStoryState(ctx sdk.Context, storyID int64) sdk.Error {
	s, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	if s.Status != story.Challenged {
		return ErrInvalidStoryState(s.Status.String())
	}

	return nil
}
