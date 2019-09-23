package account

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// Defines auth module constants
const (
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName

	EventTypeUnjailedAccount = "unjailed_account"
	AttributeKeyUser         = "user"
)

type PrimaryAccount struct {
	auth.BaseAccount

	SlashCount  int       `json:"slash_count"`
	IsJailed    bool      `json:"is_jailed"`
	JailEndTime time.Time `json:"jail_end_time"`
	CreatedTime time.Time `json:"created_time"`
}

// AppAccount is the main account for a TruStory user.
type AppAccount struct {
	Addresses   []sdk.AccAddress `json:"addresses"`
	SlashCount  int              `json:"slash_count"`
	IsJailed    bool             `json:"is_jailed"`
	JailEndTime time.Time        `json:"jail_end_time"`
	CreatedTime time.Time        `json:"created_time"`
}

func NewAppAccount(address sdk.AccAddress, createdTime time.Time) AppAccount {
	return AppAccount{
		Addresses:   []sdk.AccAddress{address},
		SlashCount:  0,
		IsJailed:    false,
		JailEndTime: time.Time{},
		CreatedTime: createdTime,
	}
}

func (acc AppAccount) PrimaryAddress() sdk.AccAddress {
	return acc.Addresses[0]
}

// String implements fmt.Stringer
func (acc AppAccount) String() string {
	return fmt.Sprintf(`
  Address:           %s
  SlashCount:        %d
  IsJailed:          %t
  JailEndTime:       %s
  CreatedTime:       %s`,
		acc.PrimaryAddress().String(), acc.SlashCount, acc.IsJailed, acc.JailEndTime.String(), acc.CreatedTime.String())
}

// AppAccounts is a slice of AppAccounts
type AppAccounts []AppAccount
