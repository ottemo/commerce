package impex

import (
	"github.com/ottemo/foundation/env"
)

// registers new command to import/export system
func RegisterImportCommand(commandName string, command InterfaceImpexImportCmd) error {
	if _, present := importCmd[commandName]; present {
		return env.ErrorNew(commandName + " already registered in impex")
	}

	importCmd[commandName] = command

	return nil
}

// un-registers command from import/export system
func UnRegisterImportCommand(commandName string) error {
	if _, present := importCmd[commandName]; !present {
		return env.ErrorNew("can't find registered command " + commandName)
	}

	delete(importCmd, commandName)

	return nil
}
