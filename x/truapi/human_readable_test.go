package truapi

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	amount sdk.Int
	output string
}

func TestHumanReadable(t *testing.T) {
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
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.output, HumanReadable(sdk.NewCoin("steak", testCase.amount)))
	}
}
