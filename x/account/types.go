package account

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ auth.Account = (*AppAccount)(nil)

// Defines auth module constants
const (
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
)

// AppAccount is the main account for a TruStory user.
type AppAccount struct {
	*auth.BaseAccount

	SlashCount  uint      `json:"slash_count"`
	IsJailed    bool      `json:"is_jailed"`
	JailEndTime time.Time `json:"jail_end_time"`
	CreatedTime time.Time `json:"created_time"`
}

// AppAccounts is a slice of AppAccounts
type AppAccounts []AppAccount
