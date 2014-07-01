package app

import (
	"github.com/ottemo/foundation/api"
)

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


func Serve() error {
	return api.GetRestService().Run()
}
