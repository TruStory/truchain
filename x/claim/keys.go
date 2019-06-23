package claim

import "encoding/binary"

// TODO
// user address -> total number of claims

// Keys for claim store
// Items are stored with the following key: values
//
// - 0x00<claimID_Bytes>: Claim_Bytes
// - 0x01: nextClaimID_Bytes
//
// - 0x10<communityID_Bytes><claimID_Bytes>: claimID_Bytes
var (
	ClaimsKeyPrefix = []byte{0x00}
	ClaimIDKey      = []byte{0x01}

	CommunityClaimsPrefix = []byte{0x10}
	CreatorClaimsPrefix   = []byte{0x11}
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
