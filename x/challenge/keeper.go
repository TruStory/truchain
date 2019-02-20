package challenge

import (
	"fmt"

	"github.com/TruStory/truchain/x/backing"

	"github.com/cosmos/cosmos-sdk/x/bank"

	params "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for challenges
	StoreKey = "challenges"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	Challenge(
		ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error)

	ChallengesByStoryID(
		ctx sdk.Context, storyID int64) (challenges []Challenge, err sdk.Error)

	ChallengeByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (challenge Challenge, err sdk.Error)

	Tally(ctx sdk.Context, gameID int64) (falseVotes []Challenge, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
		creator sdk.AccAddress) (int64, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	// list of games waiting to meet challenge threshold
	pendingGameListKey sdk.StoreKey

	backingKeeper backing.ReadKeeper
	bankKeeper    bank.Keeper
	storyKeeper   story.WriteKeeper

	challengeList app.UserList // challenge <-> story mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	pendingGameListKey sdk.StoreKey,
	backingKeeper backing.ReadKeeper,
	bankKeeper bank.Keeper,
	storyKeeper story.WriteKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		pendingGameListKey,
		backingKeeper,
		bankKeeper,
		storyKeeper,
		app.NewUserList(storyKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new challenge on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
	creator sdk.AccAddress) (challengeID int64, err sdk.Error) {

	logger := ctx.Logger().With("module", "challenge")

	if amount.Denom != params.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid backing token.")
	}

	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for challenge.")
	}

	// check if challenge amount is greater than minimum stake
	if amount.Amount.LT(game.DefaultParams().MinChallengeStake) {
		return 0, sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	// make sure creator hasn't already challenged
	if k.challengeList.Includes(ctx, k, storyID, creator) {
		return 0, ErrDuplicateChallenge(storyID, creator)
	}

	// create implicit false vote
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		StoryID:   storyID,
		Amount:    amount,
		Argument:  argument,
		Creator:   creator,
		Vote:      false,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new challenge with embedded vote
	challenge := Challenge{vote}

	// persist challenge
	k.GetStore(ctx).Set(
		k.GetIDKey(challenge.ID()),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(challenge))

	// persist challenge <-> story mapping
	k.challengeList.Append(ctx, k, storyID, creator, challenge.ID())

	// deduct challenge amount from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return 0, err
	}

	msg := fmt.Sprintf("Challenged story %d with %s by %s",
		storyID, amount.String(), creator.String())
	logger.Info(msg)

	return challenge.ID(), nil
}

// Challenge gets the challenge for the given id
func (k Keeper) Challenge(
	ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error) {

	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(challengeID))
	if bz == nil {
		return challenge, ErrNotFound(challengeID)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &challenge)

	return
}

// ChallengesByStoryID returns the list of challenges for a story id
func (k Keeper) ChallengesByStoryID(
	ctx sdk.Context, storyID int64) (challenges []Challenge, err sdk.Error) {

	// iterate over and return challenges for a game
	err = k.challengeList.Map(ctx, k, storyID, func(challengeID int64) sdk.Error {
		challenge, err := k.Challenge(ctx, challengeID)
		if err != nil {
			return err
		}
		challenges = append(challenges, challenge)

		return nil
	})

	return
}

// ChallengeByStoryIDAndCreator returns a challenge for a given story id and creator
func (k Keeper) ChallengeByStoryIDAndCreator(
	ctx sdk.Context,
	storyID int64,
	creator sdk.AccAddress) (challenge Challenge, err sdk.Error) {

	// get the story
	s, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return challenge, story.ErrInvalidStoryID(storyID)
	}

	// get the challenge
	challengeID := k.challengeList.Get(ctx, k, s.ID, creator)
	challenge, err = k.Challenge(ctx, challengeID)

	return
}

// Tally challenges for voting
func (k Keeper) Tally(
	ctx sdk.Context, storyID int64) (falseVotes []Challenge, err sdk.Error) {

	err = k.challengeList.Map(ctx, k, storyID, func(challengeID int64) sdk.Error {
		challenge, err := k.Challenge(ctx, challengeID)
		if err != nil {
			return err
		}

		if challenge.VoteChoice() == true {
			return ErrInvalidVote()
		}
		falseVotes = append(falseVotes, challenge)

		return nil
	})

	return
}

// ============================================================================

func (k Keeper) pendingGameList(ctx sdk.Context) list.List {
	return list.NewList(
		k.GetCodec(),
		ctx.KVStore(k.pendingGameListKey))
}
