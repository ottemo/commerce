// Package "errors" is a default implementation for "I_ErrorBus" interface.
package errors

import (
	"github.com/ottemo/foundation/env"
	"regexp"
)

const (
	// flag to indicate that stack trace collection required
	COLLECT_STACK = true
)

var (
	// regular expression used to parse error message mask (error level and error code, encodes in message)
	MSG_REGEXP = regexp.MustCompile(`\s*[\[{(]?\s*(?:([0-9]+)?[-: ]([0-9a-fA-F]+)?)?\s*[\]})]?\s*[:\->]*\s*(.+)`)
)

// I_ErrorBus implementer class
type DefaultErrorBus struct {
	listeners []env.F_ErrorListener
}

// I_OttemoError implementer class
type OttemoError struct {
	Message string
	Code    string
	Level   int

	Stack string

	handled bool
}
