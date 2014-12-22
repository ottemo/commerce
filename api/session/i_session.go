package session

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns current session id
func (it *DefaultSession) GetID() string {
	if it == nil {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401030a4ca49540a4e2e72cc736", "nil session instance")
		return ""
	}

	return it.id
}

// Get returns session value by a given key or nil - if not set
func (it *DefaultSession) Get(key string) interface{} {
	if it == nil {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401030a4ca49540a4e2e72cc736", "nil session instance")
		return nil
	}

	if value, ok := it.Data[key]; ok == true {
		return value
	}
	return nil
}

// Set assigns value to session key
func (it *DefaultSession) Set(key string, value interface{}) {
	if it == nil {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401030a4ca49540a4e2e72cc736", "nil session instance")
	}

	it.Data[key] = value
	it.UpdatedAt = time.Now()
}

// Close clears session data
func (it *DefaultSession) Close() error {
	if it == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401030a4ca49540a4e2e72cc736", "nil session instance")
	}

	eventData := map[string]interface{}{"session": it, "sessionID": it.GetID()}
	env.Event("session.close", eventData)

	sessionID := it.GetID()

	// removing array references
	if _, present := sessionService.Sessions[sessionID]; present {
		sessionService.sessionsMutex.Lock()
		delete(sessionService.Sessions, sessionID)
		sessionService.sessionsMutex.Unlock()
	}

	// removing storage references
	fileName := ConstStorageFolder + sessionID
	if _, err := os.Stat(fileName); os.IsExist(err) {
		err := os.Remove(fileName)
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// Save stores session for a long period
func (it *DefaultSession) Save() error {
	if it == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401030a4ca49540a4e2e72cc736", "nil session instance")
	}

	sessionID := it.GetID()
	if sessionID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "af9b84a18c82443c919f8dd75f956887", "session id is blank")
	}

	sessionFile, err := os.OpenFile(ConstStorageFolder+sessionID, os.O_WRONLY|os.O_CREATE, 0660)
	defer sessionFile.Close()

	if err != nil {
		return env.ErrorDispatch(err)
	}

	var writer io.Writer = sessionFile
	if ConstCryptSession {
		writer, err = utils.EncryptWriter(writer)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	jsonEncoder := json.NewEncoder(writer)
	err = jsonEncoder.Encode(it)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Load restores session from a long period storage
func (it *DefaultSession) Load(sessionID string) error {
	if sessionID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "af9b84a18c82443c919f8dd75f956887", "session id is blank")
	}
	filename := ConstStorageFolder + sessionID

	// loading session data to instance
	_, err := os.Stat(filename)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	sessionFile, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	defer sessionFile.Close()

	var reader io.Reader = sessionFile
	if ConstCryptSession {
		reader, err = utils.EncryptReader(reader)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	if it == nil {
		it = new(DefaultSession)
	}

	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(it)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.id = sessionID
	it.UpdatedAt = time.Now()

	// updating sessions cache
	sessionService.sessionsMutex.Lock()
	sessionService.Sessions[sessionID] = it
	sessionService.sessionsMutex.Unlock()

	return nil
}
