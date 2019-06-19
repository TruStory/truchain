package claim

import (
	"net/url"
	"testing"
	"time"

	// sdk "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAddGetClaim(t *testing.T) {
	ctx, keeper := mockDB()

	// test getting a non-existent claim
	_, err := keeper.Claim(ctx, uint64(5))
	assert.NotNil(t, err)

	claim := createFakeClaim(ctx, keeper)

	// test getting an existing claim
	_, err = keeper.Claim(ctx, claim.ID)
	assert.NoError(t, err)

	// test incrementing id by adding another story
	claim = createFakeClaim(ctx, keeper)
	assert.Equal(t, uint64(2), claim.ID)
}

func createFakeClaim(ctx sdk.Context, keeper Keeper) Claim {
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})
	body := "Preethi can handle liquor better than Aamir."
	communityID := uint64(1)
	creator := sdk.AccAddress([]byte{1, 2})
	source := url.URL{}

	claim, err := keeper.NewClaim(ctx, body, communityID, creator, source)
	if err != nil {
		panic(err)
	}

	return claim
}
