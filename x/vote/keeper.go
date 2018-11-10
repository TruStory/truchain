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
	GetByGame(ctx sdk.Context, gameID int64) ([]app.Vote, sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	Create(
		ctx sdk.Context, storyID int64, amount sdk.Coin,
		choice bool, comment string, creator sdk.AccAddress,
		evidence []url.URL) (int64, sdk.Error)

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

	storyKeeper story.WriteKeeper
	gameKeeper  game.WriteKeeper
	bankKeeper  bank.Keeper

	voterList app.VoterList
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
		app.NewVoterList(gameKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new challenge on a story in the KVStore
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, amount sdk.Coin,
	choice bool, comment string, creator sdk.AccAddress,
	evidence []url.URL) (int64, sdk.Error) {

	// get the story
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// TODO: check if this voter has already cast a vote

	// create a new vote
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		Amount:    amount,
		Comment:   comment,
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
		// TODO: add error
		return vote, sdk.ErrInternal("ErrNotFound(id)")
	}
	k.GetCodec().MustUnmarshalBinary(bz, &vote)

	return
}

// GetByGame returns a list of votes for a given game
func (k Keeper) GetByGame(
	ctx sdk.Context, gameID int64) (votes []app.Vote, err sdk.Error) {

	// iterate over voter list and get votes
	k.voterList.Map(ctx, k, gameID, func(voterID int64) {
		vote, err := k.Get(ctx, voterID)
		if err != nil {
			return
		}
		votes = append(votes, vote)
	})

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
