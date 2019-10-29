package bank

import (
	"github.com/TruStory/truchain/x/distribution"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	app "github.com/TruStory/truchain/types"
)

// Keeper is the model object for the package bank module
type Keeper struct {
	storeKey     sdk.StoreKey
	codec        *codec.Codec
	paramStore   params.Subspace
	bankKeeper   bank.Keeper
	codespace    sdk.CodespaceType
	supplyKeeper supply.Keeper
}

// NewKeeper creates a bank keeper.
func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey, bankKeeper bank.Keeper,
	paramStore params.Subspace, codespace sdk.CodespaceType, supplyKeeper supply.Keeper) Keeper {
	return Keeper{
		storeKey:     storeKey,
		codec:        codec,
		bankKeeper:   bankKeeper,
		paramStore:   paramStore.WithKeyTable(ParamKeyTable()),
		codespace:    codespace,
		supplyKeeper: supplyKeeper,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// AddCoin adds a coin to an address and adds the transaction to the association list.
func (k Keeper) AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
	referenceID uint64, txType TransactionType, txSetters ...TransactionSetter) (sdk.Coins, sdk.Error) {
	tx := Transaction{}
	for _, setter := range txSetters {
		setter(&tx)
	}
	if !txType.AllowedForAddition() {
		return sdk.Coins{}, ErrInvalidTransactionType(txType)
	}
	var err sdk.Error
	coins := sdk.Coins{amt}
	if tx.FromModuleAccount != "" {
		err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, tx.FromModuleAccount, addr, sdk.Coins{amt})
	}
	if tx.FromModuleAccount == "" {
		coins, err = k.bankKeeper.AddCoins(ctx, addr, sdk.Coins{amt})
	}

	if err != nil {
		return coins, err
	}
	transactionID, err := k.transactionID(ctx)
	if err != nil {
		return sdk.Coins{}, err
	}
	tx.ID = transactionID
	tx.Type = txType
	tx.ReferenceID = referenceID
	tx.Amount = amt
	tx.AppAccountAddress = addr
	tx.CreatedTime = ctx.BlockHeader().Time

	k.setTransaction(ctx, tx)
	k.setTransactionID(ctx, transactionID+1)
	k.setUserTransaction(ctx, addr, tx.CreatedTime, tx.ID)
	return coins, nil
}

// SubtractCoin subtracts a coin from an address and adds the transaction to the association list.
func (k Keeper) SubtractCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
	referenceID uint64, txType TransactionType, txSetters ...TransactionSetter) (sdk.Coins, sdk.Error) {
	tx := Transaction{}
	for _, setter := range txSetters {
		setter(&tx)
	}
	if !txType.AllowedForDeduction() {
		return sdk.Coins{}, ErrInvalidTransactionType(txType)
	}
	var err sdk.Error
	coins := sdk.Coins{amt}
	if tx.ToModuleAccount != "" {
		err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, tx.ToModuleAccount, sdk.Coins{amt})
	}
	if tx.ToModuleAccount == "" {
		coins, err = k.bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{amt})
	}
	if err != nil {
		return coins, err
	}

	transactionID, err := k.transactionID(ctx)
	if err != nil {
		return sdk.Coins{}, err
	}

	tx.ID = transactionID
	tx.Type = txType
	tx.ReferenceID = referenceID
	tx.Amount = amt
	tx.AppAccountAddress = addr
	tx.CreatedTime = ctx.BlockHeader().Time

	k.setTransaction(ctx, tx)
	k.setTransactionID(ctx, transactionID+1)
	k.setUserTransaction(ctx, addr, tx.CreatedTime, tx.ID)
	return coins, nil
}

// SafeSubtractCoin subtracts a coin without going below zero
func (k Keeper) SafeSubtractCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
	referenceID uint64, txType TransactionType, txSetters ...TransactionSetter) (sdk.Coins, sdk.Coin, sdk.Error) {

	if amt.IsNegative() {
		return sdk.Coins{}, sdk.Coin{}, sdk.ErrInvalidCoins("amount can't be negative")
	}

	adjustedCoin := amt
	balanceCoins := k.bankKeeper.GetCoins(ctx, addr)
	balance := balanceCoins.AmountOf(amt.Denom)
	if balance.LT(amt.Amount) {
		adjustedCoin = sdk.NewCoin(amt.Denom, balance)
	}

	if adjustedCoin.IsPositive() {
		coins, err := k.SubtractCoin(ctx, addr, adjustedCoin, referenceID, txType, txSetters...)
		if err != nil {
			return coins, adjustedCoin, err
		}
		return coins, adjustedCoin, nil
	} else {
		return balanceCoins, adjustedCoin, nil
	}
}

func (k Keeper) GetCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bankKeeper.GetCoins(ctx, address)
}

func (k Keeper) rewardBrokerAddress(ctx sdk.Context) sdk.AccAddress {
	address := sdk.AccAddress{}
	k.paramStore.GetIfExists(ctx, ParamKeyRewardBrokerAddress, &address)
	return address
}

func (k Keeper) sendGift(ctx sdk.Context,
	sender sdk.AccAddress, recipient sdk.AccAddress,
	amount sdk.Coin) sdk.Error {

	if !k.rewardBrokerAddress(ctx).Equals(sender) {
		return ErrInvalidRewardBrokerAddress(sender)
	}
	if amount.Denom != app.StakeDenom {
		return sdk.ErrInvalidCoins("Invalid denomination coin")
	}
	_, err := k.AddCoin(ctx, recipient, amount, 0, TransactionGift, FromModuleAccount(distribution.UserGrowthPoolName))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) payReward(ctx sdk.Context,
	sender sdk.AccAddress, recipient sdk.AccAddress,
	amount sdk.Coin, inviteID uint64) sdk.Error {
	if !k.rewardBrokerAddress(ctx).Equals(sender) {
		return ErrInvalidRewardBrokerAddress(sender)
	}
	_, err := k.AddCoin(ctx, recipient, amount, inviteID, TransactionRewardPayout,
		FromModuleAccount(distribution.UserRewardPoolName))
	if err != nil {
		return err
	}

	return nil
}

// Transactions gets all the transactions
func (k Keeper) Transactions(ctx sdk.Context) []Transaction {
	transactions := make([]Transaction, 0)
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), TransactionsKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var transaction Transaction
		k.codec.MustUnmarshalBinaryBare(iterator.Value(), &transaction)
		transactions = append(transactions, transaction)
	}
	return transactions
}

// TransactionsByAddress gets transactions for a given address and applies sent filters.
func (k Keeper) TransactionsByAddress(ctx sdk.Context, address sdk.AccAddress, filterSetters ...Filter) []Transaction {
	filters := GetFilters(filterSetters...)
	transactions := make([]Transaction, 0)
	filterByType := len(filters.TransactionTypes) > 0

	offsetCount := filters.Offset
	count := 0
	callbackFunc := func(tx Transaction) bool {
		if filterByType && !tx.Type.OneOf(filters.TransactionTypes) {
			return false
		}
		if offsetCount > 0 {
			offsetCount = offsetCount - 1
			return false
		}
		if filters.Limit > 0 && count == filters.Limit {
			return true
		}
		count++
		transactions = append(transactions, tx)
		return false
	}
	k.IterateUserTransactions(ctx, address, filters.SortOrder == SortDesc, callbackFunc)
	return transactions
}

func (k Keeper) transactionID(ctx sdk.Context) (uint64, sdk.Error) {
	id, err := k.getID(ctx, TransactionIDKey)
	if err != nil {
		return 0, ErrCodeUnknownTransaction(id)
	}
	return id, nil
}

func (k Keeper) setTransactionID(ctx sdk.Context, transactionID uint64) {
	k.setID(ctx, TransactionIDKey, transactionID)
}

func (k Keeper) setID(ctx sdk.Context, key []byte, length uint64) {
	b := k.codec.MustMarshalBinaryBare(length)
	k.store(ctx).Set(key, b)
}

func (k Keeper) getID(ctx sdk.Context, key []byte) (uint64, sdk.Error) {
	var id uint64
	b := k.store(ctx).Get(key)
	if b == nil {
		return 0, sdk.ErrInternal("unknown id")
	}
	k.codec.MustUnmarshalBinaryBare(b, &id)
	return id, nil
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) getTransaction(ctx sdk.Context, transactionID uint64) (Transaction, bool) {
	transaction := Transaction{}
	bz := k.store(ctx).Get(transactionKey(transactionID))
	if bz == nil {
		return transaction, false
	}
	k.codec.MustUnmarshalBinaryBare(bz, &transaction)
	return transaction, true
}

func (k Keeper) setTransaction(ctx sdk.Context, transaction Transaction) {
	bz := k.codec.MustMarshalBinaryBare(transaction)
	k.store(ctx).Set(transactionKey(transaction.ID), bz)
}
