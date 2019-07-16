package params

import (
	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/slashing"
	"github.com/TruStory/truchain/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	accountKeeper   account.Keeper
	communityKeeper community.Keeper
	claimKeeper     claim.Keeper
	bankKeeper      bank.Keeper
	stakingKeeper   staking.Keeper
	slashingKeeper  slashing.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	accountKeeper account.Keeper,
	communityKeeper community.Keeper,
	claimKeeper claim.Keeper,
	bankKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
	slashingKeeper slashing.Keeper) Keeper {

	return Keeper{
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
		AccountParams:  account.Params{},
		BankParams:     bank.Params{},
		SlashingParams: slashing.Params{},
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
	}
}
