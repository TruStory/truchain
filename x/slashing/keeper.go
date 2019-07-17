package slashing

import (
	"fmt"

	"github.com/TruStory/truchain/x/claim"

	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/staking"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/bank"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the package slashing module
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	bankKeeper    bank.Keeper
	stakingKeeper staking.Keeper
	accountKeeper account.Keeper
	claimKeeper   claim.Keeper
}

// NewKeeper creates a new keeper of the slashing Keeper
func NewKeeper(
	storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec,
	bankKeeper bank.Keeper, stakingKeeper staking.Keeper, accountKeeper account.Keeper, claimKeeper claim.Keeper,
) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		bankKeeper,
		stakingKeeper,
		accountKeeper,
		claimKeeper,
	}
}

// CreateSlash creates a new slash on an argument (mark as "Unhelpful" in app)
func (k Keeper) CreateSlash(ctx sdk.Context, argumentID uint64, slashType SlashType, slashReason SlashReason, slashDetailedReason string, creator sdk.AccAddress) (slash Slash, err sdk.Error) {
	logger := getLogger(ctx)

	err = k.validateParams(ctx, argumentID, slashDetailedReason, creator)
	if err != nil {
		return
	}

	slashID, err := k.slashID(ctx)
	if err != nil {
		return
	}

	slash = Slash{
		ID:             slashID,
		ArgumentID:     argumentID,
		Type:           slashType,
		Reason:         slashReason,
		DetailedReason: slashDetailedReason,
		Creator:        creator,
		CreatedTime:    ctx.BlockHeader().Time,
	}

	// persist the slash
	k.setSlash(ctx, slash)
	// increment the slash ID for the next slash
	k.setSlashID(ctx, slashID+1)
	// persist associations
	k.setCreatorSlash(ctx, creator, slashID)
	k.incrementSlashCount(ctx, argumentID)
	k.setArgumentSlash(ctx, argumentID, slashID)
	k.setArgumentSlasherSlash(ctx, argumentID, slashID, creator)

	err = k.stakingKeeper.DownvoteArgument(ctx, argumentID)
	if err != nil {
		return slash, err
	}

	slashCount := k.getSlashCount(ctx, argumentID)
	if slashCount >= k.GetParams(ctx).MinSlashCount || k.isAdmin(ctx, creator) {
		err = k.stakingKeeper.MarkUnhelpfulArgument(ctx, argumentID)
		if err != nil {
			return slash, err
		}
		err = k.punish(ctx, argumentID)
		if err != nil {
			return slash, err
		}
	}

	logger.Info(fmt.Sprintf("Created new slash: %s", slash.String()))

	return
}

func (k Keeper) punish(ctx sdk.Context, argumentID uint64) sdk.Error {
	stakingPool := sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())
	var communityID string
	for _, stake := range k.stakingKeeper.ArgumentStakes(ctx, argumentID) {
		communityID = stake.CommunityID
		stakingPool = stakingPool.Add(stake.Amount)
		if k.stakingKeeper.IsStakeActive(ctx, stake.ID, stake.EndTime) {
			k.stakingKeeper.RemoveFromActiveStakeQueue(ctx, stake.ID, stake.EndTime)
		} else {
			if stake.Result != nil {
				stakeInterest := stake.Result.ArgumentCreatorReward.Add(stake.Result.StakeCreatorReward)
				_, err := k.bankKeeper.SafeSubtractCoin(
					ctx,
					stake.Creator,
					stakeInterest,
					stake.ID,
					bank.TransactionInterestSlashed,
					WithCommunityID(communityID))
				if err != nil {
					return err
				}
			}
		}
		slashMagnitude := int64(k.GetParams(ctx).SlashMagnitude)
		slashCoin := sdk.NewCoin(app.StakeDenom, stake.Amount.Amount.MulRaw(slashMagnitude))

		_, err := k.bankKeeper.SafeSubtractCoin(
			ctx,
			stake.Creator,
			slashCoin,
			stake.ID,
			bank.TransactionStakeSlashed,
			WithCommunityID(communityID))
		if err != nil {
			return err
		}

		if stake.Type == staking.StakeBacking {
			err = k.claimKeeper.SubtractBackingStake(ctx, stake.ID, stake.Amount)
			if err != nil {
				return err
			}
		}
		if stake.Type == staking.StakeChallenge {
			err = k.claimKeeper.SubtractChallengeStake(ctx, stake.ID, stake.Amount)
			if err != nil {
				return err
			}
		}

		// increment slash count for user (and jail if needed)
		_, err = k.accountKeeper.IncrementSlashCount(ctx, stake.Creator)
		if err != nil {
			return err
		}
	}

	if !stakingPool.IsPositive() {
		return sdk.ErrInsufficientCoins("staking pool cannot be empty")
	}

	// reward curators who marked "unhelpful"
	curatorShareDec := k.GetParams(ctx).CuratorShare
	totalCuratorAmountDec := stakingPool.Amount.ToDec().Mul(curatorShareDec)

	slashes := k.ArgumentSlashes(ctx, argumentID)
	curatorAmount := totalCuratorAmountDec.QuoInt64(int64(len(slashes))).TruncateInt()
	curatorCoin := sdk.NewCoin(app.StakeDenom, curatorAmount)
	for _, slash := range slashes {
		_, err := k.bankKeeper.AddCoin(
			ctx,
			slash.Creator,
			curatorCoin,
			slash.ID,
			bank.TransactionCuratorReward,
			WithCommunityID(communityID))
		if err != nil {
			return err
		}
	}

	return nil
}

// Slash returns a slash by its ID
func (k Keeper) Slash(ctx sdk.Context, id uint64) (slash Slash, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	slashBytes := store.Get(key(id))
	if slashBytes == nil {
		return slash, ErrSlashNotFound(id)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(slashBytes, &slash)

	return slash, nil
}

// Slashes gets all slashes from the KVStore
func (k Keeper) Slashes(ctx sdk.Context) (slashes []Slash) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, SlashesKeyPrefix)

	return k.iterate(iterator)
}

// slashID gets the highest slash ID
func (k Keeper) slashID(ctx sdk.Context) (slashID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(SlashIDKey)
	if bz == nil {
		return 0, ErrSlashNotFound(slashID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &slashID)
	return slashID, nil
}

// setSlash sets a slash in store
func (k Keeper) setSlash(ctx sdk.Context, slash Slash) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(slash)
	store.Set(key(slash.ID), bz)
}

// set the slash ID
func (k Keeper) setSlashID(ctx sdk.Context, slashID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(slashID)
	store.Set(SlashIDKey, bz)
}

// sets the association between the creator and the slash
func (k Keeper) setCreatorSlash(ctx sdk.Context, creator sdk.AccAddress, slashID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(slashID)
	store.Set(creatorSlashKey(creator, slashID), bz)
}

// increments the slash count for a given stake
func (k Keeper) incrementSlashCount(ctx sdk.Context, stakeID uint64) {
	k.setSlashCount(ctx, stakeID, uint64(k.getSlashCount(ctx, stakeID)+1))
}

// sets the association between the stake and the slash count
func (k Keeper) setSlashCount(ctx sdk.Context, stakeID uint64, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(count)
	store.Set(slashCountKey(stakeID), bz)
}

// getSlashCount gets the number of slashes for a stake
func (k Keeper) getSlashCount(ctx sdk.Context, stakeID uint64) (count int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(slashCountKey(stakeID))
	if bz == nil {
		return 0
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	return count
}

// iterate iterates over the iterator
func (k Keeper) iterate(iterator sdk.Iterator) (slashes Slashes) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var slash Slash
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &slash)
		slashes = append(slashes, slash)
	}

	return
}

func (k Keeper) validateParams(ctx sdk.Context, argumentID uint64, detailedReason string, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	_, ok := k.stakingKeeper.Argument(ctx, argumentID)
	if !ok {
		return ErrInvalidArgument(argumentID)
	}

	if k.getSlashCount(ctx, argumentID) > params.MinSlashCount {
		return ErrMaxSlashCountReached(argumentID)
	}

	if len(detailedReason) > params.MaxDetailedReasonLength {
		return ErrInvalidSlashReason(fmt.Sprintf("Detailed reason must be under %d chars.", params.MaxDetailedReasonLength))
	}

	// validating creator
	isAdmin := k.isAdmin(ctx, creator)
	hasEnoughCoins := k.hasEnoughEarnedStake(ctx, creator, params.SlashMinStake)

	if !isAdmin && !hasEnoughCoins {
		return ErrNotEnoughEarnedStake(creator)
	}

	if !isAdmin && k.hasPreviouslySlashed(ctx, argumentID, creator) {
		return ErrAlreadySlashed()
	}

	return nil
}

func (k Keeper) hasEnoughEarnedStake(ctx sdk.Context, address sdk.AccAddress, requirement sdk.Coin) bool {
	balance := k.bankKeeper.GetCoins(ctx, address)

	return balance.AmountOf(app.StakeDenom).GTE(requirement.Amount)
}

func (k Keeper) hasPreviouslySlashed(ctx sdk.Context, argumentID uint64, creator sdk.AccAddress) bool {
	slashes := k.ArgumentSlashes(ctx, argumentID)
	for _, slash := range slashes {
		if slash.Creator.Equals(creator) {
			return true
		}
	}

	return false
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).SlashAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

// setArgumentSlash sets a argument <-> slash association in store
func (k Keeper) setArgumentSlash(ctx sdk.Context, argumentID, slashID uint64) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(slashID)
	k.store(ctx).Set(argumentSlashKey(argumentID, slashID), bz)
}

func (k Keeper) ArgumentSlashes(ctx sdk.Context, argumentID uint64) []Slash {
	slashes := make([]Slash, 0)
	k.IterateArgumentSlashes(ctx, argumentID, func(slash Slash) bool {
		slashes = append(slashes, slash)
		return false
	})
	return slashes
}

func (k Keeper) IterateArgumentSlashes(ctx sdk.Context, argumentID uint64, cb slashCallback) {
	iterator := k.store(ctx).Iterator(ArgumentSlashesPrefix, sdk.PrefixEndBytes(argumentSlashPrefix(argumentID)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var slashID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &slashID)
		slash, err := k.Slash(ctx, slashID)
		if err != nil {
			panic(err)
		}
		if cb(slash) {
			break
		}
	}
}

func (k Keeper) setArgumentSlasherSlash(ctx sdk.Context, argumentID, slashID uint64, slasher sdk.AccAddress) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(slashID)
	k.store(ctx).Set(argumentSlasherSlashKey(argumentID, slasher, slashID), bz)
}

func (k Keeper) ArgumentSlasherSlashes(ctx sdk.Context, slasher sdk.AccAddress, argumentID uint64) []Slash {
	slashes := make([]Slash, 0)
	k.IterateArgumentSlasherSlashes(ctx, argumentID, slasher, func(slash Slash) bool {
		slashes = append(slashes, slash)
		return false
	})
	return slashes
}

type slashCallback func(slash Slash) (stop bool)

func (k Keeper) IterateArgumentSlasherSlashes(ctx sdk.Context, argumentID uint64, address sdk.AccAddress, cb slashCallback) {
	iterator := k.store(ctx).Iterator(ArgumentCreatorPrefix, sdk.PrefixEndBytes(argumentSlasherPrefix(argumentID, address)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var slash Slash
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &slash)
		if cb(slash) {
			break
		}
	}
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
