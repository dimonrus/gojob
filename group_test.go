package gojob

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

func TestGroup_Schedule(t *testing.T) {
	g := NewGroup(time.Millisecond*500, GroupModeConsistently)

	repeatPeriod := time.Second * 5
	job1 := NewJob("test.group.job.1", func(ctx context.Context, args ...any) error {
		t.Log("run test.group.job.1: " + time.Now().Format(time.DateTime))
		return nil
	}, repeatPeriod)
	job1.SetSortOrder(10)

	job2 := NewJob("test.group.job.2", func(ctx context.Context, args ...any) error {
		t.Log("run test.group.job.2: " + time.Now().Format(time.DateTime))
		return errors.New("test error required")
	}, repeatPeriod*2)
	job2.SetSortOrder(1)

	job3 := NewJob("test.group.job.3", func(ctx context.Context, args ...any) error {
		t.Log("run test.group.job.3: " + time.Now().Format(time.DateTime))
		return nil
	}, time.Second*15)
	job3.SetSortOrder(2)

	job3.SetCondition(NewCondition(OperatorAND, func() bool {
		return time.Now().Weekday() == 4
	}))

	g.AddJob(job1, job2, job3)
	g.SortJobs()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*16)
	defer cancel()
	ctx = context.WithValue(ctx, "logger", log.Default())
	g.Schedule(ctx)

	g.ResetJobs()
	if len(g.jobs) > 0 {
		t.Fatal("Jobs must cleaned")
	}
}

func TestGroup_Panic(t *testing.T) {
	t.Run("panic_1", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("must be panic 1")
			}
		}()
		func() {
			g := NewGroup(time.Millisecond*500, GroupModeConsistently)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*16)
			defer cancel()
			ctx = context.WithValue(ctx, "logger", "aaa")
			g.SetRepeatDuration(time.Second)
			g.Schedule(ctx)
		}()
	})
	t.Run("panic_2", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("must be panic 2")
			}
		}()
		func() {
			g := NewGroup(time.Millisecond*500, GroupModeConsistently)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*16)
			defer cancel()
			g.SetRepeatDuration(time.Second)
			g.Schedule(ctx)
		}()
	})
}
