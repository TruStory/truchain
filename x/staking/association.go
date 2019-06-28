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
		arg, ok := k.getArgument(ctx, argumentID)
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
		if cb(k.getStake(ctx, stakeID)) {
			break
		}
	}
}

// serUserArgument sets a user <-> argument association in the store
func (k Keeper) serUserArgument(ctx sdk.Context, creator sdk.AccAddress, argumentID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(argumentID)
	k.store(ctx).Set(userArgumentKey(creator, argumentID), bz)
}

func (k Keeper) IterateUserArguments(ctx sdk.Context, creator sdk.AccAddress, cb func(argument Argument) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), userArgumentsPrefix(creator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var argumentID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &argumentID)
		arg, ok := k.getArgument(ctx, argumentID)
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
		if cb(k.getStake(ctx, stakeID)) {
			break
		}
	}
}

func (k Keeper) IterateAfterCreatedTimeUserStakes(ctx sdk.Context,
	creator sdk.AccAddress, createdTime time.Time,
	cb func(stake Stake) (stop bool)) {
	iterator := k.store(ctx).Iterator(userStakesCreatedTimePrefix(creator, createdTime),
		sdk.PrefixEndBytes(UserStakesKeyPrefix),
	)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stakeID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stakeID)
		if cb(k.getStake(ctx, stakeID)) {
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
		stake := k.getStake(ctx, stakeID)
		if cb(stake) {
			break
		}
	}
}
