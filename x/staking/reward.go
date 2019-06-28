package staking

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/bank"

	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) splitReward(ctx sdk.Context, interest sdk.Dec) (creator, staker sdk.Int) {
	p := k.GetParams(ctx)
	creatorShare := interest.Mul(p.CreatorShare)
	stakerShare := interest.Sub(creatorShare)
	return creatorShare.RoundInt(), stakerShare.RoundInt()
}

type RewardResultType byte

const (
	RewardResultArgumentCreation RewardResultType = iota
	RewardResultUpvoteSplit
)

type RewardResult struct {
	Type                  RewardResultType `json:"type"`
	ArgumentCreator       sdk.AccAddress   `json:"argument_creator"`
	ArgumentCreatorReward sdk.Coin         `json:"argument_creator_reward"`
	StakeCreator          sdk.AccAddress   `json:"stake_creator"`
	StakeCreatorReward    sdk.Coin         `json:"stake_creator_reward"`
}

func (k Keeper) distributeReward(ctx sdk.Context, stake Stake) (RewardResult, sdk.Error) {
	argument, ok := k.getArgument(ctx, stake.ArgumentID)
	if !ok {
		return RewardResult{}, ErrCodeUnknownArgument(stake.ArgumentID)
	}
	interest := k.interest(ctx, stake.Amount, ctx.BlockHeader().Time.Sub(stake.CreatedTime))
	// creator receives 100% interest of his own stake
	if !argument.Creator.Equals(stake.Creator) {
		reward := sdk.NewCoin(app.StakeDenom, interest.RoundInt())
		_, err := k.bankKeeper.AddCoin(ctx,
			argument.Creator,
			reward,
			argument.ID,
			bank.TransactionInterestArgumentCreation)
		if err != nil {
			return RewardResult{}, err
		}
		return RewardResult{Type: RewardResultArgumentCreation,
			ArgumentCreator:       argument.Creator,
			ArgumentCreatorReward: reward}, nil
	}
	creatorReward, stakerReward := k.splitReward(ctx, interest)
	creatorRewardCoin := sdk.NewCoin(app.StakeDenom, creatorReward)
	stakerRewardCoin := sdk.NewCoin(app.StakeDenom, stakerReward)
	_, err := k.bankKeeper.AddCoin(ctx,
		argument.Creator,
		creatorRewardCoin,
		argument.ID,
		bank.TransactionInterestUpvoteReceived)
	if err != nil {
		return RewardResult{}, err
	}
	_, err = k.bankKeeper.AddCoin(ctx,
		stake.Creator,
		stakerRewardCoin,
		argument.ID,
		bank.TransactionInterestUpvoteGiven)
	if err != nil {
		return RewardResult{}, err
	}
	rewardResult := RewardResult{
		Type:                  RewardResultUpvoteSplit,
		ArgumentCreator:       argument.Creator,
		ArgumentCreatorReward: creatorRewardCoin,
		StakeCreator:          stake.Creator,
		StakeCreatorReward:    stakerRewardCoin,
	}
	return rewardResult, nil
}

func (k Keeper) interest(ctx sdk.Context, amount sdk.Coin, period time.Duration) sdk.Dec {
	// TODO: https://github.com/TruStory/truchain/issues/677
	// use interest from distribution module
	interestRate := k.GetParams(ctx).InterestRate
	periodDec := sdk.NewDec(period.Nanoseconds())
	amountDec := sdk.NewDecFromInt(amount.Amount)
	oneYear := time.Hour * 24 * 365
	oneYearDec := sdk.NewDec(oneYear.Nanoseconds())
	interest := interestRate.Mul(periodDec.Quo(oneYearDec)).Mul(amountDec)
	return interest
}
