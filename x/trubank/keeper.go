package trubank

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	cat "github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for bank
	StoreKey = "trubank"
)

// ReadKeeper defines a module interface that facilitates read only access
type ReadKeeper interface {
	app.ReadKeeper

	TransactionsByCreator(ctx sdk.Context, creator sdk.AccAddress) (transactions []Transaction, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
type WriteKeeper interface {
	ReadKeeper

	AddCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64, status Status) (coins sdk.Coins, err sdk.Error)
	SubtractCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64, status Status) (coins sdk.Coins, err sdk.Error)
	MintAndAddCoin(ctx sdk.Context, creator sdk.AccAddress, catID int64, amt sdk.Int) (sdk.Coins, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	bankKeeper     bank.Keeper
	categoryKeeper cat.WriteKeeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	bankKeeper bank.Keeper,
	categoryKeeper cat.WriteKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		bankKeeper,
		categoryKeeper,
	}
}

// AddCoin wraps around adding coins via the bank keeper and adds the transaction
func (k Keeper) AddCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64, status Status) (coins sdk.Coins, err sdk.Error) {
	coins, _, err = k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{coin})

	transaction := Transaction{
		ID:              k.GetNextID(ctx),
		TransactionType: transactionType,
		ReferenceID:     referenceID,
		Creator:         creator,
		Status:          status,
		Timestamp:       app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setTransaction(ctx, transaction)

	return coins, err
}

// SubtractCoin wraps around adding coins via the bank keeper and adds the transaction
func (k Keeper) SubtractCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64, status Status) (coins sdk.Coins, err sdk.Error) {
	coins, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{coin})

	transaction := Transaction{
		ID:              k.GetNextID(ctx),
		TransactionType: transactionType,
		ReferenceID:     referenceID,
		Creator:         creator,
		Status:          status,
		Timestamp:       app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setTransaction(ctx, transaction)

	return coins, err
}

// MintAndAddCoin adds coins to a user's account and to the total category supply
func (k Keeper) MintAndAddCoin(
	ctx sdk.Context,
	creator sdk.AccAddress,
	catID int64,
	amt sdk.Int) (sdk.Coins, sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	cat, err := k.categoryKeeper.GetCategory(ctx, catID)
	if err != nil {
		return sdk.Coins{}, err
	}

	coin := sdk.NewCoin(cat.Slug, amt)

	coins, _, err := k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{coin})
	if err != nil {
		return nil, ErrTransferringCoinsToUser(creator)
	}

	err = k.categoryKeeper.AddToTotalCred(ctx, catID, coin)
	if err != nil {
		return sdk.Coins{}, ErrTransferringCoinsToCategory(catID)
	}

	logger.Info(fmt.Sprintf("Updated total supply of category %d to %s", catID, amt))

	return coins, nil
}

// NewCategoryCoin creates a new category coin type
func NewCategoryCoin(toDenom string, from sdk.Coin) sdk.Coin {
	rate := exchangeCoinsBetweenDenoms(from, toDenom)

	return sdk.NewCoin(
		toDenom,
		sdk.NewDecFromInt(from.Amount).Mul(rate).TruncateInt())
}

// exchangeCoinsBetweenDenoms exchanges coins from trustake to cred
// TODO [Shane]: https://github.com/TruStory/truchain/issues/21
func exchangeCoinsBetweenDenoms(from sdk.Coin, toDenom string) sdk.Dec {

	if from.Denom == toDenom {
		return sdk.NewDec(1)
	}
	return sdk.NewDec(1)
}

func (k Keeper) setTransaction(ctx sdk.Context, transaction Transaction) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(transaction.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(transaction))
}

// TransactionsByCreator returns all the transactions for a user
func (k Keeper) TransactionsByCreator(ctx sdk.Context, creator sdk.AccAddress) (transactions []Transaction, err sdk.Error) {
	// get store
	store := k.GetStore(ctx)

	// builds prefix "trubank:creator:"
	searchKey := fmt.Sprintf("%s:creator:", k.GetStoreKey().Name())
	searchPrefix := []byte(searchKey)

	// setup iterator
	iter := sdk.KVStorePrefixIterator(store, searchPrefix)
	defer iter.Close()

	// iterates through keyspace to find all stories
	for ; iter.Valid(); iter.Next() {
		var transaction Transaction
		k.GetCodec().MustUnmarshalBinaryLengthPrefixed(
			iter.Value(), &transaction)
		transactions = append(transactions, transaction)
	}

	return
}

// Transaction returns a single transaction from the K-V Store
func (k Keeper) Transaction(
	ctx sdk.Context, transactionID int64) (transaction Transaction, err sdk.Error) {

	store := k.GetStore(ctx)
	val := store.Get(k.GetIDKey(transactionID))
	if val == nil {
		return transaction, ErrTransactionNotFound(transactionID)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &transaction)

	return
}
