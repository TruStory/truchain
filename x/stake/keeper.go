package stake

import (
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/trubank"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// StoreKey is string representation of the store key
	StoreKey = "stake"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	truBankKeeper trubank.WriteKeeper
	paramStore    params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(truBankKeeper trubank.WriteKeeper, paramStore params.Subspace) Keeper {
	return Keeper{
		truBankKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

// DistributeInterest distributes interest to backers and challengers
// func (k Keeper) DistributeInterest(ctx sdk.Context, storyID int64) sdk.Error {
// 	logger := ctx.Logger().With("module", "expiration")

// 	for _, backing := range backings {
// 		// give the principal back to the user (in trustake)
// 		_, _, err := k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Amount()})
// 		if err != nil {
// 			return err
// 		}
// 		// give the interest earned to the user (in cred)
// 		period := story.ExpireTime.Sub(backing.Timestamp.CreatedTime)
// 		maxPeriod := story.ExpireTime.Sub(story.Timestamp.CreatedTime)
// 		logger.Info(fmt.Sprintf(
// 			"Backing period: %s, max period: %s", period, maxPeriod))

// 		interest := k.interest(ctx, backing.Amount(), period, maxPeriod, denom)
// 		_, _, err = k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{interest})
// 		if err != nil {
// 			return err
// 		}

// 		logger.Info(fmt.Sprintf(
// 			"Distributed earnings of %s with interest of %s to %s",
// 			backing.Amount().String(),
// 			backing.Interest.String(),
// 			backing.Creator().String()))
// 	}

// DistributePrincipalAndInterest distributes funds to backers and challengers
func (k Keeper) DistributePrincipalAndInterest(
	ctx sdk.Context, stakes []app.Voter, categoryID int64) sdk.Error {

	for _, stake := range stakes {
		// give principal back to user in trustake
		_, err := k.truBankKeeper.AddCoin(ctx, stake.Creator(), stake.Amount())
		if err != nil {
			return err
		}
		// mint interest earned and distribute
		// period := ctx.BlockHeader().Time.Sub(stake.)
	}

	return nil
}

// ValidateArgument validates the length of an argument
func (k Keeper) ValidateArgument(ctx sdk.Context, argument string) sdk.Error {
	len := len([]rune(argument))
	minArgumentLength := k.GetParams(ctx).MinArgumentLength
	maxArgumentLength := k.GetParams(ctx).MaxArgumentLength

	if len > 0 && (len < minArgumentLength) {
		return ErrArgumentTooShortMsg(argument, minArgumentLength)
	}

	if len > 0 && (len > maxArgumentLength) {
		return ErrArgumentTooLongMsg(maxArgumentLength)
	}

	return nil
}

func (k Keeper) interest(
	ctx sdk.Context,
	amount sdk.Coin,
	period time.Duration,
	maxPeriod time.Duration,
	credDenom string) sdk.Coin {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22

	totalSupply := sdk.NewDec(1000000000000000)

	// inputs
	maxAmount := totalSupply
	amountWeight := k.GetParams(ctx).AmountWeight
	periodWeight := k.GetParams(ctx).PeriodWeight
	maxInterestRate := k.GetParams(ctx).MaxInterestRate

	// type cast values to unitless decimals for math operations
	periodDec := sdk.NewDec(int64(period))
	maxPeriodDec := sdk.NewDec(int64(maxPeriod))
	amountDec := sdk.NewDecFromInt(amount.Amount)

	// normalize amount and period to 0 - 1
	normalizedAmount := amountDec.Quo(maxAmount)
	normalizedPeriod := periodDec.Quo(maxPeriodDec)

	// apply weights to normalized amount and period
	weightedAmount := normalizedAmount.Mul(amountWeight)
	weightedPeriod := normalizedPeriod.Mul(periodWeight)

	// calculate interest
	interestRate := maxInterestRate.Mul(weightedAmount.Add(weightedPeriod))
	// convert rate to a value
	minInterestRate := k.GetParams(ctx).MinInterestRate
	if interestRate.LT(minInterestRate) {
		interestRate = minInterestRate
	}
	interest := amountDec.Mul(interestRate)

	// return cred coin with rounded interest
	cred := sdk.NewCoin(credDenom, interest.RoundInt())

	return cred
}
