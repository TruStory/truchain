package staking

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/bank"

	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/TruStory/truchain/x/claim"
)

// Keeper is the model object for the package staking module
type Keeper struct {
	storeKey    sdk.StoreKey
	codec       *codec.Codec
	paramStore  params.Subspace
	codespace   sdk.CodespaceType
	bankKeeper  bank.Keeper
	authKeeper  AuthKeeper
	claimKeeper ClaimKeeper
}

// NewKeeper creates a staking keeper.
func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey,
	authKeeper AuthKeeper, bankKeeper bank.Keeper, claimKeeper ClaimKeeper,
	paramStore params.Subspace,
	codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:    storeKey,
		codec:       codec,
		paramStore:  paramStore.WithKeyTable(ParamKeyTable()),
		codespace:   codespace,
		bankKeeper:  bankKeeper,
		authKeeper:  authKeeper,
		claimKeeper: claimKeeper,
	}
}

func (k Keeper) Arguments(ctx sdk.Context) []Argument {
	arguments := make([]Argument, 0)
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), ArgumentsKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var argument Argument
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &argument)
		arguments = append(arguments, argument)
	}
	return arguments
}

func (k Keeper) Stakes(ctx sdk.Context) []Stake {
	stakes := make([]Stake, 0)
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), StakesKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stake Stake
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stake)
		stakes = append(stakes, stake)
	}
	return stakes
}

func (k Keeper) ClaimArguments(ctx sdk.Context, claimID uint64) []Argument {
	arguments := make([]Argument, 0)
	k.IterateClaimArguments(ctx, claimID, func(argument Argument) bool {
		arguments = append(arguments, argument)
		return false
	})
	return arguments
}

func (k Keeper) ArgumentStakes(ctx sdk.Context, argumentID uint64) []Stake {
	stakes := make([]Stake, 0)
	k.IterateArgumentStakes(ctx, argumentID, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) UserStakes(ctx sdk.Context, address sdk.AccAddress) []Stake {
	stakes := make([]Stake, 0)
	k.IterateUserStakes(ctx, address, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) UserArguments(ctx sdk.Context, address sdk.AccAddress) []Argument {
	arguments := make([]Argument, 0)
	k.IterateUserArguments(ctx, address, func(argument Argument) bool {
		arguments = append(arguments, argument)
		return false
	})
	return arguments
}

func (k Keeper) SubmitUpvote(ctx sdk.Context, argumentID uint64, creator sdk.AccAddress) (Stake, sdk.Error) {
	argument, ok := k.getArgument(ctx, argumentID)
	if !ok {
		return Stake{}, ErrCodeUnknownArgument(argumentID)
	}
	stakes := k.ArgumentStakes(ctx, argumentID)
	for _, s := range stakes {
		if s.Creator.Equals(creator) {
			return Stake{}, ErrCodeDuplicateStake(argumentID)
		}
	}
	upvoteStake := k.GetParams(ctx).UpvoteStake
	stake, err := k.newStake(ctx, upvoteStake, creator, StakeUpvote, argumentID)
	if err != nil {
		return stake, err
	}
	argument.UpvotedCount = argument.UpvotedCount + 1
	argument.UpvotedStake = argument.UpvotedStake.Add(stake.Amount)
	argument.TotalStake = argument.TotalStake.Add(stake.Amount)
	argument.UpdatedTime = ctx.BlockHeader().Time
	k.setArgument(ctx, argument)
	return stake, nil
}

func (k Keeper) SubmitArgument(ctx sdk.Context, body, summary string,
	creator sdk.AccAddress, claimID uint64, stakeType StakeType) (Argument, sdk.Error) {
	// only backing or challenge
	if !stakeType.ValidForArgument() {
		return Argument{}, ErrCodeInvalidStakeType(stakeType)
	}
	//  reject jailed accounts
	jailed, err := k.authKeeper.IsJailed(ctx, creator)
	if err != nil {
		return Argument{}, err
	}
	if jailed {
		return Argument{}, ErrCodeAccountJailed(creator)
	}
	_, ok := k.claimKeeper.Claim(ctx, claimID)
	if !ok {
		return Argument{}, claim.ErrUnknownClaim(claimID)
	}

	arguments := k.ClaimArguments(ctx, claimID)
	count := 0
	for _, a := range arguments {
		if a.Creator.Equals(creator) {
			count++
		}
	}
	p := k.GetParams(ctx)
	if count >= p.MaxArgumentsPerClaim {
		return Argument{}, ErrCodeMaxNumOfArgumentsReached(p.MaxArgumentsPerClaim)
	}

	creationAmount := p.ArgumentCreationStake
	argumentID, err := k.argumentID(ctx)
	if err != nil {
		return Argument{}, err
	}
	argument := Argument{
		ID:           argumentID,
		Creator:      creator,
		ClaimID:      claimID,
		Summary:      summary,
		Body:         body,
		StakeType:    stakeType,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		UpvotedStake: sdk.NewInt64Coin(app.StakeDenom, 0),
		TotalStake:   creationAmount,
	}
	_, err = k.newStake(ctx, creationAmount, creator, stakeType, argument.ID)
	if err != nil {
		return Argument{}, err
	}

	k.setArgument(ctx, argument)
	k.setArgumentID(ctx, argumentID+1)
	k.setClaimArgument(ctx, claimID, argument.ID)
	k.serUserArgument(ctx, creator, argument.ID)
	return argument, nil
}

func (k Keeper) getArgument(ctx sdk.Context, argumentID uint64) (Argument, bool) {
	argument := Argument{}
	bz := k.store(ctx).Get(ArgumentKey(argumentID))
	if bz == nil {
		return Argument{}, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &argument)
	return argument, true
}

func (k Keeper) setArgument(ctx sdk.Context, argument Argument) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(argument)
	k.store(ctx).Set(ArgumentKey(argument.ID), bz)
}

func (k Keeper) checkStakeThreshold(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	balance := k.bankKeeper.GetCoins(ctx, address).AmountOf(app.StakeDenom)
	if balance.IsZero() {
		return sdk.ErrInsufficientFunds("Insufficient coins")
	}
	p := k.GetParams(ctx)
	period := p.StakeLimitDays

	staked := sdk.NewInt(0)
	fromDate := ctx.BlockHeader().Time.Add(time.Duration(-1) * period)
	k.IterateAfterCreatedTimeUserStakes(ctx, address,
		fromDate, func(stake Stake) bool {
			// only account for non expired since expired would already have refunded the stake
			if stake.Expired {
				return false
			}
			staked = staked.Add(stake.Amount.Amount)
			return false
		},
	)

	total := sdk.NewDecFromInt(balance.Add(staked))
	stakedDec := sdk.NewDecFromInt(staked)
	if stakedDec.Quo(total).GTE(p.StakeLimitPercent) {
		return ErrCodeMaxAmountStakingReached(int(ctx.BlockHeader().Time.Sub(fromDate).Hours()))
	}
	return nil
}

func (k Keeper) newStake(ctx sdk.Context, amount sdk.Coin, creator sdk.AccAddress,
	stakeType StakeType, argumentID uint64) (Stake, sdk.Error) {
	if !stakeType.Valid() {
		return Stake{}, ErrCodeInvalidStakeType(stakeType)
	}
	err := k.checkStakeThreshold(ctx, creator)
	if err != nil {
		return Stake{}, err
	}
	period := k.GetParams(ctx).Period
	stakeID, err := k.stakeID(ctx)
	if err != nil {
		return Stake{}, err
	}
	_, err = k.bankKeeper.SubtractCoin(ctx, creator, amount, argumentID, stakeType.BankTransactionType())
	if err != nil {
		return Stake{}, err
	}
	stake := Stake{
		ID:          stakeID,
		ArgumentID:  argumentID,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(period),
		Creator:     creator,
		Amount:      amount,
		Type:        stakeType,
	}
	k.setStake(ctx, stake)
	k.setStakeID(ctx, stakeID+1)
	k.InsertActiveStakeQueue(ctx, stakeID, stake.EndTime)
	k.setArgumentStake(ctx, argumentID, stake.ID)
	k.setUserStake(ctx, creator, stake.CreatedTime, stake.ID)
	return stake, nil
}

func (k Keeper) getStake(ctx sdk.Context, stakeID uint64) Stake {
	stake := Stake{}
	k.codec.MustUnmarshalBinaryLengthPrefixed(k.store(ctx).Get(StakeKey(stakeID)), &stake)
	return stake
}
func (k Keeper) setStake(ctx sdk.Context, stake Stake) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stake)
	k.store(ctx).Set(StakeKey(stake.ID), bz)
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) setStakeID(ctx sdk.Context, stakeID uint64) {
	k.setID(ctx, StakeIDKey, stakeID)
}

func (k Keeper) setArgumentID(ctx sdk.Context, argumentID uint64) {
	k.setID(ctx, ArgumentIDKey, argumentID)
}

func (k Keeper) setID(ctx sdk.Context, key []byte, length uint64) {
	b := k.codec.MustMarshalBinaryBare(length)
	k.store(ctx).Set(key, b)
}

func (k Keeper) stakeID(ctx sdk.Context) (uint64, sdk.Error) {
	id, err := k.getID(ctx, StakeIDKey)
	if err != nil {
		return 0, ErrCodeUnknownStake(id)
	}
	return id, nil
}

func (k Keeper) argumentID(ctx sdk.Context) (uint64, sdk.Error) {
	id, err := k.getID(ctx, ArgumentIDKey)
	if err != nil {
		return 0, ErrCodeUnknownArgument(id)
	}
	return id, nil
}

func (k Keeper) getID(ctx sdk.Context, key []byte) (uint64, sdk.Error) {
	var id uint64
	b := k.store(ctx).Get(key)
	if b == nil {
		return 0, sdk.ErrInternal("unknown id")
	}
	k.codec.MustUnmarshalBinaryBare(b, &id)
	return id, nil
}

// InsertActiveStakeQueue inserts a stakeID into the active stake queue at endTime
func (k Keeper) InsertActiveStakeQueue(ctx sdk.Context, stakeID uint64, endTime time.Time) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stakeID)
	k.store(ctx).Set(ActiveStakeQueueKey(stakeID, endTime), bz)
}

// RemoveFromActiveStakeQueue removes a stakeID from the Active Stake Queue
func (k Keeper) RemoveFromActiveStakeQueue(ctx sdk.Context, stakeID uint64, endTime time.Time) {
	k.store(ctx).Delete(ActiveStakeQueueKey(stakeID, endTime))
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ModuleName)
}
