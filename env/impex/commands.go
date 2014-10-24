package impex

import (
	"github.com/ottemo/foundation/env"
)

type ImpexImportCmdInsert struct{}
type ImpexImportCmdUpdate struct{}
type ImpexImportCmdDelete struct{}
type ImpexImportCmdStore struct{}

func (it *ImpexImportCmdInsert) Init(args []string, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdInsert) Process(data []map[string]interface{}, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdUpdate) Init(args []string, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdUpdate) Process(data []map[string]interface{}, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdDelete) Init(args []string, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdDelete) Process(data []map[string]interface{}, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdStore) Init(args []string, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdStore) Process(data []map[string]interface{}, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}
