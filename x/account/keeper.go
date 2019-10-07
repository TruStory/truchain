package account

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/supply"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	bankKeeper    BankKeeper
	accountKeeper auth.AccountKeeper
	supplyKeeper  supply.Keeper
}

// NewKeeper creates a new keeper of the auth Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, bankKeeper BankKeeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		bankKeeper,
		accountKeeper,
		supplyKeeper,
	}
}

// CreateAppAccount creates a new account on chain for a user
func (k Keeper) CreateAppAccount(ctx sdk.Context, address sdk.AccAddress,
	coins sdk.Coins, pubKey crypto.PubKey) (appAccnt AppAccount, sdkErr sdk.Error) {

	// first create a base account
	baseAccount := auth.NewBaseAccountWithAddress(address)
	err := baseAccount.SetPubKey(pubKey)
	if err != nil {
		return appAccnt, ErrAppAccountCreateFailed(address)
	}
	k.accountKeeper.SetAccount(ctx, &baseAccount)

	//  then create an app account
	appAccnt = NewAppAccount(address, ctx.BlockHeader().Time)
	k.setAppAccount(ctx, appAccnt)

	// set initial coins
	initialCoinAmount := coins.AmountOf(app.StakeDenom)
	if initialCoinAmount.IsPositive() {
		coin := sdk.NewCoin(app.StakeDenom, initialCoinAmount)
		_, sdkErr := k.bankKeeper.AddCoin(ctx, address, coin, 0, TransactionGift)
		if sdkErr != nil {
			return appAccnt, sdkErr
		}
	} else {
		return appAccnt, sdk.ErrInvalidCoins("Invalid initial coins")
	}

	k.Logger(ctx).Info(fmt.Sprintf("Created %s", appAccnt.String()))

	return appAccnt, nil
}

// AppAccounts returns all app accounts
func (k Keeper) AppAccounts(ctx sdk.Context) (appAccounts []AppAccount) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), AppAccountKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var appAccount AppAccount
		k.codec.MustUnmarshalBinaryBare(iterator.Value(), &appAccount)
		appAccounts = append(appAccounts, appAccount)
	}

	return appAccounts
}

// PrimaryAccount gets the primary base account
func (k Keeper) PrimaryAccount(ctx sdk.Context, addr sdk.AccAddress) (pAcc PrimaryAccount, err sdk.Error) {
	appAcc, ok := k.getAppAccount(ctx, addr)
	if !ok {
		return pAcc, ErrAppAccountNotFound(addr)
	}
	acc := k.accountKeeper.GetAccount(ctx, addr)

	pAcc = PrimaryAccount{
		BaseAccount: auth.BaseAccount{
			Address:       acc.GetAddress(),
			Coins:         acc.GetCoins(),
			PubKey:        acc.GetPubKey(),
			AccountNumber: acc.GetAccountNumber(),
			Sequence:      acc.GetSequence(),
		},
		SlashCount:  appAcc.SlashCount,
		IsJailed:    appAcc.IsJailed,
		JailEndTime: appAcc.JailEndTime,
		CreatedTime: appAcc.CreatedTime,
	}

	return pAcc, nil
}

// JailedAccountsBefore returns all jailed accounts before jailEndTime
func (k Keeper) JailedAccountsBefore(ctx sdk.Context, jailEndTime time.Time) (accounts AppAccounts, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(JailEndTimeAccountPrefix, jailEndTimeAccountsKey(jailEndTime))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr := iterator.Value()
		user, ok := k.getAppAccount(ctx, addr)
		if ok {
			accounts = append(accounts, user)
		}
	}

	return accounts, nil
}

// JailUntil puts an AppAccount in jail until a time
func (k Keeper) JailUntil(ctx sdk.Context, address sdk.AccAddress, until time.Time) sdk.Error {
	user, ok := k.getAppAccount(ctx, address)
	if !ok {
		return ErrAppAccountNotFound(address)
	}

	// delete previous jail time
	if user.IsJailed {
		k.deleteJailEndTimeAccount(ctx, user.JailEndTime, user.Addresses[0])
	}
	user.IsJailed = true
	user.JailEndTime = until

	k.setAppAccount(ctx, user)

	// persist in jail list (sorted by jail end time)
	k.setJailEndTimeAccount(ctx, until, address)

	return nil
}

// UnJail unjails an AppAccount.
func (k Keeper) UnJail(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	user, ok := k.getAppAccount(ctx, address)
	if !ok {
		return ErrAppAccountNotFound(address)
	}
	user.IsJailed = false
	k.deleteJailEndTimeAccount(ctx, user.JailEndTime, user.Addresses[0])
	k.setAppAccount(ctx, user)

	return nil
}

// IsJailed tells whether an AppAccount is jailed by its address
func (k Keeper) IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error) {
	user, ok := k.getAppAccount(ctx, address)
	if !ok {
		return false, ErrAppAccountNotFound(address)
	}

	return user.IsJailed, nil
}

// IncrementSlashCount increments the slash count of the user
func (k Keeper) IncrementSlashCount(ctx sdk.Context, address sdk.AccAddress) (jailed bool, err sdk.Error) {
	user, ok := k.getAppAccount(ctx, address)
	if !ok {
		return false, ErrAppAccountNotFound(address)
	}

	user.SlashCount++
	k.setAppAccount(ctx, user)

	if user.SlashCount >= k.GetParams(ctx).MaxSlashCount {
		jailEndTime := ctx.BlockHeader().Time.Add(k.GetParams(ctx).JailDuration)
		err := k.JailUntil(ctx, user.Addresses[0], jailEndTime)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// IterateAppAccounts iterates over all the stored app accounts and performs a callback function
func (k Keeper) IterateAppAccounts(ctx sdk.Context, cb func(acc AppAccount) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), AppAccountKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var account AppAccount
		k.codec.MustUnmarshalBinaryBare(iterator.Value(), &account)

		if cb(account) {
			break
		}
	}
}

func (k Keeper) getAppAccount(ctx sdk.Context, addr sdk.AccAddress) (acc AppAccount, ok bool) {
	accBytes := k.store(ctx).Get(key(addr))
	if accBytes == nil {
		return
	}
	k.codec.MustUnmarshalBinaryBare(accBytes, &acc)

	return acc, true
}

func (k Keeper) setAppAccount(ctx sdk.Context, acc AppAccount) {
	accBytes := k.codec.MustMarshalBinaryBare(acc)
	k.store(ctx).Set(key(acc.PrimaryAddress()), accBytes)
}

func (k Keeper) setJailEndTimeAccount(ctx sdk.Context, jailEndTime time.Time, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(jailEndTimeAccountKey(jailEndTime, addr), addr)
}

func (k Keeper) deleteJailEndTimeAccount(ctx sdk.Context, jailEndTime time.Time, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(jailEndTimeAccountKey(jailEndTime, addr))
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}
