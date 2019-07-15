package params

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params errors reserve 800 ~ 899.
const (
	DefaultCodespace     sdk.CodespaceType = ModuleName
	ErrorCodeJSONParsing sdk.CodeType      = 803
)

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		fmt.Sprintf("JSON parsing error: %s", err.Error()))
}
