package truchain

// import (
// 	"testing"
// 	"time"

// 	ts "github.com/TruStory/truchain/x/truchain/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/auth"
// 	"github.com/cosmos/cosmos-sdk/x/bank"
// 	"github.com/go-kit/kit/log"
// )

// func TestCannotUnjailUnlessJailed(t *testing.T) {
// 	// initial setup
// 	ctx, ck, sk, _, keeper := createTestInput(t)
// 	slh := NewHandler(keeper)
// 	amtInt := int64(100)
// 	addr, val, amt := addrs[0], pks[0], sdk.NewInt(amtInt)
// 	msg := newTestMsgCreateValidator(addr, val, amt)
// 	got := stake.NewHandler(sk)(ctx, msg)
// 	require.True(t, got.IsOK())
// 	stake.EndBlocker(ctx, sk)
// 	require.Equal(t, ck.GetCoins(ctx, sdk.AccAddress(addr)), sdk.Coins{{sk.GetParams(ctx).BondDenom, initCoins.Sub(amt)}})
// 	require.True(t, sdk.NewDecFromInt(amt).Equal(sk.Validator(ctx, addr).GetPower()))

// 	// assert non-jailed validator can't be unjailed
// 	got = slh(ctx, NewMsgUnjail(addr))
// 	require.False(t, got.IsOK(), "allowed unjail of non-jailed validator")
// 	require.Equal(t, sdk.ToABCICode(DefaultCodespace, CodeValidatorNotJailed), got.Code)
// }

// func TestSubmitStoryMsg(t *testing.T) {
// 	ctx, _, _, k := mockDB()

// 	h := NewHandler(k)

// 	body := "fake story"
// 	cat := ts.DEX
// 	creator := sdk.AccAddress([]byte{1, 2})
// 	storyType := ts.Default
// 	msg := ts.NewSubmitStoryMsg(body, cat, creator, storyType)

// 	// _, err := k.GetBacking(ctx, id)
// 	// assert.NotNil(t, err)
// 	// assert.Equal(t, ts.ErrBackingNotFound(id).Code(), err.Code(), "Should get error")
// }

// func mockDB() (sdk.Context, sdk.MultiStore, auth.AccountMapper, TruKeeper) {
// 	ms, accKey, storyKey, backingKey := setupMultiStore()
// 	cdc := makeCodec()
// 	am := auth.NewAccountMapper(cdc, accKey, auth.ProtoBaseAccount)
// 	ck := bank.NewBaseKeeper(am)
// 	k := NewTruKeeper(storyKey, backingKey, ck, cdc)

// 	// create fake context with fake block time in header
// 	// time := time.Date(2018, time.September, 14, 23, 0, 0, 0, time.UTC)
// 	time := time.Now().Add(5 * time.Hour)
// 	header := abci.Header{Time: time}
// 	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

// 	return ctx, ms, am, k
// }
