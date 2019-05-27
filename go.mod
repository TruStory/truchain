module github.com/TruStory/truchain

require (
	github.com/cosmos/cosmos-sdk v0.34.4
	github.com/gogo/protobuf v1.1.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/iavl v0.12.2 // indirect
	github.com/tendermint/tendermint v0.31.5
	github.com/tendermint/tmlibs v0.9.0
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
