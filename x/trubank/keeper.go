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
}

// WriteKeeper defines a module interface that facilities write only access
type WriteKeeper interface {
	ReadKeeper

	MintAndAddCoin(ctx sdk.Context, creator sdk.AccAddress, catID int64, amt sdk.Int) (sdk.Coins, sdk.Error)
	NewCategoryCoin(toDenom string, from sdk.Coin) sdk.Coin
	ExchangeCoinsBetweenDenoms(from sdk.Coin, toDenom string) sdk.Dec
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

// MintAndAddCoin adds coins to a user's account and to the total category supply
func (k Keeper) MintAndAddCoin(
	ctx sdk.Context,
	creator sdk.AccAddress,
	catID int64,
	amt sdk.Int) (sdk.Coins, sdk.Error) {

	logger := ctx.Logger().With("module", "trubank")

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
	rate := ExchangeCoinsBetweenDenoms(from, toDenom)

	return sdk.NewCoin(
		toDenom,
		sdk.NewDecFromInt(from.Amount).Mul(rate).TruncateInt())
}

// TODO [Shane]: https://github.com/TruStory/truchain/issues/21
// ExchangeRate exchanges coins from trustake to cred
func ExchangeCoinsBetweenDenoms(from sdk.Coin, toDenom string) sdk.Dec {

	if from.Denom == toDenom {
		return sdk.NewDec(1)
	}
	return sdk.NewDec(1)
}
