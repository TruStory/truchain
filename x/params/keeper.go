package params

import (
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// StoreKey is string representation of the store key
	StoreKey = "params"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	argumentKeeper   argument.Keeper
	backingKeeper    backing.ReadKeeper
	challengeKeeper  challenge.ReadKeeper
	expirationKeeper expiration.Keeper
	stakeKeeper      stake.Keeper
	storyKeeper      story.ReadKeeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	argumentKeeper argument.Keeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	expirationKeeper expiration.Keeper,
	stakeKeeper stake.Keeper,
	storyKeeper story.WriteKeeper) Keeper {

	return Keeper{
		argumentKeeper,
		backingKeeper,
		challengeKeeper,
		expirationKeeper,
		stakeKeeper,
		storyKeeper,
	}
}

// Params returns all parameters for clients
func (k Keeper) Params(ctx sdk.Context) Params {
	return Params{
		ArgumentParams: argument.Params{
			MinArgumentLength: k.argumentKeeper.GetParams(ctx).MinArgumentLength,
			MaxArgumentLength: k.argumentKeeper.GetParams(ctx).MaxArgumentLength,
		},
		ChallengeParams: challenge.Params{
			MinChallengeStake: k.challengeKeeper.GetParams(ctx).MinChallengeStake,
		},
		StakeParams: stake.Params{
			MaxAmount:       k.stakeKeeper.GetParams(ctx).MaxAmount,
			MinInterestRate: k.stakeKeeper.GetParams(ctx).MinInterestRate,
			MaxInterestRate: k.stakeKeeper.GetParams(ctx).MaxInterestRate,
			MajorityPercent: k.stakeKeeper.GetParams(ctx).MajorityPercent,
			AmountWeight:    k.stakeKeeper.GetParams(ctx).AmountWeight,
			PeriodWeight:    k.stakeKeeper.GetParams(ctx).PeriodWeight,
		},
		StoryParams: story.Params{
			ExpireDuration: k.storyKeeper.GetParams(ctx).ExpireDuration,
			MinStoryLength: k.storyKeeper.GetParams(ctx).MinStoryLength,
			MaxStoryLength: k.storyKeeper.GetParams(ctx).MaxStoryLength,
		},
	}
}
