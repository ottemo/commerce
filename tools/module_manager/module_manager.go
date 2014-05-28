package module_manager

import (
	"errors"
)

type I_Module interface {
	GetModuleName() string
	GetModuleDepends() []string

	ModuleMakeSysInit() error
	ModuleMakeConfig() error
	ModuleMakeInit() error
	ModuleMakeVerify() error
	ModuleMakeLoad() error
	ModuleMakeInstall() error
	ModuleMakePostInstall() error
}

var registeredModules = map[string]I_Module {}

func RegisterModule(Module I_Module) error {
	moduleName := Module.GetModuleName()
	if _, present := registeredModules[moduleName]; !present {
		registeredModules[moduleName] = Module
	} else {
		return errors.New("module '" + moduleName + "' already registered")
	}

	return nil
}

func GetModule(Name string) I_Module {
		return registeredModules[Name]
}

func InitModules() error {
	loadedModules := map[string]bool {}
	modulesLoadOrder := make( []string, 0, len(registeredModules) )

	// determination of modules load order list based on their dependency
	for loadedCount := 1; loadedCount > 0; {
		// loop over all registered modules
		for moduleName, module := range registeredModules {
			// if module have depends - make sure all of them loaded
			dependsLoaded := true
			for _, dependName := range module.GetModuleDepends() {
				if loadedModules[dependName] != true {
					dependsLoaded = false
				}
			}

			// start loading module
			if dependsLoaded {
				loadedModules[moduleName] = true
				modulesLoadOrder = append(modulesLoadOrder, moduleName)
				loadedCount = loadedCount+1

			}
		}
		loadedCount = 0
	}

	if len(modulesLoadOrder) != len(registeredModules) {
		return errors.New("modules dependency can not be resolved")
	}

	for i := 0; i < 7 ; i++ {
		var moduleError error

		for _,moduleName := range modulesLoadOrder {
			module := registeredModules[moduleName]

			switch i {
			case 0:
				moduleError = module.ModuleMakeSysInit()
			case 1:
				moduleError = module.ModuleMakeConfig()
			case 2:
				moduleError = module.ModuleMakeInit()
			case 3:
				moduleError = module.ModuleMakeVerify()
			case 4:
				moduleError = module.ModuleMakeLoad()
			case 5:
				moduleError = module.ModuleMakeInstall()
			case 6:
				moduleError = module.ModuleMakePostInstall()
			}

			if moduleError != nil {
				return moduleError
			}
		}
	}

	return nil
}
