package slashing

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/auth"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the package slashing module
type Keeper struct {
	recordkeeper.RecordKeeper

	codec      *codec.Codec
	paramStore params.Subspace

	claimKeeper      ClaimKeeper
	stakeKeeper      StakeKeeper
	appAccountKeeper auth.Keeper
}

// NewKeeper creates a new keeper of the slashing Keeper
func NewKeeper(
	storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec,
	claimKeeper ClaimKeeper, stakeKeeper StakeKeeper, appAccountKeeper auth.Keeper,
) Keeper {
	return Keeper{
		recordkeeper.NewRecordKeeper(storeKey, codec),
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		claimKeeper,
		stakeKeeper,
		appAccountKeeper,
	}
}

// NewSlash creates a new slash
func (k Keeper) NewSlash(ctx sdk.Context, stakeID uint64, creator sdk.AccAddress) (slash Slash, err sdk.Error) {
	logger := getLogger(ctx)

	err = k.validateParams(ctx, stakeID, creator)
	if err != nil {
		return
	}

	slash = Slash{
		ID:        k.IncrementID(ctx),
		StakeID:   stakeID,
		Creator:   creator,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	k.Set(ctx, slash.ID, slash)
	logger.Info(fmt.Sprintf("Created new slash: %s", slash.String()))

	return
}

// Slash returns a slash by its ID
func (k Keeper) Slash(ctx sdk.Context, id uint64) (slash Slash, err sdk.Error) {
	err = k.Get(ctx, id, &slash)
	if err != nil {
		return slash, ErrSlashNotFound(id)
	}

	return slash, nil
}

// Slashes gets all slashes from the KVStore
func (k Keeper) Slashes(ctx sdk.Context) (slashes []Slash) {
	err := k.Each(ctx, func(bytes []byte) bool {
		var slash Slash
		k.codec.MustUnmarshalBinaryLengthPrefixed(bytes, &slash)
		slashes = append(slashes, slash)
		return true
	})

	if err != nil {
		return
	}

	return
}

func (k Keeper) validateParams(ctx sdk.Context, stakeID uint64, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// validating stake
	stake, stakeErr := k.stakeKeeper.Stake(ctx, stakeID)
	if stakeErr != nil {
		return ErrInvalidStake(stakeID)
	}
	if stake.SlashCount > params.MaxStakeSlashCount {
		return ErrMaxSlashCountReached(stakeID)
	}

	// validating claim
	claim, claimErr := k.claimKeeper.Claim(ctx, stake.ClaimID)
	if claimErr != nil {
		return ErrInvalidClaim(stake.ClaimID)
	}

	// validating creator
	appAccount, authErr := k.appAccountKeeper.AppAccount(ctx, creator)
	if authErr != nil {
		return ErrInvalidCreator(creator)
	}
	if !hasEnoughEarnedStake(appAccount, claim.CommunityID, params.SlashMinStake) {
		return ErrNotEnoughEarnedStake(creator)
	}

	if !isAdmin(creator, params.SlashAdmins) {
		return ErrInvalidCreator(creator)
	}

	if hasPreviouslySlashed(creator) {
		return ErrAlreadySlashed()
	}

	return nil
}

func hasEnoughEarnedStake(appAccount auth.AppAccount, communityID uint64, minimumEarnedCoins auth.EarnedCoins) bool {
	// first, loop through all the earned coins of the user
	for _, earnedCoin := range appAccount.EarnedStake {

		// for the coin earned for the given community,
		if earnedCoin.CommunityID == communityID {

			// find the minimum required earned coin by looping over the requirements
			for _, minimumEarnedCoin := range minimumEarnedCoins {

				// when found,
				if minimumEarnedCoin.CommunityID == communityID {

					// make sure the user's earned coin is greater than or equal to the requirment
					if earnedCoin.IsGTE(minimumEarnedCoin.Coin) {
						return true
					}
				}
			}
		}
	}

	// in all other cases, return false
	return true
}

func hasPreviouslySlashed(creator sdk.AccAddress) bool {
	// TODO: write the real logic
	return false
}

func isAdmin(address sdk.AccAddress, admins []sdk.AccAddress) bool {
	for _, admin := range admins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
