package app

import (
	"github.com/ottemo/foundation/api"
)

var callbacksOnAppInit = []func() error {}

func OnAppInit(callback func() error) {
	callbacksOnAppInit = append(callbacksOnAppInit, callback)
}

var callbacksOnAppStart = []func() error {}

func OnAppStart(callback func() error) {
	callbacksOnAppStart = append(callbacksOnAppStart, callback)
}


func Start() error {
	for _, callback := range callbacksOnAppStart {
		if err := callback(); err != nil {
			return err
		}
	}

	return nil
}


func Init() error {
	for _, callback := range callbacksOnAppInit {
		if err := callback(); err != nil {
			return err
		}
	}

	return nil
}


func Serve() error {
	return api.GetRestService().Run()
}
