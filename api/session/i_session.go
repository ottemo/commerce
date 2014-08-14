package session

import (
	"time"
)

type Session struct {
	id string

	values map[string]interface{}

	time time.Time
}



// returns currens session id
func (it *Session) GetId() string {
	return it.id
}



// returns session value by key
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
