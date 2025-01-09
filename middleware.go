package gojob

import (
	"context"
	"errors"
	"fmt"
)

// Middleware function
type Middleware func(job *Job, callback JobCallback) JobCallback

// RecoverMiddleware specify how to app must catch common panic exceptions
func RecoverMiddleware(job *Job, callback JobCallback) JobCallback {
	return JobCallback(func(ctx context.Context, args ...any) error {
		defer func() {
			var err error
			if r := recover(); r != nil {
				switch r.(type) {
				case error:
					err = errors.New("panic: " + r.(error).Error())
				case string:
					err = errors.New("panic: " + r.(string))
				default:
					err = errors.New(fmt.Sprintf("panic: %T", r))
				}
			}
			if err != nil {
				ctx.Value("logger").(Logger).Printf("\x1b[31;1mjob: %s %s \x1b[31;1m", job.GetName(), err)
			}
		}()
		return callback(ctx, args...)
	})
}

// LogMiddleware log of executed middleware
func LogMiddleware(job *Job, callback JobCallback) JobCallback {
	return JobCallback(func(ctx context.Context, args ...any) error {
		ctx.Value("logger").(Logger).Printf("job: %s scheduled", job.GetName())
		return callback(ctx, args...)
	})
}
