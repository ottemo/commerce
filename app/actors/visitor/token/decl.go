// Package token allows to create and use tokens
package token

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameVisitorToken = "visitor_token"

	ConstErrorModule = "visitor/token"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultVisitorCard is a default implementer of InterfaceVisitorCard
type DefaultVisitorCard struct {
	id        string
	visitorID string

	Holder  string
	Payment string

	Type           string
	Number         string
	ExpirationDate string

	Token string
}
