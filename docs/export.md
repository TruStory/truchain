# Export App State
To export the app's current state:

```bash
$ ./truchaind export > filname.json
```

For a sample state `json` file see [sample.json](./sample.json)

The command exports the entire app's state as defined in `app/genesis.go`:

```go
// GenesisState struct defines the current state
type GenesisState struct {
	Accounts   []GenesisAccount      `json:"accounts"`
	Stories    []story.Story         `json:"stories"`
	Categories []category.Category   `json:"categories"`
	Backings   backing.GenesisState  `json:"backings"`
	Challenges []challenge.Challenge `json:"challenges"`
	Games      game.GenesisState     `json:"games"`
	Votes      []vote.TokenVote      `json:"votes"`
}
```

The relevant key-value store keys for the genesis state are defined in application type in `app/app.go`:
```go
	// create your application type
	var app = &TruChain{
		categories:         categories,
		codec:              codec,
		BaseApp:            bam.NewBaseApp(params.AppName, logger, db, auth.DefaultTxDecoder(codec), options...),
		keyMain:            sdk.NewKVStoreKey("main"),
		keyAccount:         sdk.NewKVStoreKey("acc"),
		keyIBC:             sdk.NewKVStoreKey("ibc"),
		keyStory:           sdk.NewKVStoreKey("stories"),
		keyCategory:        sdk.NewKVStoreKey("categories"),
		keyBacking:         sdk.NewKVStoreKey("backings"),
		keyBackingList:     sdk.NewKVStoreKey("backingList"),
		keyChallenge:       sdk.NewKVStoreKey("challenges"),
		keyFee:             sdk.NewKVStoreKey("collectedFees"),
		keyGame:            sdk.NewKVStoreKey("game"),
		keyPendingGameList: sdk.NewKVStoreKey("pendingGameList"),
		keyGameQueue:       sdk.NewKVStoreKey("gameQueue"),
		keyVote:            sdk.NewKVStoreKey("vote"),
		api:                nil,
		apiStarted:         false,
		blockCtx:           nil,
		blockHeader:        abci.Header{},
		registrarKey:       loadRegistrarKey(),
	}
```


## Key-Value Store Key
Each individual store sate is defined in its corresponding `genesis.go` file. For example, the `backing`  state is defined in `x/backing/genesis.go`:
```go
// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Backings     []Backing `json:"backings"`
	BackingsList []int64   `json:"backing_list"`
}
```

which implements the keys defined above in the application type. 
