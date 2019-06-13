package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	appAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0)

	assert.Equal(t, appAccount.BaseAccount.Address, address)
	assert.Equal(t, appAccount.BaseAccount.Coins, coins)
	assert.Equal(t, appAccount.BaseAccount.PubKey, publicKey)
}

func TestAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0)

	returnedAppAccount, err := keeper.AppAccount(ctx, createdAppAccount.BaseAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount.BaseAccount, createdAppAccount.BaseAccount)
}

func TestAppAccounts_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()
	appAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0)

	_, publicKey2, address2, coins2, _ := getFakeAppAccountParams()
	appAccount2 := keeper.NewAppAccount(ctx, address2, coins2, publicKey2, 0, 0)

	all := keeper.AppAccounts(ctx)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0].BaseAccount, appAccount.BaseAccount)
	assert.Equal(t, all[1].BaseAccount, appAccount2.BaseAccount)
}

func TestJailUntil_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0)
	isJailed, err := keeper.IsJailed(ctx, createdAppAccount.BaseAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, false, isJailed)
	
	keeper.JailUntil(ctx, createdAppAccount.BaseAccount.Address, time.Now().AddDate(0, 0, 10))
	isJailed, err = keeper.IsJailed(ctx, createdAppAccount.BaseAccount.Address)
	assert.Nil(t, err)
	assert.Equal(t, true, isJailed)
}