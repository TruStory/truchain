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
	// ListStoreKey is string representation of the store key for bank
	ListStoreKey = "trubankList"
)

// ReadKeeper defines a module interface that facilitates read only access
type ReadKeeper interface {
	app.ReadKeeper

	TransactionsByCreator(ctx sdk.Context, creator sdk.AccAddress) (transactions []Transaction, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
type WriteKeeper interface {
	ReadKeeper

	AddCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64) (coins sdk.Coins, err sdk.Error)
	SubtractCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64) (coins sdk.Coins, err sdk.Error)
	MintAndAddCoin(ctx sdk.Context, creator sdk.AccAddress, catID int64, storyID int64, transactionType TransactionType, amt sdk.Int) (sdk.Coins, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	bankKeeper     bank.Keeper
	categoryKeeper cat.WriteKeeper
	trubankList    app.UserList // transactions <-> user mappings
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
		app.NewUserList(storeKey),
	}
}

// AddCoin wraps around adding coins via the bank keeper and adds the transaction
func (k Keeper) AddCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64) (coins sdk.Coins, err sdk.Error) {
	coins, _, err = k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{coin})

	if coin.IsZero() {
		return sdk.Coins{}, nil
	}

	transaction := Transaction{
		ID:              k.GetNextID(ctx),
		TransactionType: transactionType,
		ReferenceID:     referenceID,
		GroupID:         storyID,
		Amount:          coin,
		Creator:         creator,
		Timestamp:       app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setTransaction(ctx, transaction)
	k.trubankList.AppendToUser(ctx, k, creator, transaction.ID)

	return coins, err
}

// SubtractCoin wraps around subtracting coins via the bank keeper and adds the transaction
func (k Keeper) SubtractCoin(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, storyID int64, transactionType TransactionType, referenceID int64) (coins sdk.Coins, err sdk.Error) {
	coins, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{coin})

	coin.Amount = coin.Amount.Mul(sdk.NewInt(-1))
	if coin.IsZero() {
		return sdk.Coins{}, nil
	}

	transaction := Transaction{
		ID:              k.GetNextID(ctx),
		TransactionType: transactionType,
		ReferenceID:     referenceID,
		GroupID:         storyID,
		Amount:          coin,
		Creator:         creator,
		Timestamp:       app.NewTimestamp(ctx.BlockHeader()),
	}

	k.setTransaction(ctx, transaction)
	k.trubankList.AppendToUser(ctx, k, creator, transaction.ID)

	return coins, err
}

// MintAndAddCoin adds coins to a user's account and to the total category supply
func (k Keeper) MintAndAddCoin(
	ctx sdk.Context,
	creator sdk.AccAddress,
	catID int64,
	storyID int64,
	transactionType TransactionType,
	amt sdk.Int) (sdk.Coins, sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	if amt.IsZero() {
		return sdk.Coins{}, nil
	}

	cat, err := k.categoryKeeper.GetCategory(ctx, catID)
	if err != nil {
		return sdk.Coins{}, err
	}

	coin := sdk.NewCoin(cat.Slug, amt)

	coins, err := k.AddCoin(ctx, creator, coin, storyID, transactionType, 0)
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

	err = k.trubankList.MapByUser(ctx, k, creator, func(transactionID int64) sdk.Error {
		transaction, err := k.Transaction(ctx, transactionID)
		if err != nil {
			return err
		}

		transactions = append(transactions, transaction)
		return nil
	})

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
