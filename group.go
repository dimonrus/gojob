package gojob

import (
	"context"
	"sort"
	"time"
)

// GroupMode type of parallel mode
type GroupMode int

const (
	// GroupModeAllParallel All jobs wil be executed in Parallel mode
	GroupModeAllParallel GroupMode = -1
	// GroupModeConsistently All jobs wil be executed in Consistently mode
	GroupModeConsistently GroupMode = 0
)

// Group of periodical schedule jobs
type Group struct {
	// Period of next job run
	d time.Duration
	// List of jobs for run
	jobs Jobs
	// how many parallel jobs can run
	// parallel == -1 - all jobs will run in parallel mode
	// parallel == 0 - no jobs will run in parallel mode
	// parallel == N - specify the number N of jobs that can run in parallel mode
	parallel GroupMode
}

// Schedule run periodical scheduler
func (g *Group) Schedule(ctx context.Context, middlewares ...Middleware) {
	ticker := time.NewTicker(g.d)
	defer ticker.Stop()
	l := ctx.Value("logger")
	if l == nil {
		panic("logger is not defined")
	}
	logger, ok := l.(Logger)
	if !ok {
		panic("logger must implement Logger interface")
	}
	for i := range g.jobs {
		for _, middleware := range middlewares {
			g.jobs[i].callback = middleware(g.jobs[i], g.jobs[i].callback)
		}
	}
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for _, job := range g.jobs {
				if g.parallel == GroupModeAllParallel {
					go func(j *Job, x context.Context, l Logger, t time.Time) {
						e := j.RunAt(x, t)
						if e != nil {
							l.Println(e.Error())
						}
					}(job, ctx, logger, now)
				} else if g.parallel == GroupModeConsistently {
					e := job.RunAt(ctx, now)
					if e != nil {
						logger.Println(e.Error())
					}
				}

			}
		case <-ctx.Done():
			return
		}
	}
}

// Add add jobs
func (g *Group) AddJob(job ...*Job) *Group {
	g.jobs = append(g.jobs, job...)
	return g
}

// SetRepeatDuration set repeat duration
func (g *Group) SetRepeatDuration(d time.Duration) *Group {
	g.d = d
	return g
}

// ResetJobs reset jobs
func (g *Group) ResetJobs() *Group {
	g.jobs = g.jobs[:0]
	return g
}

// SortJobs reset jobs
func (g *Group) SortJobs() *Group {
	sort.Slice(g.jobs, func(i, j int) bool {
		return g.jobs[i].GetSortOrder() < g.jobs[j].GetSortOrder()
	})
	return g
}

// NewGroup Create new group
// d - repeat duration for next time run
// mode - parallel mode (-1, 0, N)
func NewGroup(d time.Duration, mode GroupMode) *Group {
	return &Group{
		d:        d,
		parallel: mode,
	}
}
