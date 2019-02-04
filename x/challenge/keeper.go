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

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	Challenge(
		ctx sdk.Context, challengeID int64) (challenge Challenge, err sdk.Error)

	ChallengesByGameID(
		ctx sdk.Context, gameID int64) (challenges []Challenge, err sdk.Error)

	ChallengeByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (challenge Challenge, err sdk.Error)

	Tally(ctx sdk.Context, gameID int64) (falseVotes []Challenge, err sdk.Error)
	Challenges(ctx sdk.Context) (challenges []Challenge)
	ExportState(ctx sdk.Context, dnh string, bh int64)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin, argument string,
		creator sdk.AccAddress) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	// list of games waiting to meet challenge threshold
	pendingGameListKey sdk.StoreKey

	backingKeeper backing.ReadKeeper
	bankKeeper    bank.Keeper
	gameKeeper    game.WriteKeeper
	storyKeeper   story.WriteKeeper

	challengeList app.UserList // challenge <-> game mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	pendingGameListKey sdk.StoreKey,
	backingKeeper backing.ReadKeeper,
	bankKeeper bank.Keeper,
	gameKeeper game.WriteKeeper,
	storyKeeper story.WriteKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		pendingGameListKey,
		backingKeeper,
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
	creator sdk.AccAddress) (challengeID int64, err sdk.Error) {

	logger := ctx.Logger().With("module", "challenge")

	// check is user has the coins they are staking
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for challenging story.")
	}

	// get category coin name
	coinName, err := k.storyKeeper.CategoryDenom(ctx, storyID)
	if err != nil {
		return
	}

	catCoin := app.NewCategoryCoin(coinName, amount)

	// check if challenge amount is greater than minimum stake
	if catCoin.Amount.LT(game.DefaultParams().MinChallengeStake) {
		return 0, sdk.ErrInsufficientFunds("Does not meet minimum stake amount.")
	}

	// get the story
	story, err := k.storyKeeper.Story(ctx, storyID)
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
		Amount:    catCoin,
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

	// persist challenge <-> game mapping
	k.challengeList.Append(ctx, k, gameID, creator, challenge.ID())

	// convert from trustake if needed
	if amount.Denom == params.StakeDenom {
		err = app.SwapCoin(ctx, k.bankKeeper, amount, catCoin, creator)
	}

	// deduct challenge amount from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{catCoin})
	if err != nil {
		return 0, err
	}

	// add another amount to the challenge pool
	err = k.gameKeeper.AddToChallengePool(ctx, gameID, catCoin)
	if err != nil {
		return 0, err
	}

	msg := fmt.Sprintf("Challenged story %d with %s by %s",
		storyID, catCoin.String(), creator.String())
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

// ChallengesByGameID returns the list of challenges for a game id
func (k Keeper) ChallengesByGameID(
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
	challengeID := k.challengeList.Get(ctx, k, s.GameID, creator)
	challenge, err = k.Challenge(ctx, challengeID)

	return
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

// Challenges returns all challenges in the order they appear in the store
func (k Keeper) Challenges(ctx sdk.Context) (challenges []Challenge) {

	// get store
	store := k.GetStore(ctx)

	// builds prefix "challenges:id"
	searchKey := fmt.Sprintf("%s:id", k.GetStoreKey().Name())
	searchPrefix := []byte(searchKey)

	// setup Iterator
	iter := sdk.KVStorePrefixIterator(store, searchPrefix)
	defer iter.Close()

	// iterates through keyspace to find all challenges
	for ; iter.Valid(); iter.Next() {
		var challenge Challenge

		k.GetCodec().MustUnmarshalBinaryLengthPrefixed(
			iter.Value(), &challenge)

		challenges = append(challenges, challenge)

	}

	return challenges
}

// ExportState gets all the current challenges and calls app.WriteJSONtoNodeHome() to write data to file.
func (k Keeper) ExportState(ctx sdk.Context, dnh string, bh int64) {
	challenges := k.Challenges(ctx)
	app.WriteJSONtoNodeHome(challenges, dnh, bh, fmt.Sprintf("%s.json", k.GetStoreKey().Name()))

}

// ============================================================================

func (k Keeper) pendingGameList(ctx sdk.Context) list.List {
	return list.NewList(
		k.GetCodec(),
		ctx.KVStore(k.pendingGameListKey))
}
