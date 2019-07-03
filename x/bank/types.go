package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines bank module constants
const (
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
)

// Association list keys
var (
	accountKey = sdk.NewKVStoreKey("account")
)
