package database

import ("errors")

// Interfaces declaration
//-----------------------

type I_DBEngine interface {
	I_DBStorage
}

type I_DBStorage interface {
	GetCollection(Name string) (I_DBCollection, error)
	GetCollectionFor(Object I_DBObject) (I_DBCollection, error)

}

type I_DBObject interface {
	Load(id string) error
	Save() error
	Delete() error
}

type I_DBCollection interface {
	Save( HashMap map[string]interface{} ) error
	SaveObject( Object I_MappableObject ) error

	Load() error

	AddFilter(Attribute string, Operator string, Value interface{}) error

	ListAttrubutes() []string
}

type I_MappableObject interface {
	ImportAttrubutes(map[string]interface{}) error
	ExportAttributes(map[string]interface{}) error
}


// Delegate routines
//------------------

var dbEngines = map[string]I_DBEngine{}
var currentDbEngine string

func GetDBEngine() I_DBEngine {
	return dbEngines[currentDbEngine]
}

func RegisterDatabaseEngine(Name string, Engine I_DBEngine) error {
	if _, present := dbEngines[Name]; present {
		errors.New("DB engine [" + Name + "] already registered")
	} else {
		dbEngines[Name] = Engine
	}

	return nil
}

func UnregisterDatabaseEngine(Name string) error {
	if _, present := dbEngines[Name]; present {
		delete(dbEngines, Name)
	} else {
		errors.New("can not find registered DB engine [" + Name + "]")
	}
	return nil
}
