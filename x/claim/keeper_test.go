package claim

import (
	"net/url"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/TruStory/truchain/types"
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

	for i := 0; i <= 1000; i++ {
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

	for i := 0; i <= 1000; i++ {
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

	for i := 0; i <= 1000; i++ {
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

	for i := 0; i <= 1000; i++ {
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

func TestKeeper_AddBackingChallengeStake(t *testing.T) {
	ctx, keeper := mockDB()
	claim := createFakeClaim(ctx, keeper)
	keeper.AddBackingStake(ctx, claim.ID, sdk.NewInt64Coin(app.StakeDenom, 50*app.Shanev))
	keeper.AddChallengeStake(ctx, claim.ID, sdk.NewInt64Coin(app.StakeDenom, 50*app.Shanev))
	keeper.AddChallengeStake(ctx, claim.ID, sdk.NewInt64Coin(app.StakeDenom, 10*app.Shanev))
	keeper.AddChallengeStake(ctx, claim.ID, sdk.NewInt64Coin(app.StakeDenom, 15*app.Shanev))
	keeper.AddBackingStake(ctx, claim.ID, sdk.NewInt64Coin(app.StakeDenom, 35*app.Shanev))
	c, ok := keeper.Claim(ctx, claim.ID)
	assert.True(t, ok)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, 85*app.Shanev).String(), c.TotalBacked.String())
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, 75*app.Shanev).String(), c.TotalChallenged.String())
}

func TestAddAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	creator := keeper.GetParams(ctx).ClaimAdmins[0]
	newAdmin := getFakeAdmin()

	err := keeper.AddAdmin(ctx, newAdmin, creator)
	assert.Nil(t, err)

	newAdmins := keeper.GetParams(ctx).ClaimAdmins
	assert.Subset(t, newAdmins, []sdk.AccAddress{newAdmin})
}

func TestAddAdmin_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidCreator := sdk.AccAddress([]byte{1, 2})
	newAdmin := getFakeAdmin()

	err := keeper.AddAdmin(ctx, newAdmin, invalidCreator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}

func TestRemoveAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	currentAdmins := keeper.GetParams(ctx).ClaimAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.RemoveAdmin(ctx, adminToRemove, adminToRemove) // removing self
	assert.Nil(t, err)
	newAdmins := keeper.GetParams(ctx).ClaimAdmins
	assert.Equal(t, len(currentAdmins)-1, len(newAdmins))
}

func TestRemoveAdmin_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidRemover := sdk.AccAddress([]byte{1, 2})
	currentAdmins := keeper.GetParams(ctx).ClaimAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.AddAdmin(ctx, adminToRemove, invalidRemover)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}
