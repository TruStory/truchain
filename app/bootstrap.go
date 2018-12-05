package app

import (
	"bufio"
	"encoding/csv"
	"net/url"
	"os"
	"path/filepath"
	"time"

	tru "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
)

func createUser() {
	// bacc := auth.NewBaseAccountWithAddress(msg.Address)

}

func createStory(
	ctx sdk.Context,
	sk story.WriteKeeper,
	claim string,
	source string,
	evidence string,
	argument string) int64 {

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	catID := int64(1)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	sourceURL, _ := url.Parse(source)

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

func loadTestDB(ctx sdk.Context, storyKeeper story.WriteKeeper) {
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

	for _, record := range records {
		createStory(ctx, storyKeeper, record[0], record[1], record[2], record[3])
	}
}
