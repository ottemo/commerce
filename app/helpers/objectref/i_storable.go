package objectref

import "github.com/ottemo/foundation/env"

// GetID returns current object id
func (it *DBObjectRef) GetID() string {
	return it.id
}

// SetID sets new id to current object
func (it *DBObjectRef) SetID(id string) {
	it.id = id
}

// Save stores current product to DB
func (it *DBObjectRef) Save() error {
	return env.ErrorNew("not implemented")
}

// Load loads information from DB
func (it *DBObjectRef) Load(id string) error {
	return env.ErrorNew("not implemented")
}

// Delete removes current object instance from DB
func (it *DBObjectRef) Delete() error {
	return env.ErrorNew("not implemented")
}
