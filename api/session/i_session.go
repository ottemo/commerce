package session

import (
	"encoding/gob"
	"github.com/ottemo/foundation/env"
	"os"
)

// GetID returns current session id
func (it *DefaultSession) GetID() string {
	return it.ID
}

// Get returns session value by a given key or nil - if not set
func (it *DefaultSession) Get(key string) interface{} {
	if value, ok := it.Data[key]; ok == true {
		return value
	}
	return nil
}

// Set assigns value to session key
func (it *DefaultSession) Set(key string, value interface{}) {
	it.Data[key] = value
}

// Close clears session data
func (it *DefaultSession) Close() error {

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
	sessionID := it.GetID()

	sessionFile, err := os.OpenFile(ConstStorageFolder+sessionID, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	defer sessionFile.Close()

	if err != nil {
		return env.ErrorDispatch(err)
	}

	encoder := gob.NewEncoder(sessionFile)

	err = encoder.Encode(it)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Load restores session from a long period storage
func (it *DefaultSession) Load(sessionID string) error {
	filename := ConstStorageFolder + sessionID

	_, err := os.Stat(filename)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	sessionFile, err := os.OpenFile(ConstStorageFolder+sessionID, os.O_RDONLY, 0660)
	defer sessionFile.Close()

	decoder := gob.NewDecoder(sessionFile)

	err = decoder.Decode(it)
	if err != nil {
		it.ID = sessionID
		return env.ErrorDispatch(err)
	}

	return nil
}
