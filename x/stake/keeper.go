package stake

import (
	"fmt"
	"time"

	"github.com/TruStory/truchain/x/story"

	app "github.com/TruStory/truchain/types"
	trubank "github.com/TruStory/truchain/x/trubank"
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

// RedistributeStake distributes stake according to story resolution
//
// Example:
// true: 100, 100, 100
// false: 100
// true wins
// total pool: 400
// win pool: 300
// lose pool: 100
// winner eq: total pool * (stake amount / total win pool)
// winners: 400 * (100/300) = 133.33
func (k Keeper) RedistributeStake(ctx sdk.Context, votes []Voter) (app.StakeDistributionResults, sdk.Error) {
	truePool := sdk.ZeroInt()
	falsePool := sdk.ZeroInt()
	rewardPool := sdk.ZeroInt()
	for _, v := range votes {
		voteStake := v.Amount().Amount
		rewardPool = rewardPool.Add(voteStake)
		if v.VoteChoice() == true {
			truePool = truePool.Add(voteStake)
		} else {
			falsePool = falsePool.Add(voteStake)
		}
	}

	winPool := falsePool
	totalPool := truePool.Add(falsePool)
	result := app.StakeDistributionResults{
		TotalAmount: sdk.NewCoin(app.StakeDenom, totalPool),
		Type:        app.DistributionMajorityNotReached,
	}
	rewards := make([]app.StakeReward, 0)
	truePoolDec := sdk.NewDecFromInt(truePool)
	truePoolPercentOfTotalPool := truePoolDec.QuoInt(totalPool)

	falsePoolDec := sdk.NewDecFromInt(falsePool)
	falsePoolPercentOfTotalPool := falsePoolDec.QuoInt(totalPool)

	if truePoolPercentOfTotalPool.GTE(k.majorityPercent(ctx)) {
		result.Type = app.DistributionBackersWin
		// true pool >= 51% total pool
		winPool = truePool
		for _, v := range votes {
			if v.VoteChoice() == true {
				r, err := k.rewardStaker(ctx, v, winPool, rewardPool)
				if err != nil {
					return result, err
				}
				rewards = append(rewards, app.StakeReward{Account: v.Creator(), Amount: *r})
			}
		}
	} else if falsePoolPercentOfTotalPool.GTE(k.majorityPercent(ctx)) {
		result.Type = app.DistributionChallengersWin
		// false pool >= 51% total pool
		for _, v := range votes {
			if v.VoteChoice() == false {
				r, err := k.rewardStaker(ctx, v, winPool, rewardPool)
				if err != nil {
					return result, err
				}
				rewards = append(rewards, app.StakeReward{Account: v.Creator(), Amount: *r})
			}
		}
	} else {
		// 51% majority not met, return stake
		for _, v := range votes {
			transactionType := trubank.BackingReturned
			if v.VoteChoice() == false {
				transactionType = trubank.ChallengeReturned
			}
			_, err := k.truBankKeeper.AddCoin(ctx, v.Creator(), v.Amount(), v.StoryID(), transactionType, 0)
			if err != nil {
				return result, err
			}
		}
	}
	result.Rewards = rewards
	return result, nil
}

// DistributeInterest distributes interest for staking
func (k Keeper) DistributeInterest(ctx sdk.Context, votes []Voter) (app.InterestDistributionResults, sdk.Error) {
	logger := ctx.Logger().With("module", StoreKey)

	result := app.InterestDistributionResults{}
	total := sdk.NewCoin(app.StakeDenom, sdk.NewInt(0))
	interests := make([]app.Interest, 0)
	for _, v := range votes {
		period := ctx.BlockHeader().Time.Sub(v.Timestamp().CreatedTime)
		interest := k.interest(ctx, v.Amount(), period)
		interestCoin := sdk.NewCoin(app.StakeDenom, interest)
		total.Plus(interestCoin)
		_, err := k.truBankKeeper.AddCoin(ctx, v.Creator(), interestCoin, v.StoryID(), trubank.Interest, 0)
		if err != nil {
			return result, err
		}
		interests = append(interests, app.Interest{Account: v.Creator(), Amount: interestCoin, Rate: interest})
		logger.Info(fmt.Sprintf("Distributed interest %s to %s", interestCoin, v.Creator()))
	}
	result.TotalAmount = total
	result.Interests = interests
	return result, nil
}

// ValidateAmount validates the stake amount
func (k Keeper) ValidateAmount(ctx sdk.Context, amount sdk.Coin) sdk.Error {
	maxAmount := k.GetParams(ctx).MaxAmount
	if maxAmount.IsLT(amount) {
		return ErrOverMaxAmount()
	}

	return nil
}

// ValidateStoryState makes sure only a pending story can be staked
func (k Keeper) ValidateStoryState(ctx sdk.Context, storyID int64) sdk.Error {
	s, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	if s.Status != story.Pending {
		return ErrInvalidStoryState(s.Status.String())
	}

	return nil
}

// Interest calculates interest for staked amount
func (k Keeper) interest(ctx sdk.Context, amount sdk.Coin, period time.Duration) sdk.Int {
	interestRate := k.GetParams(ctx).InterestRate

	periodDec := sdk.NewDec(period.Nanoseconds())
	amountDec := sdk.NewDecFromInt(amount.Amount)

	oneYear := time.Hour * 24 * 365
	oneYearDec := sdk.NewDec(oneYear.Nanoseconds())
	interest := interestRate.Mul(periodDec.Quo(oneYearDec)).Mul(amountDec)

	return interest.RoundInt()
}

// distribute stake proportionally to winner
func (k Keeper) rewardStaker(ctx sdk.Context, staker Voter, winPool sdk.Int, rewardPool sdk.Int) (*sdk.Coin, sdk.Error) {
	logger := ctx.Logger().With("module", StoreKey)
	rewardAmount := rewardAmount(staker.Amount().Amount, winPool, rewardPool)
	rewardCoin := sdk.NewCoin(app.StakeDenom, rewardAmount)
	transactionType := trubank.BackingReturned
	if staker.VoteChoice() == false {
		transactionType = trubank.ChallengeReturned
	}
	_, err := k.truBankKeeper.AddCoin(ctx, staker.Creator(), staker.Amount(), staker.StoryID(), transactionType, 0)
	if err != nil {
		return nil, err
	}
	_, err = k.truBankKeeper.AddCoin(ctx, staker.Creator(), rewardCoin, staker.StoryID(), trubank.RewardPool, 0)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Distributed stake reward %s to %s", rewardCoin, staker.Creator()))
	return &rewardCoin, nil
}

// reward stake = reward pool * (stake amount / winner pool)
func rewardAmount(stakeAmount sdk.Int, winPool sdk.Int, rewardPool sdk.Int) sdk.Int {
	totalAmount := sdk.NewDecFromInt(rewardPool).
		MulInt(stakeAmount).
		QuoInt(winPool).
		TruncateInt()
	return totalAmount.Sub(stakeAmount)
}
