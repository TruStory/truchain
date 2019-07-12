package params

import (
	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/slashing"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/staking"
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
	accountKeeper    account.Keeper
	communityKeeper  community.Keeper
	claimKeeper      claim.Keeper
	bankKeeper       bank.Keeper
	stakingKeeper    staking.Keeper
	slashingKeeper   slashing.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	argumentKeeper argument.Keeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	expirationKeeper expiration.Keeper,
	stakeKeeper stake.Keeper,
	storyKeeper story.WriteKeeper,
	accountKeeper account.Keeper,
	communityKeeper community.Keeper,
	claimKeeper claim.Keeper,
	bankKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
	slashingKeeper slashing.Keeper) Keeper {

	return Keeper{
		argumentKeeper,
		backingKeeper,
		challengeKeeper,
		expirationKeeper,
		stakeKeeper,
		storyKeeper,
		accountKeeper,
		communityKeeper,
		claimKeeper,
		bankKeeper,
		stakingKeeper,
		slashingKeeper,
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
			MaxAmount:        k.stakeKeeper.GetParams(ctx).MaxAmount,
			InterestRate:     k.stakeKeeper.GetParams(ctx).InterestRate,
			MajorityPercent:  k.stakeKeeper.GetParams(ctx).MajorityPercent,
			StakeToCredRatio: k.stakeKeeper.GetParams(ctx).StakeToCredRatio,
		},
		StoryParams: story.Params{
			ExpireDuration: k.storyKeeper.GetParams(ctx).ExpireDuration,
			MinStoryLength: k.storyKeeper.GetParams(ctx).MinStoryLength,
			MaxStoryLength: k.storyKeeper.GetParams(ctx).MaxStoryLength,
		},
		AccountParams: account.Params{
			Registrar:     k.accountKeeper.GetParams(ctx).Registrar,
			MaxSlashCount: k.accountKeeper.GetParams(ctx).MaxSlashCount,
			JailDuration:  k.accountKeeper.GetParams(ctx).JailDuration,
		},
		CommunityParams: community.Params{
			MinIDLength:          k.communityKeeper.GetParams(ctx).MinIDLength,
			MaxIDLength:          k.communityKeeper.GetParams(ctx).MaxIDLength,
			MinNameLength:        k.communityKeeper.GetParams(ctx).MinNameLength,
			MaxNameLength:        k.communityKeeper.GetParams(ctx).MaxNameLength,
			MaxDescriptionLength: k.communityKeeper.GetParams(ctx).MaxDescriptionLength,
		},
		ClaimParams: claim.Params{
			MinClaimLength: k.claimKeeper.GetParams(ctx).MinClaimLength,
			MaxClaimLength: k.claimKeeper.GetParams(ctx).MaxClaimLength,
		},
		BankParams: bank.Params{
			RewardBrokerAddress: k.bankKeeper.GetParams(ctx).RewardBrokerAddress,
		},
		StakingParams: staking.Params{
			Period:                   k.stakingKeeper.GetParams(ctx).Period,
			ArgumentCreationStake:    k.stakingKeeper.GetParams(ctx).ArgumentCreationStake,
			ArgumentBodyMaxLength:    k.stakingKeeper.GetParams(ctx).ArgumentBodyMaxLength,
			ArgumentBodyMinLength:    k.stakingKeeper.GetParams(ctx).ArgumentBodyMinLength,
			ArgumentSummaryMaxLength: k.stakingKeeper.GetParams(ctx).ArgumentSummaryMaxLength,
			ArgumentSummaryMinLength: k.stakingKeeper.GetParams(ctx).ArgumentSummaryMinLength,
			UpvoteStake:              k.stakingKeeper.GetParams(ctx).UpvoteStake,
			CreatorShare:             k.stakingKeeper.GetParams(ctx).CreatorShare,
			InterestRate:             k.stakingKeeper.GetParams(ctx).InterestRate,
			StakeLimitPercent:        k.stakingKeeper.GetParams(ctx).StakeLimitPercent,
			StakeLimitDays:           k.stakingKeeper.GetParams(ctx).StakeLimitDays,
			UnjailUpvotes:            k.stakingKeeper.GetParams(ctx).UnjailUpvotes,
			MaxArgumentsPerClaim:     k.stakingKeeper.GetParams(ctx).MaxArgumentsPerClaim,
		},
		SlashingParams: slashing.Params{
			MinSlashCount:  k.slashingKeeper.GetParams(ctx).MinSlashCount,
			SlashMagnitude: k.slashingKeeper.GetParams(ctx).SlashMagnitude,
			SlashMinStake:  k.slashingKeeper.GetParams(ctx).SlashMinStake,
			SlashAdmins:    k.slashingKeeper.GetParams(ctx).SlashAdmins,
			CuratorShare:   k.slashingKeeper.GetParams(ctx).CuratorShare,
		},
	}
}
