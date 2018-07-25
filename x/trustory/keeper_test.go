package trustory

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// for incode address generation
func testAddr(addr string) sdk.AccAddress {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	return res
}

// dummy vars used for testing
var (
	addrs = []sdk.AccAddress{
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6160"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6161"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6162"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6163"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6164"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6165"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6166"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6167"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6168"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6169"),
	}

	stories = []string{
		"@jaekwon invted proof-of-stake",
		"@bitconnect is not a ponzi scheme",
		"@bucky loves pizza",
		"@zaki invented proof-of-work",
	}
)

func TestTruStoryKeeper(t *testing.T) {
	keyStake := sdk.NewKVStoreKey("stake")
	keyStory := sdk.NewKVStoreKey("trustory")

	mapp, k, _ := CreateMockApp(100, keyStake, keyStory)

	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.NewContext(false, abci.Header{})
	if ctx.KVStore(k.TruStory) == nil {
		panic("Nil interface")
	}

	// create stories
	story := NewStory(int64(1), stories[1], addrs[1], ctx.BlockHeight())

	// ---- Test SetStory ----
	err := k.SetStory(ctx, 1, story)
	assert.Nil(t, err)
}

// CreateMockApp creates a new Mock application for testing
func CreateMockApp(
	numGenAccs int,
	stakeKey *sdk.KVStoreKey,
	storyKey *sdk.KVStoreKey,
) (
	*mock.App,
	Keeper,
	stake.Keeper,
) {
	mapp := mock.NewApp()

	stake.RegisterWire(mapp.Cdc)
	RegisterWire(mapp.Cdc)

	// coin keeper
	ck := bank.NewKeeper(mapp.AccountMapper)

	// stake keeper
	sk := stake.NewKeeper(mapp.Cdc, stakeKey, ck, mapp.RegisterCodespace(stake.DefaultCodespace))

	// story keeper
	keeper := NewKeeper(storyKey, ck, sk, mapp.RegisterCodespace(DefaultCodespace))
	mapp.Router().AddRoute("trustory", NewHandler(keeper))

	return mapp, keeper, sk
}
