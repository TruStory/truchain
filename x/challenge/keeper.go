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

	Get(ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error)
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
	challengeByGameKey := k.challengeByGameIDKey(ctx, gameID, creator)
	bz := k.GetStore(ctx).Get(challengeByGameKey)
	if bz != nil {
		return 0, ErrDuplicateChallenge(gameID, creator)
	}

	// create new challenge
	challenge := Challenge{
		ID:        k.GetNextID(ctx),
		Amount:    amount,
		Argument:  argument,
		Creator:   creator,
		Evidence:  evidence,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// persist challenge
	k.GetStore(ctx).Set(
		k.GetIDKey(challenge.ID),
		k.GetCodec().MustMarshalBinary(challenge))

	// persist game <-> challenge association
	k.GetStore(ctx).Set(
		challengeByGameKey,
		k.GetCodec().MustMarshalBinary(challenge.ID))

	// update game pool
	_, err = k.gameKeeper.Update(ctx, gameID, amount)
	if err != nil {
		return 0, err
	}

	// deduct challenge amount from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return 0, err
	}

	return challenge.ID, nil
}

// Get the challenge for the given id
func (k Keeper) Get(ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(challengeID))
	if bz == nil {
		return challenge, ErrNotFound(challengeID)
	}
	k.GetCodec().MustUnmarshalBinary(bz, &challenge)

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
