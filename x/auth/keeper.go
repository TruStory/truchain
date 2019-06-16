package auth

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
	"github.com/tendermint/tendermint/crypto"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	recordkeeper.RecordKeeper

	codec      *codec.Codec
	paramStore params.Subspace

	bankKeeper BankKeeper
}

// NewKeeper creates a new keeper of the auth Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, bankKeeper BankKeeper) Keeper {
	return Keeper{
		recordkeeper.NewRecordKeeper(storeKey, codec),
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		bankKeeper,
	}
}

// NewAppAccount creates a new account for a user
func (k Keeper) NewAppAccount(
	ctx sdk.Context,
	address sdk.AccAddress, coins sdk.Coins, pubKey crypto.PubKey, accountNumber uint64, sequence uint64,
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

		EarnedStake: nil,
		SlashCount:  0,
		IsJailed:    false,
		JailEndTime: time.Time{}, // zero value of time.Time; to check, use JailEndTime.IsZero()
	}

	k.StringSet(ctx, address.String(), appAccount)
	logger.Info(fmt.Sprintf("Created new appAccount: %s", appAccount.String()))

	k.bankKeeper.NewTransaction(ctx, address, coins)

	return appAccount
}

// AppAccount returns a AppAccount by its address
func (k Keeper) AppAccount(ctx sdk.Context, address sdk.AccAddress) (appAccount AppAccount, err sdk.Error) {
	k.StringGet(ctx, address.String(), &appAccount)
	if appAccount.BaseAccount.Address.Empty() {
		return appAccount, ErrAppAccountNotFound(address)
	}

	return appAccount, nil
}

// AppAccounts gets all AppAccounts from the KVStore
func (k Keeper) AppAccounts(ctx sdk.Context) (appAccounts []AppAccount) {
	err := k.Each(ctx, func(bytes []byte) bool {
		var appAccount AppAccount
		k.codec.MustUnmarshalBinaryLengthPrefixed(bytes, &appAccount)
		appAccounts = append(appAccounts, appAccount)
		return true
	})

	if err != nil {
		return
	}

	return
}

// AddToEarnedStake adds an EarnedCoin to the EarnedStake
func (k Keeper) AddToEarnedStake(ctx sdk.Context, address sdk.AccAddress, earnedCoin EarnedCoin) sdk.Error {
	appAccount, err := k.AppAccount(ctx, address)
	if err != nil {
		return err
	}

	added := false
	for i, current := range appAccount.EarnedStake {
		if current.CommunityID == earnedCoin.CommunityID {
			appAccount.EarnedStake[i].Coin = current.Add(earnedCoin.Coin)
			added = true
		}
	}

	if !added {
		appAccount.EarnedStake = append(appAccount.EarnedStake, earnedCoin)
	}

	k.StringSet(ctx, address.String(), appAccount)

	return nil
}

// JailUntil puts an AppAccount in jail until a time
func (k Keeper) JailUntil(ctx sdk.Context, address sdk.AccAddress, until time.Time) sdk.Error {
	appAccount, err := k.AppAccount(ctx, address)
	if err != nil {
		return err
	}

	appAccount.IsJailed = true
	appAccount.JailEndTime = until
	k.StringSet(ctx, address.String(), appAccount)

	return nil
}

// IsJailed tells whether an AppAccount is jailed by its address
func (k Keeper) IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error) {
	appAccount, err := k.AppAccount(ctx, address)
	if err != nil {
		return false, err
	}

	return appAccount.IsJailed, nil
}

// IncrementSlashCount increments the slash count of the user
func (k Keeper) IncrementSlashCount(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	appAccount, err := k.AppAccount(ctx, address)
	if err != nil {
		return err
	}

	appAccount.SlashCount++
	k.StringSet(ctx, address.String(), appAccount)

	return nil
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
