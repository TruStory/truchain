package challenge

import (
	"net/url"

	"github.com/cosmos/cosmos-sdk/x/bank"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	NewChallenge(
		ctx sdk.Context,
		storyID int64,
		amount sdk.Coin,
		argument string,
		creator sdk.AccAddress,
		evidence []url.URL,
		reason Reason) (int64, sdk.Error)
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
func NewKeeper(storeKey sdk.StoreKey, sk story.ReadWriteKeeper, bankKeeper bank.Keeper, codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec, storeKey), sk, bankKeeper}
}

// // ============================================================================

// NewChallenge adds a new challenge on a story in the KVStore
func (k Keeper) NewChallenge(
	ctx sdk.Context,
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence []url.URL,
	reason Reason) (id int64, err sdk.Error) {

	// make sure we have the story being challenged
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return
	}

	// get category coin name
	coinName, err := k.storyKeeper.GetCoinName(ctx, storyID)
	if err != nil {
		return
	}

	// load default challenge parameters
	params := NewParams()
	minStake := sdk.NewCoin(coinName, params.MinChallengeStake)

	// check if user has the stake they are claiming
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{minStake}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for challenging story.")
	}

	// check if challenge amount meets minimum stake
	if amount.IsLT(minStake) {
		return 0, sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	// check if story already has a challenge
	var challenge Challenge
	if story.ChallengeID > 0 {
		challenge, err = k.GetChallenge(ctx, story.ChallengeID)
		if err != nil {
			return 0, err
		}
		// add user to challenge
		challenge.Challengers = append(challenge.Challengers, creator)
		// add amount to challenge pool
		challenge.Pool.Plus(amount)
		// if threshold is reached, start challenge
		if challenge.Pool.Amount.GT(challenge.ThresholdAmount) {
			k.storyKeeper.StartChallenge(ctx, storyID)
		}

		// update block time for good record keeping
		challenge.UpdatedBlock = ctx.BlockHeight()
		challenge.UpdatedTime = ctx.BlockHeader().Time
	} else {
		// create new challenge type
		challenge = NewChallenge(
			k.GetNextID(ctx),
			storyID,
			amount,
			argument,
			[]sdk.AccAddress{creator},
			creator,
			evidence,
			amount,
			reason,
			getThresholdAmount(story),
			ctx.BlockHeight(),
			ctx.BlockHeader().Time)

		story.ChallengeID = challenge.ID
		k.storyKeeper.UpdateStory(ctx, story)
	}
	// persist new or updated challenge value
	k.setChallenge(ctx, challenge)

	// deduct challenge amount from  user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return
	}

	return challenge.ID, nil
}

// GetChallenge the challenge for the given id
func (k Keeper) GetChallenge(ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(challengeID))
	if bz == nil {
		return challenge, ErrNotFound(challengeID)
	}
	challenge = k.unmarshal(bz)

	return
}

// ============================================================================

// marshals the challenge and returns bytes
func (k Keeper) marshal(value Challenge) []byte {
	return k.GetCodec().MustMarshalBinary(value)
}

// unmarshals a challenge from the KVStore
func (k Keeper) unmarshal(bz []byte) (value Challenge) {
	k.GetCodec().MustUnmarshalBinary(bz, &value)

	return
}

// saves the `Challenge` in the KVStore
func (k Keeper) setChallenge(ctx sdk.Context, challenge Challenge) {
	store := k.GetStore(ctx)
	store.Set(k.GetIDKey(challenge.ID), k.marshal(challenge))
}

// ============================================================================

// TODO:
func getThresholdAmount(s story.Story) sdk.Int {
	return sdk.NewInt(10)
}
