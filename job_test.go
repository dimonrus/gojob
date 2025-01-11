package gojob

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRunJob(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		repeatPeriod := time.Second * 5
		job := NewJob("test.job", func(ctx context.Context, args ...any) error {
			if len(args) > 0 && args[0].(int) > 0 {
				t.Log("job scheduled with number")
			} else {
				t.Log("job scheduled without params")
			}
			return nil
		}, repeatPeriod)
		if job.GetRepeatPeriod() != repeatPeriod {
			t.Fatal("repeat period wrong")
		}
		job.SetRepeatPeriod(repeatPeriod)
		e := job.RunAt(context.Background(), time.Now(), 0)
		if e != nil {
			t.Fatal(e)
		}
		if !job.CanStartAt(time.Now()) {
			t.Log(job.GetName() + " can't start until repeat period will be reached")
		}
		// it wil not scheduled until repeat period duration
		t.Logf("wait for %v + 1", repeatPeriod)
		time.Sleep(repeatPeriod + time.Second)
		e = job.RunAt(context.Background(), time.Now(), 1)
		if e != nil {
			t.Fatal(e)
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gojob
// cpu: Apple M2 Max
// BenchmarkJob_RunAt
// BenchmarkJob_RunAt-12    	33538953	        33.54 ns/op	       0 B/op	       0 allocs/op
func BenchmarkJob_RunAt(b *testing.B) {
	job := NewJob("test.job", func(ctx context.Context, args ...any) error {
		return nil
	}, time.Millisecond)
	g := NewGroup(time.Microsecond, GroupModeConsistently)
	g.AddJob(job)
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		now := time.Now()
		err := job.RunAt(ctx, now)
		_ = err
	}
	b.ReportAllocs()
}

func BenchmarkGroup_Schedule(b *testing.B) {
	job := NewJob("test.job", func(ctx context.Context, args ...any) error {
		return nil
	}, time.Millisecond)
	g := NewGroup(time.Microsecond, 3)
	g.AddJob(job)
	ctx := context.WithValue(context.Background(), "logger", log.Default())
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	for i := 0; i < 1; i++ {
		g.Schedule(ctx)
	}
	b.ReportAllocs()
}
