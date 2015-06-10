package cron

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultCronScheduler)
	var _ env.InterfaceScheduler = instance

	instance.tasks = make(map[string]env.FuncCronTask)
	instance.schedules = make([]StructCronSchedule, 0)

	app.OnAppInit(instance.appInitEvent)
	app.OnAppEnd(instance.appEndEvent)

	env.RegisterScheduler(instance)
}

// routines before application end
func (it *DefaultCronScheduler) appEndEvent() error {
	return nil
}

// routines before application start (on init phase)
func (it *DefaultCronScheduler) appInitEvent() error {

	// TODO: load manually specified tasks from DB

	for _, schedule := range it.schedules {
		go it.execute(schedule)
	}

	return nil
}
