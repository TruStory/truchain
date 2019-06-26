package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()

	appAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)

	assert.Equal(t, appAccount.Address, address)
	assert.Equal(t, appAccount.Coins, coins)
	assert.Equal(t, appAccount.PubKey, publicKey)

	assert.Equal(t, false, appAccount.IsJailed)
}

func TestJailUntil_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()

	createdAppAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)
	isJailed, err := keeper.IsJailed(ctx, createdAppAccount.GetAddress())
	assert.Nil(t, err)
	assert.Equal(t, false, isJailed)

	err = keeper.JailUntil(ctx, createdAppAccount.GetAddress(), time.Now().AddDate(0, 0, 10))
	assert.NoError(t, err)
	isJailed, err = keeper.IsJailed(ctx, createdAppAccount.GetAddress())
	assert.Nil(t, err)
	assert.Equal(t, true, isJailed)

	accounts, err := keeper.JailedAccounts(ctx, time.Now().AddDate(0, 0, 10))
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)

	err = keeper.JailUntil(ctx, createdAppAccount.GetAddress(), time.Now().AddDate(0, 0, 10))
	accounts, _ = keeper.JailedAccounts(ctx, time.Now().AddDate(0, 0, 110))
	assert.Len(t, accounts, 2)

	accounts, err = keeper.JailedAccounts(ctx, time.Now())
	assert.NoError(t, err)
	assert.Len(t, accounts, 0)
}

func TestIncrementSlashCount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()

	createdAppAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)
	assert.Equal(t, createdAppAccount.SlashCount, uint(0))

	// incrementing once
	keeper.IncrementSlashCount(ctx, createdAppAccount.Address)
	returnedAppAccount, err := keeper.getAccount(ctx, address)

	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount.SlashCount, uint(1))

	// incrementing again
	keeper.IncrementSlashCount(ctx, createdAppAccount.Address)
	returnedAppAccount, err = keeper.getAccount(ctx, address)
	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount.SlashCount, uint(2))
}
