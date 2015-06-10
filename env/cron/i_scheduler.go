package cron

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/ottemo/foundation/env"
	"time"
)

// execute is a go routine function for a task (calls as separate go routine for each task)
func (it *DefaultCronScheduler) execute(schedule StructCronSchedule) {
	nextTime := schedule.expr.Next()

	c := time.Tick(time.Now().Sub(nextTime))
	_ = <-c

	schedule.task(schedule.Params)
	it.execute(schedule)
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

// ScheduleTask schedules task execution with a given params
func (it *DefaultCronScheduler) ScheduleTask(cronExpr string, taskName string, params map[string]interface{}) error {

	expr, err := cronexpr.Parse(cronExpr)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	task, present := it.tasks[taskName]
	if !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee521c4f-b84c-4238-bdac-ce61a37267a3", "unexistent task")
	}

	it.schedules = append(it.schedules, StructCronSchedule{
		CronExpr: cronExpr,
		TaskName: taskName,
		Params:   params,
		task:     task,
		expr:     expr})

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

	delete(it.schedules, idx)

	return nil
}
