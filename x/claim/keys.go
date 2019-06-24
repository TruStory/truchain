package claim

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keys for claim store
// Items are stored with the following key: values
//
// - 0x00<claimID_Bytes>: Claim_Bytes
// - 0x01: nextClaimID_Bytes
//
// - 0x10<communityID_Bytes><claimID_Bytes>: claimID_Bytes
// - 0x11<creator_Bytes><claimID_Bytes>: claimID_Bytes
// - 0x12<createdTime_Bytes><claimID_Bytes>: claimID_Bytes
var (
	ClaimsKeyPrefix = []byte{0x00}
	ClaimIDKey      = []byte{0x01}

	CommunityClaimsPrefix  = []byte{0x10}
	CreatorClaimsPrefix    = []byte{0x11}
	CreatedTimeClaimsPrefx = []byte{0x12}
)

// key for getting a specific claim from the store
func key(claimID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, claimID)
	return append(ClaimsKeyPrefix, bz...)
}

// communityClaimsKey gets the first part of the community claims key based on the communityID
func communityClaimsKey(communityID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, communityID)
	return append(CommunityClaimsPrefix, bz...)
}

// communityClaimKey key of a specific community <-> claim association from the store
func communityClaimKey(communityID uint64, claimID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, claimID)
	return append(communityClaimsKey(communityID), bz...)
}

func creatorClaimsKey(creator sdk.AccAddress) []byte {
	return append(CreatorClaimsPrefix, creator.Bytes()...)
}

func creatorClaimKey(creator sdk.AccAddress, claimID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, claimID)
	return append(creatorClaimsKey(creator), bz...)
}

func createdTimeClaimsKey(createdTime time.Time) []byte {
	return append(CreatedTimeClaimsPrefx, sdk.FormatTimeBytes(createdTime)...)
}

func createdTimeClaimKey(createdTime time.Time, claimID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, claimID)

	return append(createdTimeClaimsKey(createdTime), bz...)
}
