package gojob

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

type cTError string

func (c cTError) Text() string {
	return string(c)
}

func TestSetRepeatDuration(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		SetRepeatDuration(time.Second)
		job1, _ := Add("test.hello.job", "- * * * * * * * *", func(ctx context.Context, args ...any) error {
			panic("internal error")
			return nil
		})
		job1.SetRepeatPeriod(time.Millisecond * 500)
		job2, _ := Add("test.goodbye.job", "- * - - - - - - -", func(ctx context.Context, args ...any) error {
			return nil
		}, NewCondition(OperatorAND, func() bool {
			return true
		}))
		job2.SetRepeatPeriod(time.Millisecond * 500)
		Run(log.Default(), LogMiddleware, RecoverMiddleware)
		time.Sleep(time.Second * 5)
	})
	t.Run("error_parsing", func(t *testing.T) {
		SetRepeatDuration(time.Second)
		_, err := Add("test.parsing.error", "- //* - - - - - - -", func(ctx context.Context, args ...any) error {
			return nil
		})
		if err == nil {
			t.Fatal("must be an error parsing")
		}
		_, err = Add("test.custom.error", "- * - - - - - - -", func(ctx context.Context, args ...any) error {
			panic(cTError("error"))
			return nil
		})
		_, err = Add("test.standard.error", "- * - - - - - - -", func(ctx context.Context, args ...any) error {
			panic(errors.New("some error"))
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		Run(log.Default(), LogMiddleware, RecoverMiddleware)
		time.Sleep(time.Second * 5)
	})
}
