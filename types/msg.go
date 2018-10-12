package types

import (
	"encoding/json"
	"path"
	"reflect"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Utilities for all `sdk.Msg` types

// GetType returns the package name of the containing `Msg`
func GetType(msg sdk.Msg) string {
	pkgPath := reflect.TypeOf(msg).PkgPath()
	return path.Base(pkgPath)
}

// GetName returns the name of the `Msg` in snake_case
func GetName(msg sdk.Msg) string {
	name := reflect.TypeOf(msg).Name()
	prefix := strings.Split(toSnakeCase(name), "_")
	return strings.Join(prefix[:len(prefix)-1], "_")
}

// MustGetSignBytes serializes a `Msg` type into json bytes.
func MustGetSignBytes(msg sdk.Msg) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners returns the signers of a message
func GetSigners(addr sdk.AccAddress) []sdk.AccAddress {
	return []sdk.AccAddress{addr}
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
