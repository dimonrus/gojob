package gojob

import (
	"context"
	"errors"
	"log"
	"runtime"
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

func TestMode(t *testing.T) {
	t.Run("parallel", func(t *testing.T) {
		g := NewGroup(time.Millisecond*500, GroupModeAllParallel)
		ctx := context.WithValue(context.Background(), "logger", log.Default())
		ctx, cancel := context.WithTimeout(ctx, time.Second*6)
		defer cancel()
		g.SetRepeatDuration(time.Second)

		job1 := NewJob("job.parallel.1", func(ctx context.Context, args ...any) error {
			time.Sleep(time.Second * 2)
			t.Log("job.parallel.1 message")
			panic("error")
			return nil
		}, time.Second)

		job2 := NewJob("job.parallel.2", func(ctx context.Context, args ...any) error {
			time.Sleep(time.Second * 2)
			t.Log("job.parallel.2 message")
			return nil
		}, time.Second)

		job3 := NewJob("job.parallel.3", func(ctx context.Context, args ...any) error {
			time.Sleep(time.Second * 2)
			t.Log("job.parallel.3 message")
			return errors.New("some error")
		}, time.Second)

		g.AddJob(job1, job2, job3)
		g.Schedule(ctx, RecoverMiddleware)
	})
	t.Run("parallel", func(t *testing.T) {
		g := NewGroup(time.Millisecond*500, GroupModeConsistently)
		ctx := context.WithValue(context.Background(), "logger", log.Default())
		ctx, cancel := context.WithTimeout(ctx, time.Second*6)
		defer cancel()
		g.SetRepeatDuration(time.Second)

		job1 := NewJob("job.parallel.1", func(ctx context.Context, args ...any) error {
			time.Sleep(time.Second * 2)
			t.Log("job.parallel.1 message")
			panic("error")
			return nil
		}, time.Second)

		job2 := NewJob("job.parallel.2", func(ctx context.Context, args ...any) error {
			time.Sleep(time.Second * 2)
			t.Log("job.parallel.2 message")
			return nil
		}, time.Second)

		job3 := NewJob("job.parallel.3", func(ctx context.Context, args ...any) error {
			time.Sleep(time.Second * 2)
			t.Log("job.parallel.3 message")
			return errors.New("some error")
		}, time.Second)

		g.AddJob(job1, job2, job3)
		g.Schedule(ctx, RecoverMiddleware)
	})
	t.Run("parallel_N_error", func(t *testing.T) {
		grp := time.Millisecond * 10
		jrp := time.Millisecond * 500
		jp := time.Millisecond * 1000
		timeout := time.Second * 5
		g := NewGroup(grp, 3)
		ctx := context.WithValue(context.Background(), "logger", log.Default())
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		job1 := NewJob("job.parallel.1", func(ctx context.Context, args ...any) error {
			time.Sleep(jp)
			return errors.New("some error")
		}, jrp)

		g.AddJob(job1)
		g.Schedule(ctx)
	})

	t.Run("parallel_N", func(t *testing.T) {
		grp := time.Millisecond * 10
		jrp := time.Millisecond * 5
		jp := time.Millisecond * 10
		timeout := time.Second * 30

		printMemStat(t)

		g := NewGroup(grp, 3)
		ctx := context.WithValue(context.Background(), "logger", log.Default())
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		job1 := NewJob("job.parallel.1", func(ctx context.Context, args ...any) error {
			time.Sleep(jp)
			return nil
		}, jrp)

		job2 := NewJob("job.parallel.2", func(ctx context.Context, args ...any) error {
			time.Sleep(jp)
			return nil
		}, jrp)

		job3 := NewJob("job.parallel.3", func(ctx context.Context, args ...any) error {
			time.Sleep(jp)
			return nil
		}, jrp)

		job4 := NewJob("job.parallel.4", func(ctx context.Context, args ...any) error {
			time.Sleep(jp)
			return nil
		}, jrp)

		g.AddJob(job1, job2, job3, job4)

		printMemStat(t)

		g.Schedule(ctx)

		printMemStat(t)
	})
}

func printMemStat(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("\tAlloc = %v KB", m.Alloc/1024)
	t.Logf("\tTotalAlloc = %v KB", m.TotalAlloc/1024)
	t.Logf("\tSys = %v KB", m.Sys/1024)
	t.Logf("\tNumGC = %v\n", m.NumGC)
}
