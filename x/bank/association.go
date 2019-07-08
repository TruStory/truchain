package bank

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// setUserTransaction sets a user <-> transaction association in the store
func (k Keeper) setUserTransaction(ctx sdk.Context, creator sdk.AccAddress, creationTime time.Time, transactionID uint64) {
	bz := k.codec.MustMarshalBinaryBare(transactionID)
	k.store(ctx).Set(userTransactionKey(creator, creationTime, transactionID), bz)
}

func (k Keeper) IterateUserTransactions(ctx sdk.Context, creator sdk.AccAddress, reverse bool, cb func(transaction Transaction) (stop bool)) {
	var iterator sdk.Iterator
	prefix := userTransactionsPrefix(creator)
	if !reverse {
		iterator = sdk.KVStorePrefixIterator(k.store(ctx), prefix)
	}

	if reverse {
		iterator = sdk.KVStoreReversePrefixIterator(k.store(ctx), prefix)
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var transactionID uint64
		k.codec.MustUnmarshalBinaryBare(iterator.Value(), &transactionID)
		transaction, ok := k.getTransaction(ctx, transactionID)
		if !ok {
			panic(fmt.Sprintf("unable to retrieve transaction with id %d", transactionID))
		}
		if cb(transaction) {
			break
		}
	}
}
