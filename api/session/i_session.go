package session

import (
	"time"
)

// implementer class
type Session struct {
	id string

	values map[string]interface{}

	time time.Time
}

// returns current session id
func (it *Session) GetId() string {
	return it.id
}

// returns session value by a given key or nil - if not set
func (it *Session) Get(key string) interface{} {
	if value, ok := it.values[key]; ok == true {
		return value
	}
	return nil
}

// assigns value to session key
func (it *Session) Set(key string, value interface{}) {
	it.values[key] = value
}

// clears session data
func (it *Session) Close() error {
	sessionsMutex.Lock()

	delete(Sessions, it.GetId())

	sessionsMutex.Unlock()

	return nil
}
