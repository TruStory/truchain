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

	// NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	bankKeeper  bank.Keeper
	gameKeeper  game.WriteKeeper
	storyKeeper story.WriteKeeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey, bankKeeper bank.Keeper, gameKeeper game.WriteKeeper,
	storyKeeper story.WriteKeeper, codec *amino.Codec) Keeper {

	return Keeper{app.NewKeeper(codec, storeKey), bankKeeper, gameKeeper, storyKeeper}
}

// ============================================================================

// Create adds a new challenge on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, gameID int64, amount sdk.Coin, argument string,
	creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error) {

	// get the validation game
	game, err := k.gameKeeper.Get(ctx, gameID)
	if err != nil {
		return 0, err
	}

	// validate challenger stake before creating challenge
	if err = validateStake(ctx, k, game.StoryID, creator, amount); err != nil {
		return 0, err
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

	// deduct challenge amount from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return 0, err
	}

	k.gameKeeper.Update(ctx, game.ID, amount)

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

// Delete removes a challenge from the KVStore
func (k Keeper) delete(ctx sdk.Context, id int64) sdk.Error {
	store := k.GetStore(ctx)
	key := k.GetIDKey(id)
	bz := store.Get(key)
	if bz == nil {
		return ErrNotFound(id)
	}
	store.Delete(key)

	return nil
}

// ============================================================================

// [Shane] TODO: https://github.com/TruStory/truchain/issues/50
func thresholdAmount(s story.Story) sdk.Int {
	return sdk.NewInt(10)
}

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
