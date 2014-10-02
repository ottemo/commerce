package errors

import (
	"fmt"
	"github.com/ottemo/foundation/env"
	"regexp"
)

var (
	MSG_REGEXP = regexp.MustCompile(`\s*[\[{(]?\s*(?:([0-9]+)?[-: ]([0-9a-fA-F]+)?)?\s*[\]})]?\s*[:\->]*\s*(.+)`)
)

type DefaultErrorBus struct {
	listeners []env.F_ErrorListener
}

type OttemoError struct {
	Message string
	Code    string
	Level   int

	handled bool
}

func (it *OttemoError) Error() string {
	return fmt.Sprintf("%d:%s - %s\n", it.Level, it.Code, it.Message)
}
