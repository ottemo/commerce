package errors

import (
	"fmt"
	"github.com/ottemo/foundation/env"
	"regexp"
)

var (
	MSG_REGEXP    = regexp.MustCompile(`\s*[\[{(]?\s*(?:([0-9]+)?[-: ]([0-9a-fA-F]+)?)?\s*[\]})]?\s*[:\->]*\s*(.+)`)
	INCLUDE_STACK = true
)

type DefaultErrorBus struct {
	listeners []env.F_ErrorListener
}

type OttemoError struct {
	Message string
	Code    string
	Level   int

	Stack string

	handled bool
}

func (it *OttemoError) Error() string {
	message := it.Message
	if it.Stack != "" {
		message += "\n" + it.Stack
	}
	return fmt.Sprintf("%d:%s - %s", it.Level, it.Code, message)
}
