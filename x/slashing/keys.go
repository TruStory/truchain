package slashing

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keys for slashing store
// Items are stored with the following key: values
//
// - 0x00<slashID>: Slash{}
// - 0x01: nextSlashID
// - 0x02<argumentID>: slashCount
//
// - 0x10<creator><slashID>: slashID
// - 0x11<argumentID><slashID>: slashID
// - 0x12<argumentID><slashCreator><slashID>: slashID
var (
	SlashesKeyPrefix = []byte{0x00}
	SlashIDKey       = []byte{0x01}
	SlashCountPrefix = []byte{0x02}

	CreatorSlashesPrefix  = []byte{0x10}
	ArgumentSlashesPrefix = []byte{0x11}
	ArgumentCreatorPrefix = []byte{0x12}
)

// key for getting a specific slash from the store
func key(claimID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, claimID)
	return append(SlashesKeyPrefix, bz...)
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

func argumentSlashPrefix(argumentID uint64) []byte {
	return append(ArgumentSlashesPrefix, sdk.Uint64ToBigEndian(argumentID)...)
}

func argumentSlashKey(argumentID, slashID uint64) []byte {
	return append(argumentSlashPrefix(argumentID), sdk.Uint64ToBigEndian(slashID)...)
}

func argumentSlasherPrefix(argumentID uint64, slasher sdk.AccAddress) []byte {
	return append(ArgumentCreatorPrefix, append(sdk.Uint64ToBigEndian(argumentID), slasher.Bytes()...)...)
}

func argumentSlasherSlashKey(argumentID uint64, slasher sdk.AccAddress, slashID uint64) []byte {
	return append(argumentSlasherPrefix(argumentID, slasher), sdk.Uint64ToBigEndian(slashID)...)
}
