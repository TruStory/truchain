package slashing

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keys for slashing store
// Items are stored with the following key: values
//
// - 0x00<slashID_Bytes>: Slash_Bytes
// - 0x01: nextSlashID_Bytes
// - 0x02<stakeID_Bytes>: slashCount_Bytes
//
// - 0x10<stakeID_Bytes><slashID_Bytes>: slashID_Bytes
// - 0x11<creator_Bytes><slashID_Bytes>: slashID_Bytes
var (
	SlashesKeyPrefix = []byte{0x00}
	SlashIDKey       = []byte{0x01}
	SlashCountPrefix = []byte{0x02}

	StakeSlashesPrefix   = []byte{0x10}
	CreatorSlashesPrefix = []byte{0x11}
)

// key for getting a specific slash from the store
func key(claimID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, claimID)
	return append(SlashesKeyPrefix, bz...)
}

// stakeSlashesKey gets the first part of the stake slashes key based on the stakeID
func stakeSlashesKey(stakeID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, stakeID)
	return append(StakeSlashesPrefix, bz...)
}

// stakeSlashKey key of a specific stake <-> slash association from the store
func stakeSlashKey(stakeID uint64, slashID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, slashID)
	return append(stakeSlashesKey(stakeID), bz...)
}

// creatorSlashesKey gets the first part of the creator's slashes based on the creator
func creatorSlashesKey(creator sdk.AccAddress) []byte {
	return append(CreatorSlashesPrefix, creator.Bytes()...)
}

// creatorSlashKey key of the specific creator <-> slash association from the store
func creatorSlashKey(creator sdk.AccAddress, slashID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, slashID)
	return append(creatorSlashesKey(creator), bz...)
}

// slashCountKey gets the first part of the slash count key
func slashCountKey(stakeID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, stakeID)
	return append(SlashCountPrefix, bz...)
}
