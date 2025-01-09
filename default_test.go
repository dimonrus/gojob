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

func TestRun(t *testing.T) {
	SetRepeatDuration(time.Second)
	Add("test.hello.job", "- * * * * * * * *", func(ctx context.Context, args ...any) error {
		panic("internal error")
		return nil
	})
	Add("test.goodbye.job", "- * - - - - - - -", func(ctx context.Context, args ...any) error {
		return nil
	}, NewCondition(OperatorAND, func() bool {
		return true
	}))
	Run(log.Default(), LogMiddleware, RecoverMiddleware)
	time.Sleep(time.Second * 5)
}

func TestSetRepeatDuration(t *testing.T) {
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
