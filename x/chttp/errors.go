package chttp

import (
	"encoding/hex"
	"fmt"

	tcmn "github.com/tendermint/tendermint/libs/common"
)

func signBytesMismatchError(received tcmn.HexBytes, canonical tcmn.HexBytes) error {
	s := "Tx Error: Bytes signed were \"%s\" but canonical tx encoding would be \"%s\" (utf8: \"%s\")"
	cs, _ := hex.DecodeString(canonical.String())
	return fmt.Errorf(s, received.String(), canonical.String(), cs)
}

func expectedMessagesError(receivedCount int, providedTypes []string) error {
	s := "Tx Error: Got %i messages in expected sequence %v"
	return fmt.Errorf(s, receivedCount, providedTypes)
}

func internalDecodingError(detail string) error {
	s := "Internal Decoding Error: %s"
	return fmt.Errorf(s, detail)
}

func unsupportedMsgTypeError(name string, supported []string) error {
	s := "Tx Error: Unsupported message type \"%s\" (supported: %v)"
	return fmt.Errorf(s, name, supported)
}

func unsupportedAlgoError(name string, supported []string) error {
	s := "Tx Error: Unsupported public key algorithm \"%s\" (supported: %v)"
	return fmt.Errorf(s, name, supported)
}
