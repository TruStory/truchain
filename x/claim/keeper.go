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
	RecordKeeper

	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	accountKeeper   AccountKeeper
	communityKeeper community.Keeper
}

// NewKeeper creates a new claim keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, accountKeeper AccountKeeper, communityKeeper community.Keeper) Keeper {
	return Keeper{
		RecordKeeper{StoreKey: storeKey, Codec: codec},
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper,
		communityKeeper,
	}
}

// NewClaim creates a new claim in the claim key-value store
func (k Keeper) NewClaim(
	ctx sdk.Context,
	body string,
	communityID uint64,
	creator sdk.AccAddress,
	source url.URL) (claim Claim, err sdk.Error) {

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

	id := k.IncrementID(ctx)
	claim = NewClaim(
		id,
		communityID,
		body,
		creator,
		source,
		ctx.BlockHeader().Time)

	k.Set(ctx, id, claim)
	// add claim <-> community association
	k.Push(ctx, k.StoreKey, k.communityKeeper.StoreKey, id, communityID)

	logger(ctx).Info("Created " + claim.String())

	return claim, nil
}

// Claim gets a single claim by its ID
func (k Keeper) Claim(ctx sdk.Context, id uint64) (claim Claim, err sdk.Error) {
	err = k.Get(ctx, id, &claim)
	if err != nil {
		return
	}

	return
}

// Claims gets all the claims
func (k Keeper) Claims(ctx sdk.Context) (claims []Claim) {
	var claim Claim
	err := k.Each(ctx, func(val []byte) bool {
		k.codec.MustUnmarshalBinaryLengthPrefixed(val, &claim)
		claims = append(claims, claim)
		return true
	})
	if err != nil {
		return
	}

	return
}

// CommunityClaims gets all the claims for a given community
func (k Keeper) CommunityClaims(ctx sdk.Context, communityID uint64) (claims []Claim) {
	k.Map(ctx, k.communityKeeper.StoreKey, communityID, func(id uint64) {
		claim, err := k.Claim(ctx, id)
		if err != nil {
			panic(err)
		}
		claims = append(claims, claim)
	})

	return claims
}

// AddBackingStake adds a stake amount to the total backing amount
func (k Keeper) AddBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, err := k.Claim(ctx, id)
	if err != nil {
		return err
	}
	claim.TotalBacked.Add(stake)
	claim.TotalStakers++
	k.Set(ctx, id, claim)

	return nil
}

// AddChallengeStake adds a stake amount to the total challenge amount
func (k Keeper) AddChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, err := k.Claim(ctx, id)
	if err != nil {
		return err
	}
	claim.TotalChallenged.Add(stake)
	claim.TotalStakers++
	k.Set(ctx, id, claim)

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

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
