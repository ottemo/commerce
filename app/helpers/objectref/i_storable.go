package objectref

import "github.com/ottemo/foundation/env"

// returns current object id
func (it *DBObjectRef) GetId() string {
	return it.id
}

// sets new id to current object
func (it *DBObjectRef) SetId(id string) {
	it.id = id
}

// stores current product to DB
func (it *DBObjectRef) Save() error {
	return env.ErrorNew("not implemented")
}

// loads information from DB
func (it *DBObjectRef) Load(id string) error {
	return env.ErrorNew("not implemented")
}

// removes current object instance from DB
func (it *DBObjectRef) Delete() error {
	return env.ErrorNew("not implemented")
}
