package cron

import (
	"github.com/gorhill/cronexpr"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

func (it *DefaultCronSchedule) Execute() {
	currentTime := time.Now()

	if it.Time.Before(currentTime) && it.expr != nil {
		it.Time = it.expr.Next(currentTime)
	}

	if it.scheduler.appStarted {
		if currentTime.Before(it.Time) {
			c := time.Tick(it.Time.Sub(currentTime))
			_ = <-c
		}

		err := it.task(it.Params)
		if err != nil {
			err = env.ErrorDispatch(err)
			env.Log("cron.log", env.ConstLogPrefixError, err.Error())
		}

		if it.Repeat {
			it.Execute()
		} else {

		}
	} else {
		it.Execute()
	}
}

func (it *DefaultCronSchedule) Enable() error {
	found := false
	for _, item := range it.scheduler.schedules {
		if item == it {
			found = true
			break
		}
	}
	if !found {
		it.scheduler.schedules = append(it.scheduler.schedules, it)
	}

	return nil
}

func (it *DefaultCronSchedule) Disable() error {
	var newList []*DefaultCronSchedule
	for _, item := range it.scheduler.schedules {
		if item != it {
			newList = append(newList, item)
		}
	}
	it.scheduler.schedules = newList

	return nil
}

func (it *DefaultCronSchedule) Set(param string, value interface{}) error {
	switch param {
	case "expr":
		stringValue := utils.InterfaceToString(value)
		expr, err := cronexpr.Parse(stringValue)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		it.expr = expr
		it.CronExpr = stringValue

	case "time":
		it.Time = utils.InterfaceToTime(value)

	case "task":
		taskName := utils.InterfaceToString(value)
		if utils.IsInListStr(taskName, it.scheduler.ListTasks()) {
			it.TaskName = taskName
			it.task = it.scheduler.tasks[taskName]
		}

	case "repeat":
		it.Repeat = utils.InterfaceToBool(value)

	case "params":
		it.Params = utils.InterfaceToMap(value)
	}
	return nil
}

func (it *DefaultCronSchedule) Get(param string) interface{} {
	switch param {
	case "expr":
		return it.CronExpr

	case "time":
		return it.Time

	case "task":
		return it.TaskName

	case "repeat":
		return it.Repeat

	case "params":
		return it.Params
	}
	return nil
}

func (it *DefaultCronSchedule) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"expr":   it.CronExpr,
		"time":   it.Time,
		"task":   it.TaskName,
		"repeat": it.Repeat,
		"params": it.Params}
}
