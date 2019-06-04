package community

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	recordkeeper "github.com/shanev/cosmos-record-keeper/recordkeeper"
	amino "github.com/tendermint/go-amino"
	log "github.com/tendermint/tendermint/libs/log"
)

const (
	// StoreKey represents the KVStore for the communities
	StoreKey = "community"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	recordkeeper.RecordKeeper
}

// NewKeeper creates a new keeper of the community Keeper
func NewKeeper(storeKey sdk.StoreKey, codec *amino.Codec) Keeper {
	return Keeper{recordkeeper.NewRecordKeeper(storeKey, codec)}
}

// NewCommunity creates a new community
func (k Keeper) NewCommunity(ctx sdk.Context, name string, slug string, description string) Community {
	logger := getLogger(ctx)

	community := Community{
		ID:          k.IncrementID(ctx),
		Name:        name,
		Slug:        slug,
		Description: description,
		Timestamp:   app.NewTimestamp(ctx.BlockHeader()),
	}

	k.Set(ctx, community.ID, community)
	logger.Info("Created new community: " + community.String())

	return community
}

// Community returns a community by its ID
func (k Keeper) Community(ctx sdk.Context, id uint64) (community Community, err sdk.Error) {
	err = k.Get(ctx, id, &community)
	if err != nil {
		return community, ErrCommunityNotFound(id)
	}

	return community, nil
}

// Communities gets all communities from the KVStore
func (k Keeper) Communities(ctx sdk.Context) (communities []Community, err sdk.Error) {
	err = k.Each(ctx, func(bytes []byte) bool {
		var community Community
		k.Codec.MustUnmarshalBinaryLengthPrefixed(bytes, &community)
		communities = append(communities, community)
		return true
	})
	return
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", StoreKey)
}
