// Package "ini" is a default implementation for "I_IniConfig" interface.
package ini

const (
	INI_GLOBAL_SECTION   = ""  // ini file section name to be used as default section
	ASK_FOR_VALUE_PREFIX = "?" // prefix used before default ini value to be asked in console if not set
)

// I_IniConfig implementer class
type DefaultIniConfig struct {
	iniFilePath string

	iniFileValues  map[string]map[string]string
	currentSection string

	keysToStore map[string]bool
	storeAll    bool
}
