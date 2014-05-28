package db_sqlite

import ("github.com/ottemo/platform/interfaces/database")

type SQLiteCollection struct {

}

func (it *SQLiteCollection) Save( HashMap map[string]interface{} ) error { return nil }

func (it *SQLiteCollection) SaveObject( Object database.I_MappableObject ) error { return nil }

func (it *SQLiteCollection) Load() error { return nil }

func (it *SQLiteCollection) AddFilter(Attribute string, Operator string, Value interface{}) error { return nil }

func (it *SQLiteCollection) ListAttrubutes() []string { return []string{} }
