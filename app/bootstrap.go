package app

import (
	"bufio"
	"encoding/csv"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TruStory/truchain/x/category"

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

	spew.Dump("[DEBUG] CREATOR ADDRESS", addr)

	key, err := chttp.StdKey("ed25519", pubKey.Bytes())
	if err != nil {
		panic(err)
	}

	err = bacc.SetPubKey(key)
	if err != nil {
		panic(err)
	}

	coins, _ := sdk.ParseCoins("50000000trusteak, 30000000btc, 10000000shitcoin")

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
	ck category.ReadKeeper,
	creator sdk.AccAddress,
	claim string,
	catSlug string,
	source string,
	argument string) int64 {

	categories, _ := ck.GetAllCategories(ctx)

	var catID int64
	for _, category := range categories {
		slug := strings.ToLower(category.Slug)
		if slug == catSlug {
			catID = category.ID
		}
	}

	storyType := story.Default
	sourceURL, _ := url.Parse(source)

	// fake a block time
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	evidenceURLs := []story.Evidence{}

	spew.Dump("[DEBUG] ADDING STORY", claim, catID)

	storyID, _ := sk.Create(ctx, argument, claim, catID, creator, evidenceURLs, *sourceURL, storyType)

	return storyID
}

func loadTestDB(
	ctx sdk.Context,
	storyKeeper story.WriteKeeper,
	accountKeeper auth.AccountKeeper,
	backingKeeper backing.WriteKeeper,
	categoryKeeper category.ReadKeeper,
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

	addr1 := createUser(ctx, accountKeeper)
	addr2 := createUser(ctx, accountKeeper)
	addr3 := createUser(ctx, accountKeeper)

	for _, record := range records[1:] {
		claim := record[0]
		catSlug := record[1]
		source := record[2]
		argument := record[3]
		createStory(ctx, storyKeeper, categoryKeeper, addr1, claim, catSlug, source, argument)
	}

	// get the 1st story
	story, _ := storyKeeper.Story(ctx, 1)

	coins := bankKeeper.GetCoins(ctx, addr1)
	spew.Dump("DEBUG", coins)

	// back it
	amount, _ := sdk.ParseCoin("100000trusteak")
	argument := "this is an argument"
	duration := backing.DefaultMsgParams().MinPeriod
	testURL, _ := url.Parse("http://www.trustory.io")
	evidence := []url.URL{*testURL}

	_, err = backingKeeper.Create(ctx, story.ID, amount, argument, addr1, duration, evidence)
	if err != nil {
		panic(err)
	}

	coins = bankKeeper.GetCoins(ctx, addr1)
	spew.Dump("DEBUG", coins)

	// fake a block time
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	// challenge it
	amount, _ = sdk.ParseCoin("200000trusteak")
	challengeID, err := challengeKeeper.Create(ctx, story.ID, amount, argument, addr2, evidence)
	if err != nil {
		panic(err)
	}
	challenge, err := challengeKeeper.Challenge(ctx, challengeID)
	spew.Dump("DEBUG", challenge, err)

	// vote on it
	voteID, err := voteKeeper.Create(ctx, story.ID, amount, true, argument, addr3, evidence)
	if err != nil {
		panic(err)
	}
	vote, err := voteKeeper.TokenVote(ctx, voteID)
	spew.Dump("DEBUG", vote, err)
}
