package challenge

import (
	"fmt"
	"net/url"

	"github.com/cosmos/cosmos-sdk/x/bank"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	queue "github.com/cosmos/cosmos-sdk/store"
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
	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
		creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error)

	Update(
		ctx sdk.Context, challengeID int64, amount sdk.Coin, argument string,
		creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// ReadWriteKeeper defines a module interface that facilities read/write access to truchain data
type ReadWriteKeeper interface {
	ReadKeeper
	WriteKeeper
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	storyKeeper story.ReadWriteKeeper
	bankKeeper  bank.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey, sk story.ReadWriteKeeper, bankKeeper bank.Keeper,
	codec *amino.Codec) Keeper {

	return Keeper{app.NewKeeper(codec, storeKey), sk, bankKeeper}
}

// ============================================================================

// Create adds a new challenge on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
	creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error) {

	// get the story being challenged
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// validate challenger stake before creating challenge
	if err = validateStake(ctx, k, storyID, creator, amount); err != nil {
		return 0, err
	}

	// check if story already has a challenge
	if story.ChallengeID > 0 {
		return 0, ErrStoryAlreadyChallenged(storyID)
	}

	// create an initial empty challenge pool
	coinName, err := k.storyKeeper.GetCoinName(ctx, storyID)
	if err != nil {
		return 0, err
	}
	emptyPool := sdk.NewCoin(coinName, sdk.ZeroInt())

	// create new challenge type
	challenge := Challenge{
		k.GetNextID(ctx),
		storyID,
		creator,
		ctx.BlockHeader().Time.Add(DefaultParams().Expires),
		emptyPool,
		false,
		thresholdAmount(story),
		ctx.BlockHeight(),
		ctx.BlockHeader().Time,
		ctx.BlockHeight(),
		ctx.BlockHeader().Time,
	}

	story.ChallengeID = challenge.ID
	k.storyKeeper.UpdateStory(ctx, story)

	// push challenge id onto queue that will get checked
	// on each block tick for expired challenges
	q := queue.NewQueue(k.GetCodec(), k.GetStore(ctx))
	q.Push(challenge.ID)

	// add creator of challenge as challenger
	err = addChallenger(ctx, k, &challenge, amount, argument, creator, evidence)
	if err != nil {
		return 0, err
	}

	// set challenge in KVStore
	k.set(ctx, challenge)

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

// Update mutates an existing challenge, adding a new challenger and updating the pool
func (k Keeper) Update(
	ctx sdk.Context, challengeID int64, amount sdk.Coin, argument string,
	creator sdk.AccAddress, evidence []url.URL) (int64, sdk.Error) {

	challenge, err := k.Get(ctx, challengeID)
	if err != nil {
		return 0, err
	}

	// validate challenger stake before adding
	if err = validateStake(ctx, k, challenge.StoryID, creator, amount); err != nil {
		return 0, err
	}

	challenge, err = k.Get(ctx, challenge.ID)
	if err != nil {
		return 0, err
	}

	// check if user has already challenged
	challengerIDKey := k.getChallengerIDKey(challenge.ID, creator)
	bz := k.GetStore(ctx).Get(challengerIDKey)
	if bz != nil {
		return 0, ErrDuplicateChallenger(challenge.ID, creator)
	}

	// add this challenger to the challenger list
	err = addChallenger(ctx, k, &challenge, amount, argument, creator, evidence)
	if err != nil {
		return 0, err
	}

	// update existing challenge in KVStore
	k.set(ctx, challenge)

	return challenge.ID, nil
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

// saves the `Challenge` in the KVStore
func (k Keeper) set(ctx sdk.Context, challenge Challenge) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(challenge.ID),
		k.GetCodec().MustMarshalBinary(challenge))
}

// GetChallengerIDKey returns a key of form "challenges:id:5:userAddr:[userAddr]"
func (k Keeper) getChallengerIDKey(id int64, userAddr sdk.AccAddress) []byte {
	prefix := fmt.Sprintf("%s:id:%d:userAddr:%s", k.GetStoreKey().Name(), id, userAddr.String())
	return []byte(prefix)
}

// getChallengerPrefix returns a key of form "challenges:id:5:userAddr:"
func (k Keeper) getChallengerPrefix(id int64) []byte {
	prefix := fmt.Sprintf("%s:id:%d:userAddr:", k.GetStoreKey().Name(), id)
	return []byte(prefix)
}

// ============================================================================

// addChallenger adds a challenger to the challenge
func addChallenger(
	ctx sdk.Context, k Keeper, challenge *Challenge, amount sdk.Coin,
	argument string, creator sdk.AccAddress, evidence []url.URL) sdk.Error {

	// update block time for good record keeping
	challenge.UpdatedBlock = ctx.BlockHeight()
	challenge.UpdatedTime = ctx.BlockHeader().Time

	// create new challenger
	challenger := Challenger{
		amount,
		argument,
		creator,
		evidence,
		ctx.BlockHeight(),
		ctx.BlockHeader().Time,
	}

	// persist challenger
	k.GetStore(ctx).Set(
		k.getChallengerIDKey(challenge.ID, creator),
		k.GetCodec().MustMarshalBinary(challenger))

	// deduct challenge amount from user
	_, _, err := k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return err
	}

	// add amount to challenge pool
	challenge.Pool = challenge.Pool.Plus(amount)

	// if threshold is reached, start challenge and allow voting to begin
	if challenge.Pool.Amount.GT(challenge.ThresholdAmount) {
		err = k.storyKeeper.StartChallenge(ctx, challenge.StoryID)
		if err != nil {
			return err
		}
		challenge.Started = true
	}

	return nil
}

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
	minStake := sdk.NewCoin(coinName, DefaultParams().MinChallengeStake)
	if amount.IsLT(minStake) {
		return sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	return
}
