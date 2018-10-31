package chttp

import (
	"fmt"
)

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
