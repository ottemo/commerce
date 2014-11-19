// Package errors is a default implementation of InterfaceErrorBus declared in
// "github.com/ottemo/foundation/env" package
package errors

import (
	"github.com/ottemo/foundation/env"
	"regexp"
)

// Package global constants
const (
	// flag to indicate that stack trace collection required
	ConstCollectStack = true
)

// Package global variables
var (
	// regular expression used to parse error message mask (error level and error code, encodes in message)
	ConstMsgRegexp = regexp.MustCompile(`\s*[\[{(]?\s*(?:([0-9]+)?[-: ]([0-9a-fA-F]+)?)?\s*[\]})]?\s*[:\->]*\s*(.+)`)
)

// InterfaceErrorBus implementer class
type DefaultErrorBus struct {
	listeners []env.FuncErrorListener
}

// InterfaceOttemoError implementer class
type OttemoError struct {
	Message string
	Code    string
	Level   int

	Stack string

	handled bool
}
