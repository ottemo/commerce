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
	return it.id
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
	it.UpdatedAt = time.Now()
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

	filename := ConstStorageFolder + sessionID

	_, err := os.Stat(filename)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	sessionFile, err := os.OpenFile(ConstStorageFolder+sessionID, os.O_RDONLY, 0660)
	defer sessionFile.Close()

	var reader io.Reader = sessionFile
	if ConstCryptSession {
		reader, err = utils.EncryptReader(reader)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(it)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	sessionService.Sessions[sessionID] = it

	return nil
}
