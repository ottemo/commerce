package model_product

package db_sqlite

import (
	"errors"
	"github.com/ottemo/platform/interfaces/database"
	"github.com/ottemo/platform/interfaces/config"

	"github.com/ottemo/platform/tools/module_manager"

	"code.google.com/p/go-sqlite/go1/sqlite3"
)

func init() {
	module_manager.RegisterModule( new(SQLiteEngine) )
}

// Structures declaration
//-----------------------

type ProductModel struct { }


// I_Module interface implementation
//----------------------------------

func (it *ProductModel) GetModuleName() string { return "Sqlite3" }
func (it *ProductModel) GetModuleDepends() []string { return make([]string, 0) }

func (it *ProductModel) ModuleMakeSysInit() error { return nil }

func (it *ProductModel) ModuleMakeConfig() error {
	return nil
}

func (it *ProductModel) ModuleMakeInit() error {
}

func (it *ProductModel) ModuleMakeVerify() error {
	return nil
}

func (it *ProductModel) ModuleMakeLoad() error { return nil }
func (it *ProductModel) ModuleMakeInstall() error { return nil }
func (it *ProductModel) ModuleMakePostInstall() error { return nil }




// I_DBStorage interface implementation
//----------------------------------

func (it *ProductModel) GetModelName() string {
	return "DefaultProductModel"
}
