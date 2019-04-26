package challenge

import (
	"fmt"

	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/stake"

	"github.com/TruStory/truchain/x/backing"

	trubank "github.com/TruStory/truchain/x/trubank"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
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
	ChallengersByStoryID(
		ctx sdk.Context, storyID int64) (challengers []sdk.AccAddress, err sdk.Error)

	ChallengeByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (challenge Challenge, err sdk.Error)
	GetParams(ctx sdk.Context) Params
	TotalChallengeAmount(ctx sdk.Context, storyID int64) (
		totalCoin sdk.Coin, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context,
		storyID int64,
		amount sdk.Coin,
		argumentID int64,
		argument string,
		creator sdk.AccAddress) (int64, sdk.Error)
	SetParams(ctx sdk.Context, params Params)
	LikeArgument(ctx sdk.Context, argumentID int64, creator sdk.AccAddress, amount sdk.Coin) (*stake.LikeResult, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	argumentKeeper argument.Keeper
	stakeKeeper    stake.Keeper
	backingKeeper  backing.ReadKeeper
	bankKeeper     bank.Keeper
	trubankKeeper  trubank.Keeper
	storyKeeper    story.WriteKeeper
	paramStore     params.Subspace

	challengeList app.UserList // challenge <-> story mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	argumentKeeper argument.Keeper,
	stakeKeeper stake.Keeper,
	backingKeeper backing.ReadKeeper,
	trubankKeeper trubank.Keeper,
	bankKeeper bank.Keeper,
	storyKeeper story.WriteKeeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		argumentKeeper,
		stakeKeeper,
		backingKeeper,
		bankKeeper,
		trubankKeeper,
		storyKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
		app.NewUserList(storyKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new challenge on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context,
	storyID int64,
	amount sdk.Coin,
	argumentID int64,
	argument string,
	creator sdk.AccAddress) (challengeID int64, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	err = k.stakeKeeper.ValidateAmount(ctx, amount)
	if err != nil {
		return 0, err
	}

	err = k.stakeKeeper.ValidateStoryState(ctx, storyID)
	if err != nil {
		return 0, err
	}

	if amount.Denom != app.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid backing token.")
	}

	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for challenge.")
	}

	if amount.Amount.LT(k.minChallengeStake(ctx)) {
		return 0, sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	// make sure creator hasn't already challenged
	if k.challengeList.Includes(ctx, k, storyID, creator) {
		return 0, ErrDuplicateChallenge(storyID, creator)
	}

	stakeID := k.GetNextID(ctx)

	argumentID, err = k.argumentKeeper.Create(ctx, stakeID, storyID, argumentID, argument, creator)
	if err != nil {
		return 0, err
	}

	// create implicit false vote
	vote := stake.Vote{
		ID:         stakeID,
		StoryID:    storyID,
		Amount:     amount,
		ArgumentID: argumentID,
		Creator:    creator,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new challenge with embedded vote
	challenge := Challenge{&vote}

	// persist challenge
	k.setChallenge(ctx, challenge)

	// persist challenge <-> story mapping
	k.challengeList.Append(ctx, k, storyID, creator, challenge.ID())

	// deduct challenge amount from user
	_, err = k.trubankKeeper.SubtractCoin(ctx, creator, amount, storyID, trubank.Challenge, challenge.ID())
	if err != nil {
		return
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

// ChallengersByStoryID returns a list of addresses that challenged a specific story.
func (k Keeper) ChallengersByStoryID(ctx sdk.Context, storyID int64) (backers []sdk.AccAddress, err sdk.Error) {

	// iterate over and return challenges for a game
	err = k.challengeList.Map(ctx, k, storyID, func(challengeID int64) sdk.Error {
		challenge, err := k.Challenge(ctx, challengeID)
		if err != nil {
			return err
		}
		backers = append(backers, challenge.Creator())
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

// LikeArgument likes and argument
func (k Keeper) LikeArgument(ctx sdk.Context, argumentID int64, creator sdk.AccAddress, amount sdk.Coin) (*stake.LikeResult, sdk.Error) {
	err := k.argumentKeeper.RegisterLike(ctx, argumentID, creator)
	if err != nil {
		return nil, err
	}

	argument, err := k.argumentKeeper.Argument(ctx, argumentID)
	if err != nil {
		return nil, err
	}

	challenge, err := k.Challenge(ctx, argument.StakeID)
	if err != nil {
		return nil, err
	}

	story, err := k.storyKeeper.Story(ctx, challenge.StoryID())
	if err != nil {
		return nil, err
	}

	challengeID, err := k.Create(ctx, story.ID, amount, argumentID, "", creator)
	if err != nil {
		return nil, err
	}

	stakeToCredRatio := k.stakeKeeper.GetParams(ctx).StakeToCredRatio
	likeCredAmount := amount.Amount.Div(stakeToCredRatio)

	_, err = k.trubankKeeper.MintAndAddCoin(
		ctx,
		challenge.Creator(),
		story.CategoryID,
		story.ID,
		trubank.ChallengeLike,
		argument.StakeID,
		likeCredAmount)
	if err != nil {
		return nil, err
	}

	return &stake.LikeResult{
		StakeID:         challengeID,
		ArgumentID:      argumentID,
		ArgumentCreator: challenge.Creator(),
		CredEarned:      sdk.NewCoin(amount.Denom, likeCredAmount),
		StoryID:         story.ID,
	}, nil
}

// TotalChallengeAmount returns the total of all backings
func (k Keeper) TotalChallengeAmount(ctx sdk.Context, storyID int64) (
	totalCoin sdk.Coin, err sdk.Error) {

	logger := ctx.Logger().With("module", "challenge")
	totalAmount := sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())

	err = k.challengeList.Map(ctx, k, storyID, func(id int64) sdk.Error {
		challenge, err := k.Challenge(ctx, id)
		if err != nil {
			return err
		}
		totalAmount = totalAmount.Plus(challenge.Amount())

		return nil
	})

	if err != nil {
		return
	}

	logger.Info(fmt.Sprintf("Total Challenge Amount %s", totalAmount))

	return totalAmount, nil
}

func (k Keeper) setChallenge(ctx sdk.Context, challenge Challenge) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(challenge.ID()),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(challenge))
}
