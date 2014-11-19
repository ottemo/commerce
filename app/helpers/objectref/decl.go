// Package objectref intended to unify and simplify a way of model instance changes tracking (currently not implemented)
package objectref

// DBObjectRef is a object state tracking helper and implementer of InterfaceObject and InterfaceStorable
type DBObjectRef struct {
	id string

	loaded   bool
	modified bool

	origData map[string]interface{}
	currData map[string]interface{}
}

// marks object instance as loaded from DB
func (it *DBObjectRef) MarkAsLoaded() {
	it.loaded = true
}

// marks object instance as modified
func (it *DBObjectRef) MarkAsModified() {
	it.modified = true
}

// returns value of modification flag
func (it *DBObjectRef) IsModified() bool {
	return it.modified
}

// returns value of load from DB flag
func (it *DBObjectRef) IsLoaded() bool {
	return it.loaded
}
