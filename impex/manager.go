package impex

import (
	"github.com/ottemo/foundation/env"
)

// RegisterImportCommand registers new command to import/export system
func RegisterImportCommand(commandName string, command InterfaceImpexImportCmd) error {
	if _, present := importCmd[commandName]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f4f50cfc-3de6-4c76-8518-02ad26790f44", commandName+" already registered in impex")
	}

	importCmd[commandName] = command

	return nil
}

// UnRegisterImportCommand un-registers command from import/export system
func UnRegisterImportCommand(commandName string) error {
	if _, present := importCmd[commandName]; !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a9df49e3-0d06-4afa-8e3a-7cde5a3b2a1d", "can't find registered command "+commandName)
	}

	delete(importCmd, commandName)

	return nil
}
