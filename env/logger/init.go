package logger

import (
	"os"
	"fmt"

	"github.com/ottemo/commerce/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultLogger)
	var _ env.InterfaceLogger = instance

	if err := env.RegisterLogger(instance); err != nil {
		fmt.Println(err.Error())
	}
	env.RegisterOnConfigIniStart(startup)
	env.RegisterOnConfigStart(setupConfig)
}

// startup is a service pre-initialization stuff
func startup() error {
	if _, err := os.Stat(baseDirectory); !os.IsExist(err) {
		err := os.MkdirAll(baseDirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
