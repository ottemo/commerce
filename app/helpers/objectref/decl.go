package objectref

type DBObjectRef struct {
	id string

	loaded   bool
	modified bool

	origData map[string]interface{}
	currData map[string]interface{}
}

func (it *DBObjectRef) GetId() {
	return it.id
}

func (it *DBObjectRef) SetId(id string) {
	it.id = id
}

func (it *DBObjectRef) MarkAsLoaded() {
	it.loaded = true
}

func (it *DBObjectRef) MarkAsModified() {
	it.modified
}

func (it *DBObjectRef) IsModified() {
	return it.modified
}

func (it *DBObjectRef) IsLoaded() {
	return it.loaded
}
