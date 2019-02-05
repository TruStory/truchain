package types

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// WriteJSONtoNodeHome creates a .json file with the current state for the given context
// The file is saved in $HOME/.truchaind/bh/i.json where bh represents the
// block height at exectution time, ex: $HOME/.truchaind/1345/story.json
// If the directory does not exists it is created with permissions 0700
// The file is created with permissions 0660
func WriteJSONtoNodeHome(i interface{}, dnh string, bh int64, fn string) {
	b, _ := json.MarshalIndent(i, "", " ")
	path := filepath.Join(dnh, strconv.FormatInt(bh, 10))

	if _, err := os.Stat(path); os.IsNotExist(err) {

		err := os.Mkdir(path, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	fp := filepath.Join(path, fn)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR, 0660)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(b); err != nil {
			log.Fatal(err)
		}
	}
}
