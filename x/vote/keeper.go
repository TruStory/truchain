package vote

import (
	"net/url"

	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/x/bank"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	Get(ctx sdk.Context, id int64) (vote app.Vote, err sdk.Error)
	GetVotesByGame(ctx sdk.Context, gameID int64) ([]app.Vote, sdk.Error)
	Tally(ctx sdk.Context, gameID int64) (yes []app.Vote, no []app.Vote, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin,
		choice bool, comment string, creator sdk.AccAddress,
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

	storyKeeper story.WriteKeeper
	gameKeeper  game.WriteKeeper
	bankKeeper  bank.Keeper

	voterList app.UserList
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	gameKeeper game.WriteKeeper,
	bankKeeper bank.Keeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		storyKeeper,
		gameKeeper,
		bankKeeper,
		app.NewUserList(gameKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new vote on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin,
	choice bool, comment string, creator sdk.AccAddress,
	evidence []url.URL) (int64, sdk.Error) {

	// get the story
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// make sure validation game has started
	if story.GameID <= 0 {
		return 0, ErrGameNotStarted(storyID)
	}

	// check if this voter has already cast a vote
	if k.voterList.Includes(ctx, k, story.GameID, creator) {
		return 0, ErrDuplicateVoteForGame(story.GameID, creator)
	}

	// check if user has the funds
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds to vote on story.")
	}

	// deduct vote fee from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return 0, err
	}

	// create a new vote
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		Amount:    amount,
		Argument:  comment,
		Creator:   creator,
		Evidence:  evidence,
		Vote:      choice,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// persist vote
	k.set(ctx, vote)

	// persist game <-> vote association
	k.voterList.Append(ctx, k, story.GameID, creator, vote.ID)

	return vote.ID, nil
}

// Get returns a `Vote` from the KVStore
func (k Keeper) Get(ctx sdk.Context, id int64) (vote app.Vote, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	if bz == nil {
		return vote, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinary(bz, &vote)

	return vote, nil
}

// GetVotesByGame returns a list of votes for a given game
func (k Keeper) GetVotesByGame(
	ctx sdk.Context, gameID int64) (votes []app.Vote, err sdk.Error) {

	// iterate over voter list and get votes
	err = k.voterList.Map(ctx, k, gameID, func(voterID int64) sdk.Error {
		vote, err := k.Get(ctx, voterID)
		if err != nil {
			return err
		}
		votes = append(votes, vote)

		return nil
	})

	if err != nil {
		return votes, err
	}

	return votes, nil
}

// Tally votes
func (k Keeper) Tally(
	ctx sdk.Context, gameID int64) (yes []app.Vote, no []app.Vote, err sdk.Error) {

	err = k.voterList.Map(ctx, k, gameID, func(voteID int64) sdk.Error {
		vote, err := k.Get(ctx, voteID)
		if err != nil {
			return err
		}

		if vote.Vote == true {
			yes = append(yes, vote)
		} else {
			no = append(no, vote)
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

// ============================================================================

// saves a `Vote` in the KVStore
func (k Keeper) set(ctx sdk.Context, vote app.Vote) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(vote.ID),
		k.GetCodec().MustMarshalBinary(vote))
}
