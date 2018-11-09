package challenge

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// generates "games:id:5:challenges:user:[Address]"
func (k Keeper) challengeByGameIDKey(
	ctx sdk.Context, gameID int64, user sdk.AccAddress) []byte {

	key := fmt.Sprintf(
		"%s:id:%d:%s:user:%s",
		k.gameKeeper.GetStoreKey().Name(),
		gameID,
		k.GetStoreKey().Name(),
		user.String())

	return []byte(key)
}

// generates "games:id:5:challenges:user:"
func (k Keeper) challengeByGameIDSubspace(ctx sdk.Context, gameID int64) []byte {
	key := fmt.Sprintf(
		"%s:id:%d:%s:user:",
		k.gameKeeper.GetStoreKey().Name(),
		gameID,
		k.GetStoreKey().Name())

	return []byte(key)
}
