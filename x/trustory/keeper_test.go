package trustory

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// func TestTruStoryKeeper(t *testing.T) {
// 	keyStory := sdk.NewKVStoreKey("trustory")

// }

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
