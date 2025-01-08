package gojob

import (
	"testing"
	"time"
)

func TestSetRepeatDuration(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		SetRepeatDuration(time.Second)
		err := AddJob("test.hello.job", "- * * * * * * * *", func(args ...any) error {
			t.Log("hello. This is test.hello.job")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		err = AddJob("test.goodbye.job", "- * - - - - - - -", func(args ...any) error {
			t.Log("goodbye. This is test.goodbye.job")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		Run()
		time.Sleep(time.Second * 10)
	})
}
