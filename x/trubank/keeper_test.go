package trubank

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestAddCoins(t *testing.T) {
	ctx, k, ck := mockDB()
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})

	k.MintAndAddCoin(ctx, creator, cat.ID, sdk.NewInt(1000))
	k.MintAndAddCoin(ctx, creator, cat.ID, sdk.NewInt(1000))
	k.MintAndAddCoin(ctx, creator, cat.ID, sdk.NewInt(1000))

	cat2, _ := ck.GetCategory(ctx, cat.ID)

	assert.Equal(t, "3000trudex", cat2.TotalCred.String())
}
