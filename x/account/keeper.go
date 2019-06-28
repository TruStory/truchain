package account

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/crypto"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	bankKeeper    BankKeeper
	accountKeeper auth.AccountKeeper
}

// NewKeeper creates a new keeper of the auth Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, bankKeeper BankKeeper, accountKeeper auth.AccountKeeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		bankKeeper,
		accountKeeper,
	}
}

// NewAppAccount creates a new AppAccount
func NewAppAccount(baseAcct auth.BaseAccount) *AppAccount {
	return &AppAccount{BaseAccount: &baseAcct}
}

// CreateAppAccount creates a new account on chain for a user
func (k Keeper) CreateAppAccount(ctx sdk.Context, address sdk.AccAddress,
	coins sdk.Coins, pubKey crypto.PubKey) (acc *AppAccount, sdkErr sdk.Error) {

	baseAccount := auth.NewBaseAccountWithAddress(address)
	err := baseAccount.SetPubKey(pubKey)
	if err != nil {
		return acc, ErrAppAccountCreateFailed(address)
	}
	err = baseAccount.SetCoins(coins)
	if err != nil {
		return acc, ErrAppAccountCreateFailed(address)
	}
	acc = NewAppAccount(baseAccount)
	acc.CreatedTime = ctx.BlockHeader().Time
	k.accountKeeper.SetAccount(ctx, acc)

	getLogger(ctx).Info(fmt.Sprintf("Created %s", acc.String()))

	_, sdkErr = k.bankKeeper.AddCoin(ctx, address, coins[0], 0, TransactionRegistration)
	if sdkErr != nil {
		return
	}

	return acc, nil
}

// JailedAccountsAfter returns all jailed accounts after jailEndTime
func (k Keeper) JailedAccountsAfter(ctx sdk.Context, jailEndTime time.Time) (accounts AppAccounts, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(jailEndTimeAccountsKey(jailEndTime), sdk.PrefixEndBytes(JailEndTimeAccountPrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr := iterator.Value()
		user, err := k.getAccount(ctx, addr)
		if err != nil {
			panic(err)
		}
		accounts = append(accounts, user)
	}

	return accounts, nil
}

// JailUntil puts an AppAccount in jail until a time
func (k Keeper) JailUntil(ctx sdk.Context, address sdk.AccAddress, until time.Time) sdk.Error {
	user, err := k.getAccount(ctx, address)
	if err != nil {
		return err
	}

	user.IsJailed = true
	user.JailEndTime = until

	k.accountKeeper.SetAccount(ctx, user)

	// persist in jail list (sorted by jail end time)
	k.setJailEndTimeAccount(ctx, until, address)

	return nil
}

// IsJailed tells whether an AppAccount is jailed by its address
func (k Keeper) IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error) {
	user, err := k.getAccount(ctx, address)
	if err != nil {
		return false, err
	}

	return user.IsJailed, nil
}

// IncrementSlashCount increments the slash count of the user
func (k Keeper) IncrementSlashCount(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	user, err := k.getAccount(ctx, address)
	if err != nil {
		return err
	}

	user.SlashCount++
	k.accountKeeper.SetAccount(ctx, user)

	return nil
}

func (k Keeper) getAccount(ctx sdk.Context, addr sdk.AccAddress) (AppAccount, sdk.Error) {
	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc != nil {
		user, ok := acc.(AppAccount)
		if !ok {
			return user, ErrAppAccountNotFound(acc.GetAddress())
		}
		return user, nil
	}

	return AppAccount{}, ErrAppAccountNotFound(addr)
}

func getLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}

func (k Keeper) setJailEndTimeAccount(ctx sdk.Context, jailEndTime time.Time, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(jailEndTimeAccountKey(jailEndTime, addr), addr)
}
