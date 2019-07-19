package account

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keys for account store
// Items are stored with the following key: values
//
// - 0x00<AccAddress>: AppAccount
//
// - 0x10<jailEndTime_Bytes><AccAddress>: AccAddress
var (
	AppAccountKeyPrefix = []byte{0x00}

	JailEndTimeAccountPrefix = []byte{0x10}
)

func key(addr sdk.AccAddress) []byte {
	return append(AppAccountKeyPrefix, addr.Bytes()...)
}

func jailEndTimeAccountsKey(endTime time.Time) []byte {
	return append(JailEndTimeAccountPrefix, sdk.FormatTimeBytes(endTime)...)
}

func jailEndTimeAccountKey(endTime time.Time, addr sdk.AccAddress) []byte {
	return append(jailEndTimeAccountsKey(endTime), addr.Bytes()...)
}
