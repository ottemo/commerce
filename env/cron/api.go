package cron

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"fmt"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cron/tasks", api.ConstRESTOperationGet, getTasks)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// getTasks to get value information about config items with type [ConstConfigTypeGroup]
func getTasks(context api.InterfaceApplicationContext) (interface{}, error) {

	scheduler := env.GetScheduler()

	scheduleParams := []string {"expr", "time", "task", "repeat", "params"}

	scheduler.ListSchedules()
	for _, schedule := range scheduler.ListSchedules() {
		fmt.Println(schedule.GetInfo())
		for _, param := range scheduleParams {
			fmt.Println(schedule.Get(param))
		}
	}

	for _, task := range scheduler.ListTasks() {
		fmt.Println(task)
	}

	return scheduler.ListTasks(), nil
}
