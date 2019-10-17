package staking

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// setClaimArgument sets a claim <-> argument association in the store
func (k Keeper) setClaimArgument(ctx sdk.Context, claimID, argumentID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(argumentID)
	k.store(ctx).Set(claimArgumentKey(claimID, argumentID), bz)
}

func (k Keeper) IterateClaimArguments(ctx sdk.Context, claimID uint64, cb func(argument Argument) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), claimArgumentsPrefix(claimID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var argumentID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &argumentID)
		arg, ok := k.Argument(ctx, argumentID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve argument with id %d", argumentID))
		}
		if cb(arg) {
			break
		}
	}
}

// setArgumentStake sets a argument <-> stake association in the store
func (k Keeper) setArgumentStake(ctx sdk.Context, argumentID, stakeID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stakeID)
	k.store(ctx).Set(argumentStakeKey(argumentID, stakeID), bz)
}

func (k Keeper) IterateArgumentStakes(ctx sdk.Context, argumentID uint64, cb func(stake Stake) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), argumentStakesPrefix(argumentID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		stake, ok := k.Stake(ctx, stakeID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve stake with id %d", stakeID))
		}
		if cb(stake) {
			break
		}
	}
}

// setCommunityStake sets a community <-> stake association in the store
func (k Keeper) setCommunityStake(ctx sdk.Context, communityID string, stakeID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stakeID)
	k.store(ctx).Set(communityStakeKey(communityID, stakeID), bz)
}

func (k Keeper) IterateCommunityStakes(ctx sdk.Context, communityID string, cb func(stake Stake) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), communityStakesPrefix(communityID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		stake, ok := k.Stake(ctx, stakeID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve stake with id %d", stakeID))
		}
		if cb(stake) {
			break
		}
	}
}

// setUserArgument sets a user <-> argument association in the store
func (k Keeper) setUserArgument(ctx sdk.Context, creator sdk.AccAddress, argumentID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(argumentID)
	k.store(ctx).Set(userArgumentKey(creator, argumentID), bz)
}

func (k Keeper) IterateUserArguments(ctx sdk.Context, creator sdk.AccAddress, cb func(argument Argument) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), userArgumentsPrefix(creator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var argumentID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &argumentID)
		arg, ok := k.Argument(ctx, argumentID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve argument with id %d", argumentID))
		}
		if cb(arg) {
			break
		}
	}
}

// setUserStake sets a user <-> stake association in the store
func (k Keeper) setUserStake(ctx sdk.Context, creator sdk.AccAddress, creationTime time.Time, stakeID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stakeID)
	k.store(ctx).Set(userStakeKey(creator, creationTime, stakeID), bz)
}

func (k Keeper) IterateUserStakes(ctx sdk.Context, creator sdk.AccAddress, cb func(stake Stake) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), userStakesPrefix(creator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		stake, ok := k.Stake(ctx, stakeID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve stake with id %d", stakeID))
		}
		if cb(stake) {
			break
		}
	}
}

func (k Keeper) setUserCommunityStake(ctx sdk.Context, creator sdk.AccAddress, communityID string, stakeID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stakeID)
	k.store(ctx).Set(userCommunityStakeKey(creator, communityID, stakeID), bz)
}

func (k Keeper) IterateUserCommunityStakes(ctx sdk.Context, creator sdk.AccAddress, communityID string, cb func(stake Stake) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), userCommunityStakesPrefix(creator, communityID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		stake, ok := k.Stake(ctx, stakeID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve stake with id %d", stakeID))
		}
		if cb(stake) {
			break
		}
	}
}

func (k Keeper) IterateAfterCreatedTimeUserStakes(ctx sdk.Context,
	creator sdk.AccAddress, createdTime time.Time,
	cb func(stake Stake) (stop bool)) {
	iterator := k.store(ctx).Iterator(userStakesCreatedTimePrefix(creator, createdTime),
		sdk.PrefixEndBytes(userStakesPrefix(creator)),
	)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		stake, ok := k.Stake(ctx, stakeID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve stake with id %d", stakeID))
		}
		if cb(stake) {
			break
		}
	}
}

// ActiveStakeQueueIterator returns an sdk.Iterator for all the stakes in the Active Queue that expire by endTime
func (k Keeper) ActiveStakeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(ActiveStakeQueuePrefix, sdk.PrefixEndBytes(activeStakeByTimeKey(endTime)))
}

func (k Keeper) IterateActiveStakeQueue(ctx sdk.Context, endTime time.Time, cb func(stake Stake) (stop bool)) {
	iterator := k.ActiveStakeQueueIterator(ctx, endTime)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		stake, ok := k.Stake(ctx, stakeID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve stake with id %d", stakeID))
		}
		if cb(stake) {
			break
		}
	}
}

type userEarnedCoinsCallback func(address sdk.AccAddress, coins sdk.Coins) (stop bool)

func (k Keeper) IterateUserEarnedCoins(ctx sdk.Context, cb userEarnedCoinsCallback) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), EarnedCoinsKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		address := splitKeyWithAddress(iterator.Key())
		coins := sdk.NewCoins()

		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &coins)
		if cb(address, coins) {
			break
		}
	}
}
