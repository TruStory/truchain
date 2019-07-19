package community

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the community Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}

// NewCommunity creates a new community
func (k Keeper) NewCommunity(ctx sdk.Context, id string, name string, description string, creator sdk.AccAddress) (community Community, err sdk.Error) {
	err = k.validateParams(ctx, id, name, description, creator)
	if err != nil {
		return
	}

	community = Community{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: ctx.BlockHeader().Time,
	}

	k.setCommunity(ctx, community)
	logger(ctx).Info(fmt.Sprintf("Created %s", community))

	return
}

// Community returns a community by its ID
func (k Keeper) Community(ctx sdk.Context, id string) (community Community, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	communityBytes := store.Get(key(id))
	if communityBytes == nil {
		return community, ErrCommunityNotFound(community.ID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(communityBytes, &community)

	return community, nil
}

// Communities gets all communities from the KVStore
func (k Keeper) Communities(ctx sdk.Context) (communities []Community) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, CommunityKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var community Community
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &community)
		communities = append(communities, community)
	}

	return
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	// if already present, don't add again
	for _, currentAdmin := range params.CommunityAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.CommunityAdmins = append(params.CommunityAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.CommunityAdmins {
		if currentAdmin.Equals(admin) {
			params.CommunityAdmins = append(params.CommunityAdmins[:i], params.CommunityAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) validateParams(ctx sdk.Context, id, name, description string, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)
	if len(id) < params.MinIDLength || len(id) > params.MaxIDLength {
		err = ErrInvalidCommunityMsg(
			fmt.Sprintf("ID must be between %d-%d chars in length", params.MinIDLength, params.MaxIDLength),
		)
	}
	if len(name) < params.MinNameLength || len(name) > params.MaxNameLength {
		err = ErrInvalidCommunityMsg(
			fmt.Sprintf("Name must be between %d-%d chars in length", params.MinNameLength, params.MaxNameLength),
		)
	}
	if len(description) > params.MaxDescriptionLength {
		err = ErrInvalidCommunityMsg(
			fmt.Sprintf("Description must be less than %d chars in length", params.MaxDescriptionLength),
		)
	}

	if !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	return
}

func (k Keeper) setCommunity(ctx sdk.Context, community Community) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(community)
	store.Set(key(community.ID), bz)
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).CommunityAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
