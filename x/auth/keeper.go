package auth

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	recordkeeper "github.com/shanev/cosmos-record-keeper/recordkeeper"
	"github.com/tendermint/tendermint/crypto"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	recordkeeper.RecordKeeper
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the auth Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		recordkeeper.NewRecordKeeper(storeKey, codec),
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}

// NewAppAccount creates a new account for a user
func (k Keeper) NewAppAccount(
	ctx sdk.Context,
	address sdk.AccAddress, coins sdk.Coins, pubKey crypto.PubKey, accountNumber uint64, sequence uint64, earnedStake EarnedCoins,
) AppAccount {

	logger := getLogger(ctx)

	appAccount := AppAccount{
		BaseAccount: sdkAuth.BaseAccount{
			Address:       address,
			Coins:         coins,
			PubKey:        pubKey,
			AccountNumber: accountNumber,
			Sequence:      sequence,
		},

		ID:          k.IncrementID(ctx),
		EarnedStake: earnedStake,
		SlashCount:  0,
		IsJailed:    false,
		JailEndTime: time.Time{}, // zero value of time.Time; to check, use JailEndTime.IsZero()
	}

	k.Set(ctx, appAccount.ID, appAccount)
	logger.Info(fmt.Sprintf("Created new appAccount: %s", appAccount.String()))

	// TODO: Add a bank transaction with the the initial creation of an AppAccount
	// ...

	return appAccount
}

// AppAccount returns a AppAccount by its ID
func (k Keeper) AppAccount(ctx sdk.Context, id uint64) (appAccount AppAccount, err sdk.Error) {
	err = k.Get(ctx, id, &appAccount)
	if err != nil {
		return appAccount, ErrAppAccountNotFound(id)
	}

	return appAccount, nil
}

// AppAccounts gets all AppAccounts from the KVStore
func (k Keeper) AppAccounts(ctx sdk.Context) (appAccounts []AppAccount) {
	err := k.Each(ctx, func(bytes []byte) bool {
		var appAccount AppAccount
		k.Codec.MustUnmarshalBinaryLengthPrefixed(bytes, &appAccount)
		appAccounts = append(appAccounts, appAccount)
		return true
	})

	if err != nil {
		return
	}

	return
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
