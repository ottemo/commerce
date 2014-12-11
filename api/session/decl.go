// Package session is a default implementation of InterfaceSessionService and InterfaceSession
// declared in "github.com/ottemo/foundation/api" package
package session

import (
	"github.com/ottemo/foundation/env"
	"sync"
	"time"
)

// Package global constants
const (
	ALPHANUMERIC = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890" // sessionID allowed symbols

	ConstErrorModule = "api/session"
	ConstErrorLevel  = env.ConstErrorLevelService

	ConstStorageFolder = "./var/session/"
)

// Package global variables
var (
	sessionService *DefaultSessionService
)

// DefaultSession is a default implementer of InterfaceSession declared in
// "github.com/ottemo/foundation/api" package
type DefaultSession struct {
	ID        string
	Data      map[string]interface{}
	UpdatedAt time.Time
}

// DefaultSessionService is a default implementer of InterfaceSessionService declared in
// "github.com/ottemo/foundation/api" package
type DefaultSessionService struct {
	Sessions      map[string]*DefaultSession // active sessions set
	sessionsMutex sync.RWMutex               // syncronization on Sessions variable modification
	gcRate        int64                      // garbage collection rate
}
