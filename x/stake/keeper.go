package stake

import (
	"fmt"
	"time"

	"github.com/TruStory/truchain/x/story"
	"github.com/davecgh/go-spew/spew"

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
	storyKeeper   story.ReadKeeper
	truBankKeeper trubank.WriteKeeper
	paramStore    params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storyKeeper story.ReadKeeper,
	truBankKeeper trubank.WriteKeeper,
	paramStore params.Subspace) Keeper {

	return Keeper{
		storyKeeper,
		truBankKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

// DistributePrincipalAndInterest distributes funds to backers and challengers
func (k Keeper) DistributePrincipalAndInterest(
	ctx sdk.Context, votes []app.Voter, categoryID int64) sdk.Error {

	logger := ctx.Logger().With("module", StoreKey)

	for _, vote := range votes {
		// give principal back to user in trustake
		_, err := k.truBankKeeper.AddCoin(ctx, vote.Creator(), vote.Amount())
		if err != nil {
			return err
		}
		// mint interest earned in cred and distribute
		period := ctx.BlockHeader().Time.Sub(vote.Timestamp().CreatedTime)
		interest := k.interest(ctx, vote.Amount(), categoryID, period)

		_, err = k.truBankKeeper.MintAndAddCoin(
			ctx, vote.Creator(), categoryID, interest)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf(
			"Distributed %s with interest %s for period %s to %s",
			vote.Amount().String(),
			interest.String(),
			period.String(),
			vote.Creator().String()))
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

// TODO: pass in category keeper so we can get total cred
func (k Keeper) interest(
	ctx sdk.Context,
	amount sdk.Coin,
	categoryID int64,
	period time.Duration) sdk.Int {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22

	totalSupply := sdk.NewDec(1000000000000000)

	// inputs
	maxAmount := totalSupply
	amountWeight := k.GetParams(ctx).AmountWeight
	periodWeight := k.GetParams(ctx).PeriodWeight
	maxInterestRate := k.GetParams(ctx).MaxInterestRate
	maxPeriod := k.storyKeeper.GetParams(ctx).ExpireDuration
	spew.Dump(maxPeriod)

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

	return interest.RoundInt()
}
