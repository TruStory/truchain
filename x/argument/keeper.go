package argument

import (
	"fmt"

	"github.com/TruStory/truchain/x/story"
	amino "github.com/tendermint/go-amino"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// StoreKey is string representation of the store key
	StoreKey = "argument"

	// LikeStoreKey is a string representation of the like store
	// LikeStoreKey = "likes"
)

// Keeper stores keys and other keepers needed to read/write arguments
type Keeper struct {
	app.Keeper

	storyKeeper story.WriteKeeper
}

// NewKeeper constructs a new argument keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		storyKeeper,
	}
}

// Argument returns the argument for the given id
func (k Keeper) Argument(ctx sdk.Context, id int64) (Argument, sdk.Error) {
	return Argument{}, nil
}

// Create stores a new argument in the argument store
func (k Keeper) Create(
	ctx sdk.Context,
	stakeID int64,
	body string) (int64, sdk.Error) {

	// TODO: validate body

	arg := Argument{
		ID:        k.GetNextID(ctx),
		StoryID:   0,
		StakeID:   stakeID,
		StakeType: nil,
		Body:      body,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setArgument(ctx, arg)

	return arg.ID, nil
}

// RegisterLike registers a like for the argument
func (k Keeper) RegisterLike(ctx sdk.Context, argumentID int64, creator sdk.AccAddress) sdk.Error {

	like := Like{
		ArgumentID: argumentID,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}

	// append argument <-> like association
	// argument:id:[ID]:creator:[0xdeadbeef] = Like{}
	store := k.GetStore(ctx)
	store.Set(
		k.argumentIDByCreatorKey(argumentID, creator),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(like))

	return nil
}

func (k Keeper) setArgument(ctx sdk.Context, argument Argument) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(argument.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(argument))
}

func (k Keeper) argumentIDByCreatorKey(argumentID int64, creator sdk.AccAddress) []byte {
	key := fmt.Sprintf(
		"%s:id:%d:creator:%s",
		k.GetStoreKey().Name(),
		argumentID,
		creator.String(),
	)

	return []byte(key)
}
