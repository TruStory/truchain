package params

import (
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/voting"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// StoreKey is string representation of the store key
	StoreKey = "params"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	backingKeeper    backing.ReadKeeper
	challengeKeeper  challenge.ReadKeeper
	expirationKeeper expiration.Keeper
	stakeKeeper      stake.Keeper
	storyKeeper      story.ReadKeeper
	votingKeeper     voting.ReadKeeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	expirationKeeper expiration.Keeper,
	stakeKeeper stake.Keeper,
	storyKeeper story.WriteKeeper,
	votingKeeper voting.ReadKeeper) Keeper {

	return Keeper{
		backingKeeper,
		challengeKeeper,
		expirationKeeper,
		stakeKeeper,
		storyKeeper,
		votingKeeper,
	}
}

// Params returns all parameters for clients
func (k Keeper) Params(ctx sdk.Context) Params {
	return Params{
		ChallengeParams: challenge.Params{
			ChallengeToBackingRatio: k.challengeKeeper.GetParams(ctx).ChallengeToBackingRatio,
			MinChallengeThreshold:   k.challengeKeeper.GetParams(ctx).MinChallengeThreshold,
			MinChallengeStake:       k.challengeKeeper.GetParams(ctx).MinChallengeStake,
		},
		StakeParams: stake.Params{
			MinArgumentLength: k.stakeKeeper.GetParams(ctx).MinArgumentLength,
			MaxArgumentLength: k.stakeKeeper.GetParams(ctx).MaxArgumentLength,
			MinInterestRate:   k.stakeKeeper.GetParams(ctx).MinInterestRate,
			MaxInterestRate:   k.stakeKeeper.GetParams(ctx).MaxInterestRate,
			AmountWeight:      k.stakeKeeper.GetParams(ctx).AmountWeight,
			PeriodWeight:      k.stakeKeeper.GetParams(ctx).PeriodWeight,
		},
		StoryParams: story.Params{
			ExpireDuration: k.storyKeeper.GetParams(ctx).ExpireDuration,
			MinStoryLength: k.storyKeeper.GetParams(ctx).MinStoryLength,
			MaxStoryLength: k.storyKeeper.GetParams(ctx).MaxStoryLength,
			VotingDuration: k.storyKeeper.GetParams(ctx).VotingDuration,
		},
		VotingParams: voting.Params{
			StakerRewardPoolShare: k.votingKeeper.GetParams(ctx).StakerRewardPoolShare,
			MajorityPercent:       k.votingKeeper.GetParams(ctx).MajorityPercent,
		},
	}
}
