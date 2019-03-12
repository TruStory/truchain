package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UserList defines a list of users associated with a type.
// Users could be backers, challengers, or voters.
// Store layout:
// [foreignStoreKey]:id:[keyID]:[storeKey]:users:[user] -> [valueID]
type UserList struct {
	foreignStoreKey sdk.StoreKey
}

// NewUserList creates a new `UserList`
func NewUserList(foreignStoreKey sdk.StoreKey) UserList {
	return UserList{foreignStoreKey}
}

// Append adds a new key <-> value association
func (l UserList) Append(
	ctx sdk.Context, k WriteKeeper, keyID int64, user sdk.AccAddress, valueID int64) {
	k.GetStore(ctx).Set(
		l.typeByUserKey(ctx, k, keyID, user),
		k.GetCodec().MustMarshalBinaryBare(valueID))
}

// Delete deletes a  key <-> value association from the underlying store.
func (l UserList) Delete(
	ctx sdk.Context, k WriteKeeper, keyID int64, user sdk.AccAddress) {
	k.GetStore(ctx).Delete(l.typeByUserKey(ctx, k, keyID, user))
}

// Get gets a saved value id for the given key
func (l UserList) Get(
	ctx sdk.Context, k WriteKeeper, keyID int64, user sdk.AccAddress) (valueID int64) {

	bz := k.GetStore(ctx).Get(l.typeByUserKey(ctx, k, keyID, user))
	if bz == nil {
		return 0
	}
	k.GetCodec().MustUnmarshalBinaryBare(bz, &valueID)

	return valueID
}

// Includes returns true if the given key is found
func (l UserList) Includes(
	ctx sdk.Context, k WriteKeeper, keyID int64, user sdk.AccAddress) bool {

	return l.Get(ctx, k, keyID, user) > 0
}

// Map applies a function across the subspace of users on a key
func (l UserList) Map(
	ctx sdk.Context, k WriteKeeper, keyID int64, fn func(int64) sdk.Error) sdk.Error {

	// get store
	store := k.GetStore(ctx)

	// builds prefix
	prefix := l.typeByUserSubspaceKey(ctx, k, keyID)

	// iterates through keyspace to find all value ids
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var id int64
		k.GetCodec().MustUnmarshalBinaryBare(iter.Value(), &id)
		if err := fn(id); err != nil {
			return err
		}
	}

	return nil
}

// MapByUser applies a function
func (l UserList) MapByUser(
	ctx sdk.Context, k WriteKeeper, keyID int64, user sdk.AccAddress, fn func(int64) sdk.Error) sdk.Error {

	// get store
	store := k.GetStore(ctx)

	// builds prefix
	prefix := l.typeByUserKey(ctx, k, keyID, user)

	// iterates through keyspace to find all value ids
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var id int64
		k.GetCodec().MustUnmarshalBinaryBare(iter.Value(), &id)
		if err := fn(id); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================

// generates key "[foreignStoreKey]:id:[keyID]:[storeKey]:users:[user]"
func (l UserList) typeByUserKey(
	ctx sdk.Context, k WriteKeeper, keyID int64, user sdk.AccAddress) []byte {

	key := fmt.Sprintf(
		"%s:id:%d:%s:user:%s",
		l.foreignStoreKey.Name(),
		keyID,
		k.GetStoreKey().Name(),
		user.String())

	return []byte(key)
}

// generates key "[foreignStoreKey]:id:[keyID]:[storeKey]:users:"
func (l UserList) typeByUserSubspaceKey(
	ctx sdk.Context, k WriteKeeper, keyID int64) []byte {

	key := fmt.Sprintf(
		"%s:id:%d:%s:user:",
		l.foreignStoreKey.Name(),
		keyID,
		k.GetStoreKey().Name())

	return []byte(key)
}
