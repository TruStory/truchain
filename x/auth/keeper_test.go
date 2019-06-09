package auth

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	appAccount := keeper.NewAppAccount(ctx, address, coins, publicKey)

	assert.Equal(t, appAccount.Address, address)
	assert.Equal(t, appAccount.Coins, coins)
	assert.Equal(t, appAccount.PubKey, publicKey)
}

func TestAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey)

	returnedAppAccount, err := keeper.AppAccount(ctx, createdAppAccount.BaseAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount.BaseAccount, createdAppAccount.BaseAccount)
}

func TestAppAccounts_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()
	_ = keeper.NewAppAccount(ctx, address, coins, publicKey)

	_, publicKey2, address2, coins2, _ := getFakeAppAccountParams()
	_ = keeper.NewAppAccount(ctx, address2, coins2, publicKey2)

	all := keeper.AppAccounts(ctx)
	assert.Len(t, all, 2)
}

func TestJailUntil_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey)
	isJailed, err := keeper.IsJailed(ctx, createdAppAccount.BaseAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, false, isJailed)

	keeper.JailUntil(ctx, createdAppAccount.BaseAccount.Address, time.Now().AddDate(0, 0, 10))
	isJailed, err = keeper.IsJailed(ctx, createdAppAccount.BaseAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, true, isJailed)
}

func TestAddToEarnedStake_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey)
	assert.Len(t, createdAppAccount.EarnedStake, 0)

	earnedCoin := EarnedCoin{sdk.NewCoin("testcoin", sdk.NewInt(10)), uint64(1)}

	// adding for the first time
	keeper.AddToEarnedStake(ctx, createdAppAccount.Address, earnedCoin)
	returnedAppAccount, err := keeper.AppAccount(ctx, createdAppAccount.Address)
	assert.Nil(t, err)
	assert.Len(t, returnedAppAccount.EarnedStake, 1)
	assert.Equal(t, returnedAppAccount.EarnedStake[0], earnedCoin)

	// adding again
	keeper.AddToEarnedStake(ctx, createdAppAccount.Address, earnedCoin)
	returnedAppAccount, err = keeper.AppAccount(ctx, createdAppAccount.Address)
	assert.Nil(t, err)
	assert.Len(t, returnedAppAccount.EarnedStake, 1)
	assert.Equal(t, returnedAppAccount.EarnedStake[0].Coin, earnedCoin.Add(earnedCoin.Coin))
}

func TestIncrementSlashCount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey)
	assert.Equal(t, createdAppAccount.SlashCount, uint(0))

	// incrementing once
	keeper.IncrementSlashCount(ctx, createdAppAccount.Address)
	returnedAppAccount, err := keeper.AppAccount(ctx, createdAppAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount.SlashCount, uint(1))

	// incrementing again
	keeper.IncrementSlashCount(ctx, createdAppAccount.Address)
	returnedAppAccount, err = keeper.AppAccount(ctx, createdAppAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount.SlashCount, uint(2))
}
