package session

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"io/ioutil"
	"os"
)

// init makes package self-initialization routine
func init() {
	sessionService = new(DefaultSessionService)
	sessionService.gcRate = 10
	sessionService.Sessions = make(map[string]*DefaultSession)

	var _ api.InterfaceSessionService = sessionService

	api.RegisterSessionService(sessionService)

	app.OnAppStart(startup)
	app.OnAppEnd(shutdown)
}

// service startup routines
func startup() error {

	// checking session storage directory exists
	if _, err := os.Stat(ConstStorageFolder); !os.IsExist(err) {
		err := os.MkdirAll(ConstStorageFolder, os.ModePerm)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	// listing files within session storage folder
	files, err := ioutil.ReadDir(ConstStorageFolder)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// loading session data
	for _, fileInfo := range files {
		sessionInstance := new(DefaultSession)
		err := sessionInstance.Load(fileInfo.Name())
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// service shutdown routines
func shutdown() error {

	// saving all session to storage
	for _, session := range sessionService.Sessions {
		err := session.Save()
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}
