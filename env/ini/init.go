package ini

import (
	app "github.com/ottemo/foundation/app"
	env "github.com/ottemo/foundation/env"
	ini "github.com/vaughan0/go-ini"
)

func init() {
	instance := new(DefaultIniConfig)

	app.OnAppStart(instance.startup)
	env.RegisterIniConfig(instance)
}

func (it *DefaultIniConfig) startup() error {

	iniFile, _ := ini.LoadFile("ottemo.ini")
	it.iniFileValues = iniFile.Section("")

	err := env.OnConfigIniStart()

	return err
}
