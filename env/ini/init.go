package ini

import (
	"io/ioutil"
	"os"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	goini "github.com/vaughan0/go-ini"
)

const (
	INI_CONFIG_FILE = "ottemo.ini"
)

func init() {
	instance := new(DefaultIniConfig)

	instance.iniFileValues = make(map[string]string)
	instance.keysToStore = make([]string, 0)

	app.OnAppInit(instance.init)
	app.OnAppStart(instance.startup)

	env.RegisterIniConfig(instance)
}

func (it *DefaultIniConfig) startup() error {

	if len(it.keysToStore) > 0 {
		iniData := ""
		for _, key := range it.keysToStore {
			iniData += key + "=" + it.iniFileValues[key] + "\n"
		}

		ioerr := ioutil.WriteFile(INI_CONFIG_FILE, []byte(iniData), os.ModePerm)
		if ioerr != nil {
			return ioerr
		}
	}

	return nil
}

func (it *DefaultIniConfig) init() error {

	iniFile, _ := goini.LoadFile(INI_CONFIG_FILE)
	it.iniFileValues = iniFile.Section("")

	err := env.OnConfigIniStart()

	return err
}
