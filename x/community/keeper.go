package community

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	log "github.com/tendermint/tendermint/libs/log"
)

const (
	// StoreKey represents the KVStore for the communities
	StoreKey = "community"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	app.Keeper
}

// NewKeeper creates a new keeper of the community Keeper
func NewKeeper(storeKey sdk.StoreKey, codec *amino.Codec) Keeper {
	return Keeper{app.NewKeeper(codec, storeKey)}
}

// Create creates a new community
func (k Keeper) Create(ctx sdk.Context, name string, slug string, description string) Community {
	logger := getLogger(ctx)

	community := Community{
		ID:          k.GetNextID(ctx),
		Name:        name,
		Slug:        slug,
		Description: description,
		Timestamp:   app.NewTimestamp(ctx.BlockHeader()),
	}

	store := k.GetStore(ctx)
	store.Set(k.GetIDKey(community.ID), k.GetCodec().MustMarshalBinaryLengthPrefixed(community))

	logger.Info("Created new community: " + community.String())

	return community
}

// Get returns a community by its ID
func (k Keeper) Get(ctx sdk.Context, id int64) (community Community, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	if bz == nil {
		return community, ErrCommunityNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &community)

	return community, nil
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", StoreKey)
}
