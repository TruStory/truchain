package staking

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/supply"

	app "github.com/TruStory/truchain/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the package staking module
type Keeper struct {
	storeKey      sdk.StoreKey
	codec         *codec.Codec
	paramStore    params.Subspace
	codespace     sdk.CodespaceType
	bankKeeper    BankKeeper
	accountKeeper AccountKeeper
	claimKeeper   ClaimKeeper
	supplyKeeper  supply.Keeper
}

// NewKeeper creates a staking keeper.
func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey,
	accountKeeper AccountKeeper, bankKeeper BankKeeper, claimKeeper ClaimKeeper, supplyKeeper supply.Keeper,
	paramStore params.Subspace,
	codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:      storeKey,
		codec:         codec,
		paramStore:    paramStore.WithKeyTable(ParamKeyTable()),
		codespace:     codespace,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		claimKeeper:   claimKeeper,
		supplyKeeper:  supplyKeeper,
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

func (k Keeper) CommunityStakes(ctx sdk.Context, communityID string) []Stake {
	stakes := make([]Stake, 0)
	k.IterateCommunityStakes(ctx, communityID, func(stake Stake) bool {
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

func (k Keeper) UserCommunityStakes(ctx sdk.Context, address sdk.AccAddress, communityID string) []Stake {
	stakes := make([]Stake, 0)
	k.IterateUserCommunityStakes(ctx, address, communityID, func(stake Stake) bool {
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
	err := k.checkJailed(ctx, creator)
	if err != nil {
		return Stake{}, err
	}
	argument, ok := k.Argument(ctx, argumentID)
	if !ok {
		return Stake{}, ErrCodeUnknownArgument(argumentID)
	}
	stakes := k.ArgumentStakes(ctx, argumentID)
	for _, s := range stakes {
		if s.Creator.Equals(creator) {
			return Stake{}, ErrCodeDuplicateStake(argumentID)
		}
	}
	claim, ok := k.claimKeeper.Claim(ctx, argument.ClaimID)
	if !ok {
		return Stake{}, ErrCodeUnknownClaim(argument.ClaimID)
	}

	upvoteStake := k.GetParams(ctx).UpvoteStake
	stake, err := k.newStake(ctx, upvoteStake, creator, StakeUpvote, argumentID, claim.CommunityID)
	if err != nil {
		return stake, err
	}
	argument.UpvotedCount = argument.UpvotedCount + 1
	argument.UpvotedStake = argument.UpvotedStake.Add(stake.Amount)
	argument.TotalStake = argument.TotalStake.Add(stake.Amount)
	argument.UpdatedTime = ctx.BlockHeader().Time
	k.setArgument(ctx, argument)

	switch {
	case argument.StakeType == StakeBacking:
		err := k.claimKeeper.AddBackingStake(ctx, argument.ClaimID, stake.Amount)
		if err != nil {
			return Stake{}, err
		}
	case argument.StakeType == StakeChallenge:
		err := k.claimKeeper.AddChallengeStake(ctx, argument.ClaimID, stake.Amount)
		if err != nil {
			return Stake{}, err
		}
	}

	return stake, nil
}

func (k Keeper) checkJailed(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	jailed, err := k.accountKeeper.IsJailed(ctx, address)
	if err != nil {
		return err
	}
	if jailed {
		return ErrCodeAccountJailed(address)
	}
	return nil
}

func (k Keeper) SubmitArgument(ctx sdk.Context, body, summary string,
	creator sdk.AccAddress, claimID uint64, stakeType StakeType) (Argument, sdk.Error) {
	// only backing or challenge
	if !stakeType.ValidForArgument() {
		return Argument{}, ErrCodeInvalidStakeType(stakeType)
	}
	err := k.checkJailed(ctx, creator)
	if err != nil {
		return Argument{}, err
	}
	claim, ok := k.claimKeeper.Claim(ctx, claimID)
	if !ok {
		return Argument{}, ErrCodeUnknownClaim(claimID)
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
		CommunityID:  claim.CommunityID,
		Summary:      summary,
		Body:         body,
		StakeType:    stakeType,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		UpvotedStake: sdk.NewInt64Coin(app.StakeDenom, 0),
		TotalStake:   creationAmount,
		EditedTime:   ctx.BlockHeader().Time,
		Edited:       false,
	}
	_, err = k.newStake(ctx, creationAmount, creator, stakeType, argument.ID, claim.CommunityID)
	if err != nil {
		return Argument{}, err
	}

	k.setArgument(ctx, argument)
	k.setArgumentID(ctx, argumentID+1)
	k.setClaimArgument(ctx, claimID, argument.ID)
	k.setUserArgument(ctx, creator, argument.ID)

	if claim.FirstArgumentTime.Equal(time.Time{}) {
		err = k.claimKeeper.SetFirstArgumentTime(ctx, claimID, ctx.BlockHeader().Time)
		if err != nil {
			return Argument{}, err
		}
	}

	switch {
	case stakeType == StakeBacking:
		err := k.claimKeeper.AddBackingStake(ctx, claimID, creationAmount)
		if err != nil {
			return Argument{}, err
		}
	case stakeType == StakeChallenge:
		err := k.claimKeeper.AddChallengeStake(ctx, claimID, creationAmount)
		if err != nil {
			return Argument{}, err
		}
	}

	return argument, nil
}

func (k Keeper) Argument(ctx sdk.Context, argumentID uint64) (Argument, bool) {
	argument := Argument{}
	bz := k.store(ctx).Get(argumentKey(argumentID))
	if bz == nil {
		return Argument{}, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &argument)
	return argument, true
}

func (k Keeper) MarkUnhelpfulArgument(ctx sdk.Context, argumentID uint64) sdk.Error {
	arg, ok := k.Argument(ctx, argumentID)
	if !ok {
		return ErrCodeUnknownArgument(argumentID)
	}
	arg.IsUnhelpful = true
	k.setArgument(ctx, arg)

	return nil
}

func (k Keeper) DownvoteArgument(ctx sdk.Context, argumentID uint64) sdk.Error {
	arg, ok := k.Argument(ctx, argumentID)
	if !ok {
		return ErrCodeUnknownArgument(argumentID)
	}
	arg.DownvotedCount++
	k.setArgument(ctx, arg)
	return nil
}

func (k Keeper) SetStakeExpired(ctx sdk.Context, stakeID uint64) sdk.Error {
	stake, ok := k.Stake(ctx, stakeID)
	if !ok {
		return ErrCodeUnknownStake(stakeID)
	}
	stake.Expired = true
	k.setStake(ctx, stake)
	return nil
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.StakingAdmins) > 0 && !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.StakingAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.StakingAdmins = append(params.StakingAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.StakingAdmins {
		if currentAdmin.Equals(admin) {
			params.StakingAdmins = append(params.StakingAdmins[:i], params.StakingAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).StakingAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func (k Keeper) setArgument(ctx sdk.Context, argument Argument) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(argument)
	k.store(ctx).Set(argumentKey(argument.ID), bz)
}

var tierLimitsEarnedCoins = []sdk.Int{
	sdk.NewInt(app.Shanev * 10),
	sdk.NewInt(app.Shanev * 20),
	sdk.NewInt(app.Shanev * 30),
	sdk.NewInt(app.Shanev * 40),
	sdk.NewInt(app.Shanev * 50),
}

var tierLimitsStakeAmounts = []sdk.Int{
	sdk.NewInt(app.Shanev * 1000),
	sdk.NewInt(app.Shanev * 1500),
	sdk.NewInt(app.Shanev * 2000),
	sdk.NewInt(app.Shanev * 2500),
	sdk.NewInt(app.Shanev * 3000),
}

var defaultStakeLimit = sdk.NewInt(app.Shanev * 500)
var defaultMinimumBalance = sdk.NewInt(app.Shanev * 50)

func (k Keeper) checkStakeThreshold(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) sdk.Error {
	balance := k.bankKeeper.GetCoins(ctx, address).AmountOf(app.StakeDenom)
	if balance.IsZero() {
		return sdk.ErrInsufficientFunds("Insufficient coins")
	}
	p := k.GetParams(ctx)
	period := p.Period

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
	if balance.Sub(amount).LT(defaultMinimumBalance) {
		return ErrCodeMinBalance()
	}

	switch totalEarned := k.TotalEarnedCoins(ctx, address); {
	// if total earned >= 50
	case totalEarned.GTE(tierLimitsEarnedCoins[4]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[4]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 40
	case totalEarned.GTE(tierLimitsEarnedCoins[3]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[3]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 30
	case totalEarned.GTE(tierLimitsEarnedCoins[2]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[2]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 20
	case totalEarned.GTE(tierLimitsEarnedCoins[1]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[1]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 10
	case totalEarned.GTE(tierLimitsEarnedCoins[0]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[0]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	default:
		if staked.Add(amount).GT(defaultStakeLimit) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	}
}

func (k Keeper) TotalEarnedCoins(ctx sdk.Context, creator sdk.AccAddress) sdk.Int {
	earnedCoins := k.getEarnedCoins(ctx, creator)
	total := sdk.NewInt(0)
	for _, e := range earnedCoins {
		total = total.Add(e.Amount)
	}
	return total
}

func (k Keeper) newStake(ctx sdk.Context, amount sdk.Coin, creator sdk.AccAddress,
	stakeType StakeType, argumentID uint64, communityID string) (Stake, sdk.Error) {
	if !stakeType.Valid() {
		return Stake{}, ErrCodeInvalidStakeType(stakeType)
	}
	err := k.checkStakeThreshold(ctx, creator, amount.Amount)
	if err != nil {
		return Stake{}, err
	}
	period := k.GetParams(ctx).Period
	stakeID, err := k.stakeID(ctx)
	if err != nil {
		return Stake{}, err
	}

	_, err = k.bankKeeper.SubtractCoin(ctx, creator, amount,
		argumentID, stakeType.BankTransactionType(), WithCommunityID(communityID),
		ToModuleAccount(UserStakesPoolName),
	)
	if err != nil {
		return Stake{}, err
	}

	stake := Stake{
		ID:          stakeID,
		ArgumentID:  argumentID,
		CommunityID: communityID,
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
	k.setCommunityStake(ctx, communityID, stake.ID)
	k.setUserCommunityStake(ctx, stake.Creator, communityID, stakeID)
	return stake, nil
}

func (k Keeper) Stake(ctx sdk.Context, stakeID uint64) (Stake, bool) {
	stake := Stake{}
	bz := k.store(ctx).Get(stakeKey(stakeID))
	if bz == nil {
		return stake, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &stake)
	return stake, true
}

func (k Keeper) setStake(ctx sdk.Context, stake Stake) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stake)
	k.store(ctx).Set(stakeKey(stake.ID), bz)
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

func (k Keeper) getEarnedCoins(ctx sdk.Context, user sdk.AccAddress) sdk.Coins {
	earnedCoins := sdk.Coins{}
	bz := k.store(ctx).Get(userEarnedCoinsKey(user))
	if bz == nil {
		return sdk.NewCoins()
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &earnedCoins)
	return earnedCoins
}

func (k Keeper) setEarnedCoins(ctx sdk.Context, user sdk.AccAddress, earnedCoins sdk.Coins) {
	b := k.codec.MustMarshalBinaryLengthPrefixed(earnedCoins)
	k.store(ctx).Set(userEarnedCoinsKey(user), b)
}

func (k Keeper) addEarnedCoin(ctx sdk.Context, user sdk.AccAddress, communityID string, amount sdk.Int) {
	earnedCoins := k.getEarnedCoins(ctx, user)
	earnedCoins = earnedCoins.Add(sdk.NewCoins(sdk.NewCoin(communityID, amount)))
	k.setEarnedCoins(ctx, user, earnedCoins)
}

func (k Keeper) SubtractEarnedCoin(ctx sdk.Context, user sdk.AccAddress, communityID string, amount sdk.Int) {
	earnedCoins := k.getEarnedCoins(ctx, user)
	earnedCoins = earnedCoins.Sub(sdk.NewCoins(sdk.NewCoin(communityID, amount)))
	k.setEarnedCoins(ctx, user, earnedCoins)
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
	k.store(ctx).Set(activeStakeQueueKey(stakeID, endTime), bz)
}

// RemoveFromActiveStakeQueue removes a stakeID from the Active Stake Queue
func (k Keeper) RemoveFromActiveStakeQueue(ctx sdk.Context, stakeID uint64, endTime time.Time) {
	k.store(ctx).Delete(activeStakeQueueKey(stakeID, endTime))
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}

func (k Keeper) UsersEarnings(ctx sdk.Context) []UserEarnedCoins {
	userEarnedCoins := make([]UserEarnedCoins, 0)
	k.IterateUserEarnedCoins(ctx, func(address sdk.AccAddress, coins sdk.Coins) bool {
		userEarnedCoins = append(userEarnedCoins, UserEarnedCoins{
			Address: address,
			Coins:   coins,
		})
		return false
	})
	return userEarnedCoins
}

// EditArgument lets a creator edit an argument as long it hasn't been staked on
func (k Keeper) EditArgument(ctx sdk.Context, body, summary string,
	creator sdk.AccAddress, argumentID uint64) (Argument, sdk.Error) {

	err := k.checkJailed(ctx, creator)
	if err != nil {
		return Argument{}, err
	}

	argument, ok := k.Argument(ctx, argumentID)
	if !ok {
		return Argument{}, ErrCodeUnknownArgument(argumentID)
	}

	isAdmin := k.isAdmin(ctx, creator)

	if !argument.Creator.Equals(creator) && !isAdmin {
		return Argument{}, ErrCodeCannotEditArgumentWrongCreator(argumentID)
	}

	stakes := k.ArgumentStakes(ctx, argumentID)
	if len(stakes) > 1 && !isAdmin {
		return Argument{}, ErrCodeCannotEditArgumentAlreadyStaked(argumentID)
	}

	editedArgument := Argument{
		ID:           argumentID,
		Creator:      argument.Creator,
		ClaimID:      argument.ClaimID,
		CommunityID:  argument.CommunityID,
		Summary:      summary,
		Body:         body,
		StakeType:    argument.StakeType,
		CreatedTime:  argument.CreatedTime,
		UpdatedTime:  argument.UpdatedTime,
		UpvotedStake: argument.UpvotedStake,
		TotalStake:   argument.TotalStake,
		UpvotedCount: argument.UpvotedCount,
		EditedTime:   ctx.BlockHeader().Time,
		Edited:       true,
	}

	k.setArgument(ctx, editedArgument)
	return argument, nil
}
