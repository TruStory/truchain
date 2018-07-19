package trustory

// import (
// 	"os"
// 	"testing"

// 	// "github.com/TruStory/cosmos-sdk/x/auth"
// 	"github.com/TruStory/trucoin/app"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	dbm "github.com/tendermint/tmlibs/db"

// 	"github.com/cosmos/cosmos-sdk/x/auth"
// 	"github.com/tendermint/tmlibs/log"
// )

// func loggerAndDB() (log.Logger, dbm.DB) {
// 	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
// 	dB := dbm.NewMemDB()
// 	return logger, dB
// }

// func newTruStoryApp() *app.TruStoryApp {
// 	logger, dB := loggerAndDB()
// 	return app.NewTruStoryApp(logger, dB)
// }

// // hogpodge of all sorts of input required for testing
// func createTestInput(t *testing.T, initCoins int64) (sdk.Context, auth.AccountMapper, Keeper) {

// 	app := newTruStoryApp()
// 	keyStake := sdk.NewKVStoreKey("stake")
// 	keyAuth := sdk.NewKVStoreKey("auth")
// 	simpleGovKey := sdk.NewKVStoreKey("truStory")

// 	// db := dbm.NewMemDB()

// 	db := dbm.NewDB("filesystemDB", dbm.FSDBBackend, "dir")
// 	ms := store.NewCommitMultiStore(db)

// 	app.MountStoreWithDB(keyStake, sdk.StoreTypeIAVL, db)
// 	app.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
// 	app.MountStoreWithDB(simpleGovKey, sdk.StoreTypeIAVL, db)
// 	err := ms.LoadLatestVersion()
// 	require.Nil(t, err)

// 	ctx := app.NewContext(isCheckTx, abci.Header{ChainID: "foochainid"})
// 	app.Mo
// 	accountMapper := auth.NewAccountMapper(
// 		cdc,                 // amino codec
// 		keyAuth,             // target store
// 		&auth.BaseAccount{}, // prototype
// 	)
// 	ck := bank.NewKeeper(accountMapper)
// 	stakeKeeper := stake.NewKeeper(cdc, keyStake, ck, DefaultCodespace)

// 	// fill all the addresses with some coins
// 	for _, addr := range addrs {
// 		ck.AddCoins(ctx, addr, sdk.Coins{
// 			{"Atom", sdk.NewInt(initCoins)},
// 		})
// 	}

// 	keeper := NewKeeper(simpleGovKey, ck, stakeKeeper, DefaultCodespace)
// 	return ctx, accountMapper, keeper
// }
