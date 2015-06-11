package cron

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/ottemo/foundation/env"
	"time"
)

// execute is a go routine function for a task (calls as separate go routine for each task)
func (it *DefaultCronScheduler) execute(schedule StructCronSchedule) {

	currentTime := time.Now()

	if schedule.Time.Before(currentTime) && schedule.expr != nil {
		schedule.Time = schedule.expr.Next(currentTime)
	}
	nextTime := schedule.expr.Next(schedule.Time)

	if it.appStarted {
		if currentTime.Before(schedule.Time) {
			c := time.Tick(nextTime.Sub(currentTime))
			_ = <-c
		}

		schedule.task(schedule.Params)

		if schedule.Repeat {
			it.execute(schedule)
		}
	} else {
		it.execute(schedule)
	}
}

// ListTasks returns a list of task names currently available
func (it *DefaultCronScheduler) ListTasks() []string {
	var result []string
	for taskName := range it.tasks {
		result = append(result, taskName)
	}
	return result
}

// RegisterTask registers a new task routine by a given task name
//   - returns error no non unique name
func (it *DefaultCronScheduler) RegisterTask(name string, task env.FuncCronTask) error {
	if _, present := it.tasks[name]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "231fa82d-c357-498d-b0b3-f4daee7e25c5", "task already exists")
	}

	it.tasks[name] = task

	return nil
}

// ScheduleOnce schedules task execution with a given params
func (it *DefaultCronScheduler) ScheduleAtTime(scheduleTime time.Time, taskName string, params map[string]interface{}) error {
	expr, err := cronexpr.Parse("*/1 * * * *")
	if err != nil {
		return env.ErrorDispatch(err)
	}

	task, present := it.tasks[taskName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee521c4f-b84c-4238-bdac-ce61a37267a3", "unexistent task")
	}

	schedule := StructCronSchedule{
		CronExpr: "*/1 * * * *",
		TaskName: taskName,
		Params:   params,
		Repeat:   false,
		Time: 	  scheduleTime,
		task:     task,
		expr:     expr}

	it.schedules = append(it.schedules, schedule)

	go it.execute(schedule)

	return nil
}

// ScheduleOnce schedules task execution with a given params
func (it *DefaultCronScheduler) ScheduleOnce(cronExpr string, taskName string, params map[string]interface{}) error {
	expr, err := cronexpr.Parse(cronExpr)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	task, present := it.tasks[taskName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee521c4f-b84c-4238-bdac-ce61a37267a3", "unexistent task")
	}

	schedule := StructCronSchedule{
		CronExpr: cronExpr,
		TaskName: taskName,
		Params:   params,
		Repeat:   false,
		task:     task,
		expr:     expr}

	it.schedules = append(it.schedules, schedule)

	go it.execute(schedule)

	return nil
}

// ScheduleRepeat schedules task execution with a given params
func (it *DefaultCronScheduler) ScheduleRepeat(cronExpr string, taskName string, params map[string]interface{}) error {

	expr, err := cronexpr.Parse(cronExpr)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	task, present := it.tasks[taskName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee521c4f-b84c-4238-bdac-ce61a37267a3", "unexistent task")
	}

	schedule := StructCronSchedule{
		CronExpr: cronExpr,
		TaskName: taskName,
		Params:   params,
		Repeat:   true,
		task:     task,
		expr:     expr}

	it.schedules = append(it.schedules, schedule)

	go it.execute(schedule)

	return nil
}

// ListSchedules returns list of currently registered schedules
func (it *DefaultCronScheduler) ListSchedules() []string {
	var result []string

	for idx, schedule := range it.schedules {
		resultLine := fmt.Sprintf("%d. %s - %s (%v)", idx, schedule.CronExpr, schedule.TaskName, schedule.Params)
		result = append(result, resultLine)
	}

	return result
}

// RemoveSchedule removes schedule for a given index, which could be obtained by ListSchedules() call
func (it *DefaultCronScheduler) RemoveSchedule(idx int) error {
	if idx >= len(it.schedules) {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee521c4f-b84c-4238-bdac-ce61a37267a3", "invalid index")
	}

	var result [] StructCronSchedule
	for i, item := range it.schedules {
		if i != idx {
			result = append(result, item)
		}
	}
	it.schedules = result

	return nil
}
