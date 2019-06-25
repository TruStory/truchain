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
