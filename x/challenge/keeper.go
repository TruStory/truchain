package challenge

import (
	"net/url"

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
		evidence []url.URL) (int64, sdk.Error)
}

// ReadWriteKeeper defines a module interface that facilities read/write access to truchain data
type ReadWriteKeeper interface {
	ReadKeeper
	WriteKeeper
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	challengeKey sdk.StoreKey
	storyKeeper  story.ReadWriteKeeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(challengeKey sdk.StoreKey, sk story.ReadWriteKeeper, codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec), challengeKey, sk}
}

// // ============================================================================

// NewChallenge adds a new challenge on a story in the KVStore
func (k Keeper) NewChallenge(
	ctx sdk.Context,
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence []url.URL) (int64, sdk.Error) {

	// make sure we have the story being challenged
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	challenge := NewChallenge(
		k.GetNextID(ctx, k.challengeKey),
		storyID,
		amount,
		argument,
		creator,
		[]url.URL{},
		ctx.BlockHeight(),
		ctx.BlockHeader().Time)

	// persist challenge value
	k.setChallenge(ctx, challenge)

	// add story id to challenged stories list
	k.storyKeeper.Challenge(ctx, story)

	return challenge.ID, nil
}

// GetChallenge the challenge for the given id
func (k Keeper) GetChallenge(ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error) {
	store := ctx.KVStore(k.challengeKey)
	key := getChallengeIDKey(k, challengeID)
	bz := store.Get(key)
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
	store := ctx.KVStore(k.challengeKey)
	store.Set(getChallengeIDKey(k, challenge.ID), k.marshal(challenge))
}

// ============================================================================

// returns a key of the form "challenges:id:[ID]"
func getChallengeIDKey(k Keeper, id int64) []byte {
	return app.GetIDKey(k.challengeKey, id)
}
