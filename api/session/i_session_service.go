package session

import (
	"crypto/rand"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"math/big"
	"time"
)

// GetName returns implementation name of session service
func (it *DefaultSessionService) GetName() string {
	return "DefaultSessionService"
}

// Get returns session object for given session id or nil of not currently exists
func (it *DefaultSessionService) Get(sessionID string) (api.InterfaceSession, error) {
	if sessionInstance, present := it.Sessions[sessionID]; present {
		return sessionInstance, nil
	}

	sessionInstance := new(DefaultSession)
	err := sessionInstance.Load(sessionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return sessionInstance, nil
}

// New initializes new session instance
func (it *DefaultSessionService) New() (api.InterfaceSession, error) {

	// garbage collecting
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(it.gcRate))
	if err == nil && randomNumber.Cmp(big.NewInt(1)) == 0 {
		it.gc()
	}

	// keeping in memory only allowed amount of sessions
	numOfSessionsToClean := len(it.Sessions) - ConstSessionKeepInMemoryItems
	for sessionID, sessionInstance := range it.Sessions {
		err := sessionInstance.Save()
		if err == nil {
			it.sessionsMutex.Lock()
			delete(it.Sessions, sessionID)
			it.sessionsMutex.Unlock()
		}

		numOfSessionsToClean--
		if numOfSessionsToClean == 0 {
			break
		}
	}

	// receiving new session id
	sessionID, err := it.generateSessionID()
	if err != nil {
		return nil, err
	}

	sessionInstance := new(DefaultSession)
	sessionInstance.id = sessionID
	sessionInstance.Data = make(map[string]interface{})
	sessionInstance.UpdatedAt = time.Now()

	it.Sessions[sessionID] = sessionInstance

	return it.Sessions[sessionID], nil
}

// generateSessionID returns new session id number
func (it *DefaultSessionService) generateSessionID() (string, error) {
	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aab0911497264844984c772fb25dcb88", "can't generate sessionID")
	}

	for i := 0; i < 32; i++ {
		sessionID[i] = ALPHANUMERIC[sessionID[i]%62]
	}

	return string(sessionID), nil
}

// gc removes expired sessions
func (it *DefaultSessionService) gc() {
	for _, sessionInstance := range it.Sessions {
		secondsAfterLastUpdate := time.Now().Sub(sessionInstance.UpdatedAt).Seconds()

		// closing out of date sessions
		if secondsAfterLastUpdate > ConstSessionLifeTime {
			sessionInstance.Close()
		}

		// updating sessions information in a storage
		if secondsAfterLastUpdate > ConstSessionUpdateTime {
			sessionInstance.Save()
		}
	}
}
