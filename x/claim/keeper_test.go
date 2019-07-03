package claim

import (
	"net/url"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAddGetClaim(t *testing.T) {
	ctx, keeper := mockDB()

	// test getting a non-existent claim
	claim, ok := keeper.Claim(ctx, uint64(5))
	assert.False(t, ok)
	assert.Equal(t, Claim{}, claim)

	claim = createFakeClaim(ctx, keeper)

	// test getting an existing claim
	_, ok = keeper.Claim(ctx, claim.ID)
	assert.True(t, ok)

	// test incrementing id by adding another story
	claim = createFakeClaim(ctx, keeper)
	assert.Equal(t, uint64(2), claim.ID)

	claims := keeper.CreatorClaims(ctx, claim.Creator)
	assert.Len(t, claims, 2)

	claims = keeper.CommunityClaims(ctx, claim.CommunityID)
	assert.Len(t, claims, 2)

	claims = keeper.ClaimsBetweenIDs(ctx, 0, 100)
	assert.Len(t, claims, 2)

	claims = keeper.ClaimsBetweenIDs(ctx, 2, 100)
	assert.Len(t, claims, 1)

	tt := time.Now().UTC()
	claims = keeper.ClaimsAfterTime(ctx, tt)
	assert.Len(t, claims, 0)

	tt = tt.Add(-60 * time.Minute)
	claims = keeper.ClaimsAfterTime(ctx, tt)
	assert.Len(t, claims, 2)

	claims = keeper.ClaimsBeforeTime(ctx, tt)
	assert.Len(t, claims, 0)

	tt = tt.Add(60 * 60 * time.Minute)
	claims = keeper.ClaimsBeforeTime(ctx, tt)
	assert.Len(t, claims, 2)
}

func TestClaims(t *testing.T) {
	ctx, keeper := mockDB()

	for i := 0; i <= 1000; i++  {
		createFakeClaim(ctx, keeper)
	}

	claims := keeper.Claims(ctx)
	assert.Equal(t, uint64(1001), claims[0].ID)
	assert.Equal(t, uint64(1000), claims[1].ID)
	assert.Equal(t, uint64(2), claims[999].ID)
	assert.Equal(t, uint64(1), claims[1000].ID)
}

func TestCommunityClaims(t *testing.T) {
	ctx, keeper := mockDB()

	for i := 0; i <= 1000; i++  {
		createFakeClaim(ctx, keeper)
	}

	claims := keeper.CommunityClaims(ctx, "crypto")
	assert.Equal(t, uint64(1001), claims[0].ID)
	assert.Equal(t, uint64(1000), claims[1].ID)
	assert.Equal(t, uint64(2), claims[999].ID)
	assert.Equal(t, uint64(1), claims[1000].ID)
}

func TestCreatorClaims(t *testing.T) {
	ctx, keeper := mockDB()

	for i := 0; i <= 1000; i++  {
		createFakeClaim(ctx, keeper)
	}

	creator := sdk.AccAddress([]byte{1, 2})

	claims := keeper.CreatorClaims(ctx, creator)
	assert.Equal(t, uint64(1001), claims[0].ID)
	assert.Equal(t, uint64(1000), claims[1].ID)
	assert.Equal(t, uint64(2), claims[999].ID)
	assert.Equal(t, uint64(1), claims[1000].ID)
}

func TestCreatedTimeClaims(t *testing.T) {
	ctx, keeper := mockDB()

	for i := 0; i <= 1000; i++  {
		createFakeClaim(ctx, keeper)
	}

	createdTime := time.Now().UTC()

	claims := keeper.ClaimsBeforeTime(ctx, createdTime)
	assert.Equal(t, uint64(1), claims[0].ID)
	assert.Equal(t, uint64(2), claims[1].ID)
	assert.Equal(t, uint64(100), claims[99].ID)
	assert.Equal(t, uint64(1000), claims[999].ID)
}

func createFakeClaim(ctx sdk.Context, keeper Keeper) Claim {
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})
	body := "Preethi can handle liquor better than Aamir."
	communityID := "crypto"
	creator := sdk.AccAddress([]byte{1, 2})
	source := url.URL{}

	claim, err := keeper.SubmitClaim(ctx, body, communityID, creator, source)
	if err != nil {
		panic(err)
	}

	return claim
}
