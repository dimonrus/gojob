package gojob

import (
	"context"
	"time"
)

var (
	group      = initDefaultGroup()
	defaultCtx = context.Background()
)

// SetRepeatDuration change repeat duration
// Must be called before Add methods for sync repeatDuration between jobs
func SetRepeatDuration(d time.Duration) {
	group.SetRepeatDuration(d)
}

// Add job to default schedule
func Add(name string, expression ScheduleExpression, callback JobCallback, condition ...Condition) (*Job, error) {
	job := NewJob(name, callback, group.d)
	tp, err := expression.Parse()
	if err != nil {
		return nil, err
	}
	cond := tp.ToCondition()
	if len(condition) > 0 {
		cond = cond.Merge(OperatorAND, condition...)
	}
	job.SetCondition(cond)
	group.AddJob(job)
	return job, nil
}

// Run default schedule group
func Run(logger Logger, middleware ...Middleware) {
	defaultCtx.Done()
	defaultCtx = context.Background()
	defaultCtx = context.WithValue(defaultCtx, "logger", logger)
	go group.Schedule(defaultCtx, middleware...)
}

// Init default schedule group
func initDefaultGroup() *Group {
	return NewGroup(time.Minute, GroupModeConsistently)
}
