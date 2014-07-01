package ini

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app"
	goini "github.com/vaughan0/go-ini"
)

func init() {
	instance := new(DefaultIniConfig)

	app.OnAppStart( instance.startup )
	env.RegisterIniConfig( instance )
}

func (it *DefaultIniConfig) startup() error {

	iniFile, _ := goini.LoadFile("ottemo.ini")
	it.iniFileValues = iniFile.Section("")

	err := env.OnConfigIniStart()

	return err
}
