package truapi

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestHumanReadable(t *testing.T) {
	type testCase struct {
		amount sdk.Int
		output string
	}

	// Do not show decimals if they do not exist
	testCases := []testCase{
		// If greater than 1.0 => show two decimal digits, truncate trailing zeros
		testCase{amount: sdk.NewInt(10000000000), output: "10"},
		testCase{amount: sdk.NewInt(00000000000), output: "0"},
		testCase{amount: sdk.NewInt(2000057891), output: "2"},
		testCase{amount: sdk.NewInt(1100000000), output: "1.1"},
		testCase{amount: sdk.NewInt(1123400000), output: "1.12"},
		// If less than 1.0 => show four decimal digits, truncate trailing zeros
		testCase{amount: sdk.NewInt(100000000), output: "0.1"},
		testCase{amount: sdk.NewInt(10000000), output: "0.01"},
		testCase{amount: sdk.NewInt(123000000), output: "0.123"},
		testCase{amount: sdk.NewInt(123450000), output: "0.1234"},
		testCase{amount: sdk.NewInt(999999999), output: "0.9999"},
		// anything in the 5th decimal places is effectively 0
		testCase{amount: sdk.NewInt(99999), output: "0"},
		testCase{amount: sdk.NewInt(1), output: "0"},
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.output, HumanReadable(sdk.NewCoin("steak", testCase.amount)))
	}
}
