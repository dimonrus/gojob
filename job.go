package gojob

import (
	"time"
)

// Jobs list of job
type Jobs []*Job

// JobCallback Main job callback
type JobCallback func(args ...any) error

// Job Simple executable schedule job
type Job struct {
	// Condition for the job
	condition Condition
	// Job callback
	callback JobCallback
	// When job must be scheduled
	nextAttemptAt time.Time
	// Pause before next run
	repeatPeriod time.Duration
	// Name of job
	name string
	// Sort order
	sortOrder uint
}

// isNextTime is next Time
func (j *Job) isNextTime(t time.Time) bool {
	if j.nextAttemptAt.IsZero() {
		return true
	}
	return j.nextAttemptAt.Before(t)
}

// SetCondition set job condition
func (j *Job) SetCondition(c Condition) {
	j.condition = c
	return
}

// SetSortOrder set sort order
func (j *Job) SetSortOrder(order uint) {
	j.sortOrder = order
	return
}

// GetSortOrder get sort order
func (j *Job) GetSortOrder() uint {
	return j.sortOrder
}

// CanStartAt is possible to start job now
func (j *Job) CanStartAt(t time.Time) bool {
	if j.condition.IsEmpty() {
		return j.isNextTime(t)
	}
	return j.isNextTime(t) && j.condition.IsTrue()
}

// GetName get job name
func (j *Job) GetName() string {
	return j.name
}

// SetNextTime set next Time
func (j *Job) SetNextTime(t time.Time) {
	j.nextAttemptAt = t
}

// Run job with params
func (j *Job) Run(arg ...any) error {
	return j.callback(arg...)
}

// RunAt run at specific time
func (j *Job) RunAt(t time.Time, arg ...any) (err error) {
	if !j.CanStartAt(t) {
		return
	}
	err = j.Run(arg...)
	j.SetNextTime(t.Add(j.repeatPeriod))
	return
}

// NewJob create new job
// callback - the job
// rp - repeat period
func NewJob(name string, callback JobCallback, rp time.Duration) *Job {
	return &Job{
		condition:    Condition{},
		callback:     callback,
		repeatPeriod: rp,
		name:         name,
	}
}
