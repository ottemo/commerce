// Package ini is a default implementation of I_IniConfig declared in
// "github.com/ottemo/foundation/env" package
package ini

// Package global constants
const (
	INI_GLOBAL_SECTION   = ""  // ini file section name to be used as default section
	ASK_FOR_VALUE_PREFIX = "?" // prefix used before default ini value to be asked in console if not set

	CMD_ARG_STORE_ALL_FLAG = "--iniStoreAll"
	CMD_ARG_SECTION_NAME   = "--iniSection="
	CMD_ARG_TEST_FLAG      = "--test"

	ENVIRONMENT_INI_FILE    = "OTTEMO_INI"
	ENVIRONMENT_INI_SECTION = "OTTEMO_MODE"

	TEST_SECTION_NAME = "test"
	DEFAULT_INI_FILE  = "ottemo.ini"
)

// I_IniConfig implementer class
type DefaultIniConfig struct {
	iniFilePath string

	iniFileValues  map[string]map[string]string
	currentSection string

	keysToStore map[string]bool
	storeAll    bool
}
