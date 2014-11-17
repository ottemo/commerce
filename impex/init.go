package impex

import (
	"github.com/ottemo/foundation/api"
)

// init makes package self-initialization routine
func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	RegisterImportCommand("INSERT", new(ImpexImportCmdInsert))
	RegisterImportCommand("UPDATE", new(ImpexImportCmdUpdate))
	RegisterImportCommand("DELETE", new(ImpexImportCmdDelete))

	RegisterImportCommand("STORE", new(ImpexImportCmdStore))
	RegisterImportCommand("MEDIA", new(ImpexImportCmdMedia))

	RegisterImportCommand("ATTRIBUTE_ADD", new(ImpexImportCmdAttributeAdd))
}
