package vantagepoint

import (
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

func init() {
	initHoursList()

	db.RegisterOnDatabaseStart(onDatabaseStart)

	env.RegisterOnConfigStart(setupConfig)
}

func onDatabaseStart() error {
	app.OnAppStart(onAppStart)

	return nil
}

func onAppStart() error {
	if err := scheduleCheckNewUploads(); err != nil {
		_ = env.ErrorDispatch(err)
	}

	return nil
}
