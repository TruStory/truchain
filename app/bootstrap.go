package app

import (
	"bufio"
	"encoding/csv"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/davecgh/go-spew/spew"

	tru "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tmlibs/cli"
)

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func createUser(
	ctx sdk.Context,
	accountKeeper auth.AccountKeeper) sdk.AccAddress {

	_, pubKey, addr := keyPubAddr()
	bacc := auth.NewBaseAccountWithAddress(addr)

	spew.Dump("DEBUG -- CREATOR ADDRESS", addr)

	key, err := chttp.StdKey("ed25519", pubKey.Bytes())
	if err != nil {
		panic(err)
	}

	err = bacc.SetPubKey(key)
	if err != nil {
		panic(err)
	}

	coins, _ := sdk.ParseCoins("5000000trusteak, 3000000btc, 1000000shitcoin")

	err = bacc.SetCoins(coins)
	if err != nil {
		panic(err)
	}

	acc := tru.NewAppAccount(bacc)

	accountKeeper.SetAccount(ctx, acc)

	return addr
}

func createStory(
	ctx sdk.Context,
	sk story.WriteKeeper,
	creator sdk.AccAddress,
	claim string,
	source string,
	evidence string,
	argument string) int64 {

	catID := int64(1)
	storyType := story.Default
	sourceURL, _ := url.Parse(source)

	// fake a block time
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	url, _ := url.Parse(evidence)
	e := story.Evidence{
		Creator:   creator,
		URL:       *url,
		Timestamp: tru.NewTimestamp(ctx.BlockHeader()),
	}
	evidenceURLs := []story.Evidence{e}

	arg := story.Argument{
		Creator:   creator,
		Body:      argument,
		Timestamp: tru.NewTimestamp(ctx.BlockHeader()),
	}

	arguments := []story.Argument{arg}

	storyID, _ := sk.NewStory(ctx, arguments, claim, catID, creator, evidenceURLs, *sourceURL, storyType)

	return storyID
}

func loadTestDB(
	ctx sdk.Context,
	storyKeeper story.WriteKeeper,
	accountKeeper auth.AccountKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	voteKeeper vote.WriteKeeper,
	gameKeeper game.WriteKeeper,
	bankKeeper bank.Keeper) {

	rootdir := viper.GetString(cli.HomeFlag)
	if rootdir == "" {
		rootdir = DefaultNodeHome
	}

	path := filepath.Join(rootdir, "bootstrap.csv")
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	addr := createUser(ctx, accountKeeper)

	for _, record := range records {
		createStory(ctx, storyKeeper, addr, record[0], record[1], record[2], record[3])
	}

	// get the 1st story
	story, _ := storyKeeper.Story(ctx, 1)

	coins := bankKeeper.GetCoins(ctx, addr)
	spew.Dump("DEBUG", coins)

	// back it
	amount, _ := sdk.ParseCoin("1000trusteak")
	argument := "this is an argument"
	duration := 30 * 24 * time.Hour
	testURL, _ := url.Parse("http://www.trustory.io")
	evidence := []url.URL{*testURL}

	_, err = backingKeeper.Create(ctx, story.ID, amount, argument, addr, duration, evidence)
	if err != nil {
		panic(err)
	}

	coins = bankKeeper.GetCoins(ctx, addr)
	spew.Dump("DEBUG", coins)

	// fake a block time
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	// challenge it
	amount, _ = sdk.ParseCoin("1000trusteak")
	_, err = challengeKeeper.Create(ctx, story.ID, amount, argument, addr, evidence)
	if err != nil {
		panic(err)
	}

	// vote on it
	_, err = voteKeeper.Create(ctx, story.ID, amount, true, argument, addr, evidence)
	if err != nil {
		panic(err)
	}
}
