package impex

import (
	"github.com/ottemo/foundation/api"
)

func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	RegisterImportCommand("insert", new(ImpexImportCmdInsert))
	RegisterImportCommand("update", new(ImpexImportCmdUpdate))
	RegisterImportCommand("delete", new(ImpexImportCmdDelete))
	RegisterImportCommand("store", new(ImpexImportCmdStore))
}
