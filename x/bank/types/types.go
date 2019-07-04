package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


const (
	ModuleName        = "trubank2"
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
)

// Association list keys
var (
	AccountKey = sdk.NewKVStoreKey("account")
)
