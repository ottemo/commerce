package swatch

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/media"
)

// init makes package self-initialization routine
func init() {
	app.OnAppStart(onAppStart)
	api.RegisterOnRestServiceStart(setupAPI)
}

func onAppStart() error {

	var err error
	mediaStorage, err = media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
