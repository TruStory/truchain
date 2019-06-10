package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, earnedCoins := getFakeAppAccountParams()
	
	appAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0, earnedCoins)

	assert.NotZero(t, appAccount.ID)
	assert.Equal(t, appAccount.BaseAccount.Address, address)
	assert.Equal(t, appAccount.BaseAccount.Coins, coins)
	assert.Equal(t, appAccount.BaseAccount.PubKey, publicKey)
	assert.Equal(t, appAccount.EarnedStake, earnedCoins)
}

func TestAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, earnedCoins := getFakeAppAccountParams()
	
	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0, earnedCoins)

	returnedAppAccount, err := keeper.AppAccount(ctx, createdAppAccount.ID)
	assert.Nil(t, err)
	assert.Equal(t, returnedAppAccount, createdAppAccount)
}

func TestAppAccounts_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, earnedCoins := getFakeAppAccountParams()
	appAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0, earnedCoins)

	_, publicKey2, address2, coins2, earnedCoins2 := getFakeAppAccountParams()
	appAccount2 := keeper.NewAppAccount(ctx, address2, coins2, publicKey2, 0, 0, earnedCoins2)

	all := keeper.AppAccounts(ctx)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], appAccount)
	assert.Equal(t, all[1], appAccount2)
}
