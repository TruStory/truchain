package argument

import (
	"fmt"

	"github.com/TruStory/truchain/x/story"
	amino "github.com/tendermint/go-amino"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// StoreKey is string representation of the store key
	StoreKey = "argument"
)

// Keeper stores keys and other keepers needed to read/write arguments
type Keeper struct {
	app.Keeper

	storyKeeper story.WriteKeeper
	paramStore  params.Subspace
}

// NewKeeper constructs a new argument keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		storyKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

// Argument returns the argument for the given id
func (k Keeper) Argument(ctx sdk.Context, id int64) (argument Argument, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	if bz == nil {
		return argument, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &argument)

	return
}

// Create stores a new argument in the argument store
func (k Keeper) Create(
	ctx sdk.Context,
	stakeID int64,
	storyID int64,
	argumentID int64,
	body string,
	creator sdk.AccAddress) (int64, sdk.Error) {

	switch {
	case len(body) == 0 && argumentID == 0:
		return 0, ErrInvalidArgument("Must have either body or argumentID")
	case len(body) > 0 && argumentID > 0:
		return 0, ErrInvalidArgument("Cannot have both body and argumentID")
	case len(body) == 0 && argumentID > 0:
		return argumentID, nil
	}

	err := k.validateArgumentBody(ctx, body)
	if err != nil {
		return 0, err
	}

	arg := Argument{
		ID:        k.GetNextID(ctx),
		StoryID:   storyID,
		StakeID:   stakeID,
		Body:      body,
		Creator:   creator,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setArgument(ctx, arg)

	return arg.ID, nil
}

// RegisterLike registers a like for the argument
func (k Keeper) RegisterLike(ctx sdk.Context, argumentID int64, creator sdk.AccAddress) sdk.Error {
	like := Like{
		ArgumentID: argumentID,
		Creator:    creator,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}

	// append argument <-> like association
	// argument:id:[ID]:likes:creator:[0xdeadbeef] = Like{}
	store := k.GetStore(ctx)
	store.Set(
		k.likesKey(argumentID, creator),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(like))

	return nil
}

// LikesByArgumentID returns likes for a given argument id
func (k Keeper) LikesByArgumentID(ctx sdk.Context, argumentID int64) (likes []Like, err sdk.Error) {
	// iterate through prefix argument:id:[ID]:likes:creator:
	searchPrefix := fmt.Sprintf(
		"%s:id:%d:likes:creator:",
		k.GetStoreKey().Name(),
		argumentID,
	)

	err = k.EachPrefix(ctx, searchPrefix, func(val []byte) bool {
		var like Like
		k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &like)
		likes = append(likes, like)
		return true
	})

	return
}

func (k Keeper) setArgument(ctx sdk.Context, argument Argument) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(argument.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(argument))
}

func (k Keeper) likesKey(argumentID int64, creator sdk.AccAddress) []byte {
	key := fmt.Sprintf(
		"%s:id:%d:likes:creator:%s",
		k.GetStoreKey().Name(),
		argumentID,
		creator.String(),
	)

	return []byte(key)
}

func (k Keeper) validateArgumentBody(ctx sdk.Context, argument string) sdk.Error {
	len := len([]rune(argument))
	minArgumentLength := k.GetParams(ctx).MinArgumentLength
	maxArgumentLength := k.GetParams(ctx).MaxArgumentLength

	if len > 0 && (len < minArgumentLength) {
		return ErrArgumentTooShortMsg(argument, minArgumentLength)
	}

	if len > 0 && (len > maxArgumentLength) {
		return ErrArgumentTooLongMsg(maxArgumentLength)
	}

	return nil
}
