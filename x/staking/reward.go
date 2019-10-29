package staking

import (
	"time"

	app "github.com/TruStory/truchain/types"

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
	argument, ok := k.Argument(ctx, stake.ArgumentID)
	if !ok {
		return RewardResult{}, ErrCodeUnknownArgument(stake.ArgumentID)
	}
	claim, ok := k.claimKeeper.Claim(ctx, argument.ClaimID)
	if !ok {
		return RewardResult{}, ErrCodeUnknownClaim(claim.ID)
	}

	// refund
	var refundType TransactionType

	switch stake.Type {
	case StakeBacking:
		refundType = TransactionBackingReturned
	case StakeChallenge:
		refundType = TransactionChallengeReturned
	case StakeUpvote:
		refundType = TransactionUpvoteReturned
	default:
		return RewardResult{}, ErrCodeUnknownStakeType()
	}

	_, err := k.bankKeeper.AddCoin(ctx, stake.Creator, stake.Amount, stake.ArgumentID,
		refundType, WithCommunityID(argument.CommunityID),
		FromModuleAccount(UserStakesPoolName),
	)
	if err != nil {
		return RewardResult{}, err
	}

	if err != nil {
		return RewardResult{}, err
	}

	interest := k.interest(ctx, stake.Amount, stake.EndTime.Sub(stake.CreatedTime))
	// creator receives 100% interest of his own stake
	if argument.Creator.Equals(stake.Creator) {
		reward := sdk.NewCoin(app.StakeDenom, interest.RoundInt())
		_, err := k.bankKeeper.AddCoin(ctx,
			argument.Creator,
			reward,
			argument.ID,
			TransactionInterestArgumentCreation,
			WithCommunityID(argument.CommunityID),
			FromModuleAccount(UserRewardPoolName),
		)
		if err != nil {
			return RewardResult{}, err
		}
		k.addEarnedCoin(ctx, argument.Creator, claim.CommunityID, reward.Amount)
		return RewardResult{Type: RewardResultArgumentCreation,
			ArgumentCreator:       argument.Creator,
			ArgumentCreatorReward: reward}, nil
	}
	creatorReward, stakerReward := k.splitReward(ctx, interest)
	creatorRewardCoin := sdk.NewCoin(app.StakeDenom, creatorReward)
	stakerRewardCoin := sdk.NewCoin(app.StakeDenom, stakerReward)
	_, err = k.bankKeeper.AddCoin(ctx,
		argument.Creator,
		creatorRewardCoin,
		stake.ID,
		TransactionInterestUpvoteReceived,
		WithCommunityID(argument.CommunityID),
		FromModuleAccount(UserRewardPoolName),
	)
	if err != nil {
		return RewardResult{}, err
	}

	_, err = k.bankKeeper.AddCoin(ctx,
		stake.Creator,
		stakerRewardCoin,
		argument.ID,
		TransactionInterestUpvoteGiven,
		WithCommunityID(argument.CommunityID),
		FromModuleAccount(UserRewardPoolName),
	)
	if err != nil {
		return RewardResult{}, err
	}

	k.addEarnedCoin(ctx, argument.Creator, claim.CommunityID, creatorRewardCoin.Amount)
	k.addEarnedCoin(ctx, stake.Creator, claim.CommunityID, stakerRewardCoin.Amount)
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
	interestRate := k.GetParams(ctx).InterestRate
	return Interest(interestRate, amount, period)
}

// Interest takes an annual inflation/interest rate and calculates the return on an amount staked for a given period
func Interest(interestRate sdk.Dec, amount sdk.Coin, period time.Duration) sdk.Dec {
	periodDec := sdk.NewDec(period.Nanoseconds())
	amountDec := sdk.NewDecFromInt(amount.Amount)
	oneYear := time.Hour * 24 * 365
	oneYearDec := sdk.NewDec(oneYear.Nanoseconds())
	interest := interestRate.Mul(periodDec.Quo(oneYearDec)).Mul(amountDec)
	return interest
}
