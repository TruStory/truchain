package staking

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Define keys
var (
	StakesKeyPrefix      = []byte{0x00}
	ArgumentsKeyPrefix   = []byte{0x01}
	EarnedCoinsKeyPrefix = []byte{0x02}

	// ID Keys
	StakeIDKey    = []byte{0x10}
	ArgumentIDKey = []byte{0x11}

	// AssociationKeys
	ClaimArgumentsKeyPrefix      = []byte{0x20}
	ArgumentStakesKeyPrefix      = []byte{0x21}
	UserArgumentsKeyPrefix       = []byte{0x22}
	UserStakesKeyPrefix          = []byte{0x23}
	CommunityStakesKeyPrefix     = []byte{0x24}
	UserCommunityStakesKeyPrefix = []byte{0x25}

	// Queue
	ActiveStakeQueuePrefix = []byte{0x40}
)

// stakeKey gets a key for a stake.
// 0x00<stake_id>
func stakeKey(id uint64) []byte {
	return buildKey(StakesKeyPrefix, id)
}

// argumentKey gets a key for an argument
// 0x01<argument_id>
func argumentKey(id uint64) []byte {
	return buildKey(ArgumentsKeyPrefix, id)
}

// 0x02<user>
func userEarnedCoinsKey(user sdk.AccAddress) []byte {
	return append(EarnedCoinsKeyPrefix, user.Bytes()...)
}

func splitKeyWithAddress(key []byte) (addr sdk.AccAddress) {
	if len(key[1:]) != sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8+sdk.AddrLen))
	}
	addr = sdk.AccAddress(key[1:])
	return
}

// ClaimArgumentsPrefix
// 0x20<claim_id>
func claimArgumentsPrefix(claimID uint64) []byte {
	return buildKey(ClaimArgumentsKeyPrefix, claimID)
}

// claimArgumentKey builds the key for claim->argument association
func claimArgumentKey(claimID, argumentID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(argumentID)
	return append(claimArgumentsPrefix(claimID), bz...)
}

// ArgumentStakesPrefix
// 0x21<argument_id>
func argumentStakesPrefix(argumentID uint64) []byte {
	return buildKey(ArgumentStakesKeyPrefix, argumentID)
}

// argumentStakeKey builds the key for argument->stake association
func argumentStakeKey(argumentID, stakeID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(stakeID)
	return append(argumentStakesPrefix(argumentID), bz...)
}

// userArgumentsPrefix
// 0x22<creator>
func userArgumentsPrefix(creator sdk.AccAddress) []byte {
	return append(UserArgumentsKeyPrefix, creator.Bytes()...)
}

// userArgumentKey builds the key for user->argument association
func userArgumentKey(creator sdk.AccAddress, argumentID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(argumentID)
	return append(userArgumentsPrefix(creator), bz...)
}

// userStakesPrefix
// 0x23<creator>
func userStakesPrefix(creator sdk.AccAddress) []byte {
	return append(UserStakesKeyPrefix, creator.Bytes()...)
}

// 0x24<community_id>
func communityStakesPrefix(communityID string) []byte {
	return append(CommunityStakesKeyPrefix, []byte(communityID)...)
}

// 0x24<community_id><stake_id>
func communityStakeKey(communityID string, stakeID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(stakeID)
	return append(communityStakesPrefix(communityID), bz...)
}

func userStakesCreatedTimePrefix(creator sdk.AccAddress, createdTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(createdTime)
	return append(userStakesPrefix(creator), bz...)
}

// userStakeKey builds the key for <user><creationTime><stakeID>->stake association
func userStakeKey(creator sdk.AccAddress, createdTime time.Time, stakeID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(stakeID)
	return append(userStakesCreatedTimePrefix(creator, createdTime), bz...)
}

// 0x25<creator>
func userCommunityPrefix(creator sdk.AccAddress) []byte {
	return append(UserCommunityStakesKeyPrefix, creator.Bytes()...)
}

// 0x25<creator><community_id>
func userCommunityStakesPrefix(creator sdk.AccAddress, communityID string) []byte {
	return append(userCommunityPrefix(creator), []byte(communityID)...)
}

// 0x25<creator><community_id><stake_id>
func userCommunityStakeKey(creator sdk.AccAddress, communityID string, stakeID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(stakeID)
	return append(userCommunityStakesPrefix(creator, communityID), bz...)
}

// activeStakeQueueKey
// 0x40<end_time><stake_id>
func activeStakeQueueKey(stakeID uint64, endTime time.Time) []byte {
	bz := sdk.Uint64ToBigEndian(stakeID)
	return append(activeStakeByTimeKey(endTime), bz...)
}

// activeStakeByTimeKey gets the active proposal queue key by endTime
func activeStakeByTimeKey(endTime time.Time) []byte {
	return append(ActiveStakeQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

func buildKey(prefix []byte, id uint64) []byte {
	bz := sdk.Uint64ToBigEndian(id)
	return append(prefix, bz...)
}
