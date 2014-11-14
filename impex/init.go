package impex

import (
	"github.com/ottemo/foundation/api"
)

// module entry point
func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	RegisterImportCommand("INSERT", new(I_ImpexImportCmdInsert))
	RegisterImportCommand("UPDATE", new(I_ImpexImportCmdUpdate))
	RegisterImportCommand("DELETE", new(I_ImpexImportCmdDelete))

	RegisterImportCommand("STORE", new(I_ImpexImportCmdStore))
	RegisterImportCommand("MEDIA", new(I_ImpexImportCmdMedia))

	RegisterImportCommand("ATTRIBUTE_ADD", new(I_ImpexImportCmdAttributeAdd))
}
