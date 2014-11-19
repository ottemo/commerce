// Package errors is a default implementation of InterfaceErrorBus declared in
// "github.com/ottemo/foundation/env" package
package errors

import (
	"github.com/ottemo/foundation/env"
	"regexp"
)

// Package global constants
const (
	ConstCollectStack = true // flag to indicate that stack trace collection required
)

// Package global variables
var (
	// ConstMsgRegexp is a regular expression used to parse error message mask (error level and error code, encodes in message)
	ConstMsgRegexp = regexp.MustCompile(`\s*[\[{(]?\s*(?:([0-9]+)?[-: ]([0-9a-fA-F]+)?)?\s*[\]})]?\s*[:\->]*\s*(.+)`)
)

// DefaultErrorBus InterfaceErrorBus implementer class
type DefaultErrorBus struct {
	listeners []env.FuncErrorListener
}

// OttemoError @reconcile@ InterfaceOttemoError implementer class
type OttemoError struct {
	Message string
	Code    string
	Level   int

	Stack string

	handled bool
}
