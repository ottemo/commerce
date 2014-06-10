package default_ini_config

import (
	config "github.com/ottemo/foundation/config"
	app "github.com/ottemo/foundation/app"
	ini "github.com/vaughan0/go-ini"
)

func init() {
	instance := new(DefaultIniConfig)

	app.OnAppStart( instance.startup )
	config.RegisterIniConfig( instance )
}

func (it *DefaultIniConfig) startup() error {

	iniFile, _ := ini.LoadFile("ottemo.ini")
	it.iniFileValues = iniFile.Section("")

	err := config.OnConfigIniStart()

	return err
}
