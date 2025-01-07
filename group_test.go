package gojob

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestGroup_Schedule(t *testing.T) {
	g := NewGroup(time.Millisecond*500, GroupModeConsistently)

	repeatPeriod := time.Second * 5
	job1 := NewJob("test.group.job.1", func(args ...any) error {
		t.Log("run test.group.job.1: " + time.Now().Format(time.DateTime))
		return nil
	}, repeatPeriod)
	job1.SetSortOrder(10)

	job2 := NewJob("test.group.job.2", func(args ...any) error {
		t.Log("run test.group.job.2: " + time.Now().Format(time.DateTime))
		return nil
	}, repeatPeriod*2)
	job2.SetSortOrder(1)

	job3 := NewJob("test.group.job.3", func(args ...any) error {
		t.Log("run test.group.job.3: " + time.Now().Format(time.DateTime))
		return nil
	}, time.Second*15)
	job3.SetSortOrder(2)

	job3.SetCondition(NewCondition(OperatorAND, func() bool {
		return time.Now().Weekday() == 4
	}))

	g.AddJob(job1, job2, job3)
	g.SortJobs()

	logger := log.Default()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*16)
	g.Schedule(logger, ctx)
}
