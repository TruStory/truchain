package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO:
// make generic SubspaceList
// VoterList embeds SubspaceList

// VoterList defines a list of voters associated with a validation game.
// Voters could be Backers, Challengers, or actual Voters.
type VoterList struct {
	gameStoreKey sdk.StoreKey
}

// NewVoterList creates a new `VoterList`
func NewVoterList(
	gameStoreKey sdk.StoreKey) VoterList {

	return VoterList{gameStoreKey}
}

// Append adds a new game <-> vote association
func (l VoterList) Append(
	ctx sdk.Context, k WriteKeeper, gameID int64, user sdk.AccAddress, voteID int64) {

	k.GetStore(ctx).Set(
		l.gameByUserKey(ctx, k, gameID, user),
		k.GetCodec().MustMarshalBinary(voteID))
}

// Map applies a function across the subspace of voters on a game
func (l VoterList) Map(
	ctx sdk.Context, k WriteKeeper, gameID int64, fn func(int64)) {

	store := k.GetStore(ctx)

	// builds prefix of from "game:id:5:votes:user:"
	prefix := l.gameByUserSubspace(ctx, k, gameID)

	// iterates through keyspace to find all votes on a game
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var id int64
		k.GetCodec().MustUnmarshalBinary(iter.Value(), &id)
		fn(id)
	}
}

// TODO: Also implement Reduce, and use for tallying votes?

// generates "games:id:5:votes:user:[Address]"
func (l VoterList) gameByUserKey(
	ctx sdk.Context, k WriteKeeper, gameID int64, user sdk.AccAddress) []byte {

	key := fmt.Sprintf(
		"%s:id:%d:%s:user:%s",
		l.gameStoreKey.Name(),
		gameID,
		k.GetStoreKey().Name(),
		user.String())

	return []byte(key)
}

// generates "games:id:5:votes:user:"
func (l VoterList) gameByUserSubspace(
	ctx sdk.Context, k WriteKeeper, gameID int64) []byte {

	key := fmt.Sprintf(
		"%s:id:%d:%s:user:",
		l.gameStoreKey.Name(),
		gameID,
		k.GetStoreKey().Name())

	return []byte(key)
}
