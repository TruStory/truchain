package account

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ auth.Account = (*AppAccount)(nil)

// Defines auth module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// AppAccount is the main account for a TruStory user.
type AppAccount struct {
	*auth.BaseAccount

	SlashCount  uint
	IsJailed    bool
	JailEndTime time.Time
	CreatedTime time.Time
}

// AppAccounts is a slice of AppAccounts
type AppAccounts []AppAccount
