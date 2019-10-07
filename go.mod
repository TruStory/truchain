module github.com/TruStory/truchain

require (
	github.com/cosmos/cosmos-sdk v0.37.0
	github.com/gorilla/mux v1.7.3
	github.com/magiconair/properties v1.8.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.32.3
	github.com/tendermint/tm-db v0.2.0
	github.com/tendermint/tmlibs v0.9.0
	golang.org/x/text v0.3.2 // indirect
)

go 1.13

replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk
