package challenge

import (
	"net/url"

	"github.com/cosmos/cosmos-sdk/x/bank"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	Challenge(ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error)
	ChallengesByGame(ctx sdk.Context, gameID int64) (challenges []Challenge, err sdk.Error)
	Tally(ctx sdk.Context, gameID int64) (falseVotes []Challenge, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
		creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	gameQueueKey sdk.StoreKey

	bankKeeper  bank.Keeper
	gameKeeper  game.WriteKeeper
	storyKeeper story.WriteKeeper

	challengeList app.UserList // challenge <-> game mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey, gameQueueKey sdk.StoreKey, bankKeeper bank.Keeper,
	gameKeeper game.WriteKeeper, storyKeeper story.WriteKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		gameQueueKey,
		bankKeeper,
		gameKeeper,
		storyKeeper,
		app.NewUserList(gameKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new challenge on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
	creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error) {

	// validate stake before creating challenge
	err := validateStake(ctx, k, storyID, creator, amount)
	if err != nil {
		return 0, err
	}

	// get the story
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// create game if one doesn't exist yet
	gameID := story.GameID
	if gameID == 0 {
		gameID, err = k.gameKeeper.Create(ctx, story.ID, creator)
		if err != nil {
			return 0, err
		}
	}

	// make sure creator hasn't already challenged
	if k.challengeList.Includes(ctx, k, gameID, creator) {
		return 0, ErrDuplicateChallenge(gameID, creator)
	}

	// create implicit false vote
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		Amount:    amount,
		Argument:  argument,
		Creator:   creator,
		Evidence:  evidence,
		Vote:      false,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new challenge with embedded vote
	challenge := Challenge{vote}

	// persist challenge
	k.GetStore(ctx).Set(
		k.GetIDKey(challenge.ID()),
		k.GetCodec().MustMarshalBinary(challenge))

	// persist challenge <-> game mapping
	k.challengeList.Append(ctx, k, gameID, creator, challenge.ID())

	// deduct challenge amount from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return 0, err
	}

	// update game pool
	_, err = k.gameKeeper.Update(ctx, gameID, amount)
	if err != nil {
		return 0, err
	}

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
	k.GetCodec().MustUnmarshalBinary(bz, &challenge)

	return
}

// ChallengesByGame returns the list of challenges for a game id
func (k Keeper) ChallengesByGame(
	ctx sdk.Context, gameID int64) (challenges []Challenge, err sdk.Error) {

	// iterate over and return challenges for a game
	err = k.challengeList.Map(ctx, k, gameID, func(challengeID int64) sdk.Error {
		challenge, err := k.Challenge(ctx, challengeID)
		if err != nil {
			return err
		}
		challenges = append(challenges, challenge)

		return nil
	})

	if err != nil {
		return challenges, err
	}

	return challenges, nil
}

// Tally challenges for voting
func (k Keeper) Tally(
	ctx sdk.Context, gameID int64) (falseVotes []Challenge, err sdk.Error) {

	err = k.challengeList.Map(ctx, k, gameID, func(challengeID int64) sdk.Error {
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

// validate if a challenger has the right staking amount
func validateStake(
	ctx sdk.Context, k Keeper, storyID int64,
	creator sdk.AccAddress, amount sdk.Coin) (err sdk.Error) {

	// get category coin name
	coinName, err := k.storyKeeper.GetCoinName(ctx, storyID)
	if err != nil {
		return
	}

	// check if user has the stake they are claiming
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return sdk.ErrInsufficientFunds("Insufficient funds for challenging story.")
	}

	// check if challenge amount is greater than minimum stake
	minStake := sdk.NewCoin(coinName, game.DefaultParams().MinChallengeStake)
	if amount.IsLT(minStake) {
		return sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	return
}
