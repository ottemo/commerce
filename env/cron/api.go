package cron

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cron/tasks", api.ConstRESTOperationGet, getSchedules)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cron/functions", api.ConstRESTOperationGet, getTasks)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cron/task/:taskIndex", api.ConstRESTOperationUpdate, updateSchedule)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cron/task", api.ConstRESTOperationCreate, createSchedule)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// getSchedules to get information about current schedules
func getSchedules(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var result []interface {}

	scheduler := env.GetScheduler()

	for _, schedule := range scheduler.ListSchedules() {
		result = append(result, schedule.GetInfo())
	}

	return result, nil
}

// getTasks return scheduler registered tasks (functions that are available to execute)
func getTasks(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	scheduler := env.GetScheduler()

	return scheduler.ListTasks(), nil
}

// updateSchedule update scheduler task
//   - "taskIndex" should be specified as argument (task index can be obtained from getSchedules)
func updateSchedule(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reqTaskIndex := context.GetRequestArgument("taskIndex")
	if reqTaskIndex == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d4ee4c0c-124a-4098-aeef-23d868b0d682", "task index should be specified")
	}

	taskIndex, err := utils.StringToInteger(reqTaskIndex)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}


	scheduler := env.GetScheduler()
	scheduleParams := []string {"expr", "time", "task", "repeat", "params"}

	for index, schedule := range scheduler.ListSchedules() {
		if index == taskIndex {
			for _, param := range scheduleParams {
				if value, present := postValues[param]; present {
					err = schedule.Set(param, value)
					if err != nil {
						return nil, env.ErrorDispatch(err)
					}
				}
			}
		}
	}

	return "ok", nil
}

// createSchedule with request params
// in request params required are time or cronExpr for creating different type of schedules
func createSchedule(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cronExpression := utils.InterfaceToString(utils.GetFirstMapValue(postValues, "CronExpr", "cronExpr", "Expr", "expr"))
	scheduledTime := utils.InterfaceToTime(utils.GetFirstMapValue(postValues, "Time", "time"))

	if utils.IsZeroTime(scheduledTime) && cronExpression == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2250e0a5-e439-444d-b331-13d16f8e2401", "cronExpr or time were not specified")
	}

	scheduler := env.GetScheduler()

	taskName := utils.InterfaceToString(utils.GetFirstMapValue(postValues, "TaskName", "Task", "task"))

	if taskName == "" || !utils.IsInListStr (taskName, scheduler.ListTasks()) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "eafc6b15-e897-4d9f-a93a-f84cffa78497", "task not specified or not regisetered")
	}

	isRepeat := utils.InterfaceToBool(postValues["repeat"])
	taskParams := utils.InterfaceToMap(postValues["params"])

	var newSchedule env.InterfaceSchedule
	if !utils.IsZeroTime(scheduledTime) {
		newSchedule, err = scheduler.ScheduleAtTime(scheduledTime, taskName, taskParams)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	} else {
		if isRepeat {
			newSchedule, err = scheduler.ScheduleRepeat(cronExpression, taskName, taskParams)
		} else {
			newSchedule, err = scheduler.ScheduleOnce(cronExpression, taskName, taskParams)
		}

		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	return newSchedule, nil
}
