package truapi

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HumanReadable formats a coin amount to be human readable
// Duplicated inside of Octopus until truapi package is moved there
// If modifying, modify Octopus as well
func HumanReadable(coin sdk.Coin) string {
	shanevs := sdk.NewDecFromIntWithPrec(coin.Amount, 9).String()
	parts := strings.Split(shanevs, ".")
	number := parts[0]
	decimal := parts[1]
	// If greater than 1.0 => show two decimal digits, truncate trailing zeros
	displayDecimalPlaces := 2
	if number == "0" {
		// If less than 1.0 => show four decimal digits, truncate trailing zeros
		displayDecimalPlaces = 4
	}
	decimal = strings.TrimRight(decimal, "0")
	numberOfDecimalPlaces := len(decimal)
	if numberOfDecimalPlaces > displayDecimalPlaces {
		numberOfDecimalPlaces = displayDecimalPlaces
	}
	decimal = decimal[0:numberOfDecimalPlaces]
	decimal = strings.TrimRight(decimal, "0")
	if decimal == "" {
		return number
	}
	return fmt.Sprintf("%s%s%s", number, ".", decimal)
}
