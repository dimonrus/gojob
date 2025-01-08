package gojob

import (
	"context"
	"log"
	"time"
)

var (
	group = initDefaultGroup()
	ctx   = context.Background()
)

// SetRepeatDuration change repeat duration
// Must be called before AddJob methods for sync repeatDuration between jobs
func SetRepeatDuration(d time.Duration) {
	group.d = d
}

// AddJob add job to default schedule
func AddJob(name string, expression ScheduleExpression, callback JobCallback) error {
	job := NewJob(name, callback, group.d)
	tp, err := expression.Parse()
	if err != nil {
		return err
	}
	job.SetCondition(tp.ToCondition())
	group.AddJob(job)
	return nil
}

// Run default schedule group
func Run() {
	ctx.Done()
	ctx = context.Background()
	runScheduler()
}

// Init default schedule group
func initDefaultGroup() *Group {
	return NewGroup(time.Minute, GroupModeConsistently)
}

// run schedule
func runScheduler() {
	go group.Schedule(log.Default(), ctx)
}
