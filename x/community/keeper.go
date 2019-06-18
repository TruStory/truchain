package community

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	recordkeeper "github.com/shanev/cosmos-record-keeper/recordkeeper"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	recordkeeper.RecordKeeper

	codec      *codec.Codec
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the community Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		recordkeeper.NewRecordKeeper(storeKey, codec),
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}

// NewCommunity creates a new community
func (k Keeper) NewCommunity(ctx sdk.Context, name string, slug string, description string) (community Community, err sdk.Error) {
	logger := getLogger(ctx)

	err = k.validateParams(ctx, name, slug, description)
	if err != nil {
		return
	}

	community = Community{
		ID:          k.IncrementID(ctx),
		Name:        name,
		Slug:        slug,
		Description: description,
		Timestamp:   app.NewTimestamp(ctx.BlockHeader()),
	}

	k.Set(ctx, community.ID, community)
	logger.Info(fmt.Sprintf("Created new community: %s", community.String()))

	return
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
		k.codec.MustUnmarshalBinaryLengthPrefixed(bytes, &community)
		communities = append(communities, community)
		return true
	})

	if err != nil {
		return
	}

	return
}

func (k Keeper) validateParams(ctx sdk.Context, name, slug, description string) (err sdk.Error) {
	params := k.GetParams(ctx)
	if len(name) < params.MinNameLength || len(name) > params.MaxNameLength {
		err = ErrInvalidCommunityMsg(
			fmt.Sprintf("Name must be between %d-%d chars in length", params.MinNameLength, params.MaxNameLength),
		)
	}
	if len(slug) < params.MinSlugLength || len(slug) > params.MaxSlugLength {
		err = ErrInvalidCommunityMsg(
			fmt.Sprintf("Slug must be between %d-%d chars in length", params.MinSlugLength, params.MaxSlugLength),
		)
	}
	if len(description) > params.MaxDescriptionLength {
		err = ErrInvalidCommunityMsg(
			fmt.Sprintf("Description must be less than %d chars in length", params.MaxDescriptionLength),
		)
	}

	return
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
