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
	ChallengeByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (challenge Challenge, err sdk.Error)
	GetParams(ctx sdk.Context) Params
	Tally(ctx sdk.Context, storyID int64) (falseVotes []Challenge, err sdk.Error)
	TotalChallengeAmount(ctx sdk.Context, storyID int64) (
		totalCoin sdk.Coin, err sdk.Error)
	ChallengeThreshold(ctx sdk.Context, storyID int64) (sdk.Coin, sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
		creator sdk.AccAddress, toggled bool) (int64, sdk.Error)
	Update(ctx sdk.Context, challenge Challenge)
	Delete(ctx sdk.Context, challenge Challenge) sdk.Error
	SetParams(ctx sdk.Context, params Params)
	LikeArgument(ctx sdk.Context, argumentID int64, creator sdk.AccAddress, amount sdk.Coin) (int64, sdk.Error)
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
	ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
	creator sdk.AccAddress,
	toggled bool) (challengeID int64, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	err = k.stakeKeeper.ValidateAmount(ctx, amount)
	if err != nil {
		return 0, err
	}

	err = k.stakeKeeper.ValidateStoryState(ctx, storyID, toggled)
	if err != nil {
		return 0, err
	}

	if err = k.stakeKeeper.ValidateArgument(ctx, argument); err != nil {
		return 0, err
	}

	if amount.Denom != app.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid backing token.")
	}

	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for challenge.")
	}

	if !toggled && amount.Amount.LT(k.minChallengeStake(ctx)) {
		return 0, sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	// make sure creator hasn't already challenged
	if k.challengeList.Includes(ctx, k, storyID, creator) {
		return 0, ErrDuplicateChallenge(storyID, creator)
	}

	stakeID := k.GetNextID(ctx)

	argumentID := int64(0)
	if len(argument) > 0 {
		argumentID, err = k.argumentKeeper.Create(ctx, stakeID, argument)
		if err != nil {
			return 0, sdk.ErrInternal("Error creating argument")
		}
	}

	// create implicit false vote
	vote := stake.Vote{
		ID:         k.GetNextID(ctx),
		StoryID:    storyID,
		Amount:     amount,
		ArgumentID: argumentID,
		Weight:     sdk.NewInt(0),
		Creator:    creator,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new challenge with embedded vote
	challenge := Challenge{&vote}

	// persist challenge
	k.GetStore(ctx).Set(
		k.GetIDKey(challenge.ID()),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(challenge))

	// persist challenge <-> story mapping
	k.challengeList.Append(ctx, k, storyID, creator, challenge.ID())

	// deduct challenge amount from user
	_, err = k.trubankKeeper.SubtractCoin(ctx, creator, amount, storyID, trubank.Challenge, challenge.ID())
	if err != nil {
		return
	}

	// err = k.checkThreshold(ctx, storyID)
	// if err != nil {
	// 	return 0, err
	// }

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

// LikeArgument likes and argument
func (k Keeper) LikeArgument(ctx sdk.Context, argumentID int64, creator sdk.AccAddress, amount sdk.Coin) (int64, sdk.Error) {
	k.argumentKeeper.RegisterLike(ctx, argumentID, creator)

	argument, err := k.argumentKeeper.Argument(ctx, argumentID)
	if err != nil {
		return 0, sdk.ErrInternal("error getting argument")
	}

	challenge, err := k.Challenge(ctx, argument.StakeID)
	if err != nil {
		return 0, sdk.ErrInternal("error getting challenge")
	}

	story, err := k.storyKeeper.Story(ctx, challenge.StoryID())
	if err != nil {
		return 0, sdk.ErrInternal("can't get story")
	}

	challengeID, err := k.Create(ctx, story.ID, amount, "", creator, false)
	if err != nil {
		return 0, sdk.ErrInternal("cannot create challenge")
	}

	// amount of cred for a like
	cred := sdk.NewInt(1 * app.Shanev)

	_, err = k.trubankKeeper.MintAndAddCoin(ctx, challenge.Creator(), story.CategoryID, story.ID, trubank.Like, cred)
	if err != nil {
		return 0, sdk.ErrInternal("cant mint coins")
	}

	return challengeID, nil
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

func (k Keeper) checkThreshold(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", "challenge")

	// only check threshold if it is in pending state
	s, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}
	if s.Status != story.Pending {
		return nil
	}

	backingPool, err := k.backingKeeper.TotalBackingAmount(ctx, storyID)
	if err != nil {
		return err
	}

	challengePool, err := k.TotalChallengeAmount(ctx, storyID)
	if err != nil {
		return err
	}

	challengeThreshold, err := k.ChallengeThreshold(ctx, storyID)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf(
		"Backing pool: %s, challenge pool: %s, threshold: %s",
		backingPool, challengePool, challengeThreshold))

	if challengePool.IsGTE(challengeThreshold) {
		err := k.storyKeeper.StartVotingPeriod(ctx, storyID)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf(
			"Challenge threshold met. Voting started for story %d", storyID))
	}

	return nil
}

// ChallengeThreshold returns the challenge threshold for a story
func (k Keeper) ChallengeThreshold(ctx sdk.Context, storyID int64) (sdk.Coin, sdk.Error) {
	backingPool, err := k.backingKeeper.TotalBackingAmount(ctx, storyID)
	if err != nil {
		return sdk.NewCoin(app.StakeDenom, sdk.ZeroInt()), err
	}

	// calculate challenge threshold amount (based on total backings)
	totalBackingDec := sdk.NewDecFromInt(backingPool.Amount)
	challengeThresholdAmount := totalBackingDec.
		Mul(k.challengeToBackingRatio(ctx)).
		TruncateInt()

	// challenge threshold can't be less than min challenge stake
	minChallengeThreshold := k.minChallengeThreshold(ctx)
	if challengeThresholdAmount.LT(minChallengeThreshold) {
		return sdk.NewCoin(backingPool.Denom, minChallengeThreshold), nil
	}

	return sdk.NewCoin(backingPool.Denom, challengeThresholdAmount), nil
}

// Update updates the challenge vote
func (k Keeper) Update(ctx sdk.Context, challenge Challenge) {

	newChallenge := Challenge{
		Vote: challenge.Vote,
	}

	k.setChallenge(ctx, newChallenge)
}

// Delete deletes a challenge and restores coins to the user.
func (k Keeper) Delete(ctx sdk.Context, challenge Challenge) sdk.Error {
	challengeIDKey := k.GetIDKey(challenge.ID())

	if !k.GetStore(ctx).Has(challengeIDKey) {
		return ErrNotFound(challenge.ID())
	}

	// removes challenge
	k.GetStore(ctx).Delete(
		k.GetIDKey(challenge.ID()))

	// restore coins
	_, err := k.trubankKeeper.AddCoin(ctx, challenge.Creator(), challenge.Amount(), challenge.StoryID(), trubank.ChallengeReturned, challenge.ID())
	if err != nil {
		return err
	}

	// removes challenge association from the backing list
	k.challengeList.Delete(ctx, k, challenge.StoryID(), challenge.Creator())

	return nil
}

func (k Keeper) setChallenge(ctx sdk.Context, challenge Challenge) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(challenge.ID()),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(challenge))
}
