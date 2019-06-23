package claim

import (
	"net/url"

	"github.com/TruStory/truchain/x/community"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the module
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	accountKeeper   AccountKeeper
	communityKeeper community.Keeper
}

// NewKeeper creates a new claim keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, accountKeeper AccountKeeper, communityKeeper community.Keeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper,
		communityKeeper,
	}
}

// SubmitClaim creates a new claim in the claim key-value store
func (k Keeper) SubmitClaim(ctx sdk.Context, body string, communityID uint64,
	creator sdk.AccAddress, source url.URL) (claim Claim, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}
	if k.accountKeeper.IsJailed(ctx, creator) {
		return claim, ErrCreatorJailed(creator)
	}
	community, err := k.communityKeeper.Community(ctx, communityID)
	if err != nil {
		return claim, ErrInvalidCommunityID(community.ID)
	}

	claimID, err := k.claimID(ctx)
	if err != nil {
		return
	}
	claim = NewClaim(claimID, communityID, body, creator, source,
		ctx.BlockHeader().Time,
	)
	k.setClaim(ctx, claim)
	k.setClaimID(ctx, claimID+1)
	k.setCommunityClaim(ctx, claim.CommunityID, claimID)

	logger(ctx).Info("Created " + claim.String())
	// fmt.Println(claim)

	return claim, nil
}

// Claim gets a single claim by its ID
func (k Keeper) Claim(ctx sdk.Context, id uint64) (claim Claim, ok bool) {
	store := ctx.KVStore(k.storeKey)
	claimBytes := store.Get(key(id))
	if claimBytes == nil {
		return claim, false
	}
	k.codec.UnmarshalBinaryLengthPrefixed(claimBytes, &claim)

	return claim, true
}

// Claims gets all the claims
func (k Keeper) Claims(ctx sdk.Context) (claims Claims) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ClaimsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claim Claim
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &claim)
		claims = append(claims, claim)
	}

	return
}

// CommunityClaims gets all the claims for a given community
func (k Keeper) CommunityClaims(ctx sdk.Context, communityID uint64) (claims Claims) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, communityClaimsKey(communityID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claimID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &claimID)
		claim, ok := k.Claim(ctx, claimID)
		if ok {
			claims = append(claims, claim)
		}
	}

	return
}

// AddBackingStake adds a stake amount to the total backing amount
func (k Keeper) AddBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.TotalBacked.Add(stake)
	claim.TotalStakers++
	k.setClaim(ctx, claim)

	return nil
}

// AddChallengeStake adds a stake amount to the total challenge amount
func (k Keeper) AddChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.TotalChallenged.Add(stake)
	claim.TotalStakers++
	k.setClaim(ctx, claim)

	return nil
}

func (k Keeper) validateLength(ctx sdk.Context, body string) sdk.Error {
	var minClaimLength int
	var maxClaimLength int

	k.paramStore.Get(ctx, KeyMinClaimLength, &minClaimLength)
	k.paramStore.Get(ctx, KeyMaxClaimLength, &maxClaimLength)

	len := len([]rune(body))
	if len < minClaimLength {
		return ErrInvalidBodyTooShort(body)
	}
	if len > maxClaimLength {
		return ErrInvalidBodyTooLong()
	}

	return nil
}

// claimID gets the highest claim ID
func (k Keeper) claimID(ctx sdk.Context) (claimID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(ClaimIDKey)
	if bz == nil {
		return 0, ErrUnknownClaim(claimID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &claimID)
	return claimID, nil
}

// set the claim ID
func (k Keeper) setClaimID(ctx sdk.Context, claimID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claimID)
	store.Set(ClaimIDKey, bz)
}

// setClaim sets a claim in store
func (k Keeper) setClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(key(claim.ID), bz)
}

// setCommunityClaim sets a community <-> claim association in store
func (k Keeper) setCommunityClaim(ctx sdk.Context, communityID, claimID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claimID)
	store.Set(communityClaimKey(communityID, claimID), bz)
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
