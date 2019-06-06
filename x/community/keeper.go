package community

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	recordkeeper "github.com/shanev/cosmos-record-keeper/recordkeeper"
	log "github.com/tendermint/tendermint/libs/log"
)

const (
	// StoreKey represents the KVStore for the communities
	StoreKey = ModuleName
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	recordkeeper.RecordKeeper
}

// NewKeeper creates a new keeper of the community Keeper
func NewKeeper(storeKey sdk.StoreKey, codec *codec.Codec) Keeper {
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
	logger.Info(fmt.Sprintf("Created new community: %s", community.String()))

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
func (k Keeper) Communities(ctx sdk.Context) (communities []Community) {
	err := k.Each(ctx, func(bytes []byte) bool {
		var community Community
		k.Codec.MustUnmarshalBinaryLengthPrefixed(bytes, &community)
		communities = append(communities, community)
		return true
	})

	if err != nil {
		return
	}

	return
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", StoreKey)
}
