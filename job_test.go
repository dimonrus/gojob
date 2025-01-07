package gojob

import (
	"testing"
	"time"
)

func TestRunJob(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		repeatPeriod := time.Second * 5
		job := NewJob("test.job", func(args ...any) error {
			if len(args) > 0 && args[0].(int) > 0 {
				t.Log("job scheduled with number")
			} else {
				t.Log("job scheduled without params")
			}
			return nil
		}, repeatPeriod)
		e := job.RunAt(time.Now(), 0)
		if e != nil {
			t.Fatal(e)
		}
		if !job.CanStartAt(time.Now()) {
			t.Log(job.GetName() + " can't start until repeat period will be reached")
		}
		// it wil not scheduled until repeat period duration
		t.Logf("wait for %v + 1", repeatPeriod)
		time.Sleep(repeatPeriod + time.Second)
		e = job.RunAt(time.Now(), 1)
		if e != nil {
			t.Fatal(e)
		}
	})
}
