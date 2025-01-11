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

type parallelData struct {
	logger Logger
	job    *Job
	t      time.Time
}

// apply middlewares to jobs
func (g *Group) applyMiddleware(middlewares ...Middleware) {
	for i := range g.jobs {
		for _, middleware := range middlewares {
			g.jobs[i].callback = middleware(g.jobs[i], g.jobs[i].callback)
		}
	}
}

// extract logger
func (g *Group) getLoggerFromContext(ctx context.Context) Logger {
	l := ctx.Value("logger")
	if l == nil {
		panic("logger is not defined")
	}
	logger, ok := l.(Logger)
	if !ok {
		panic("logger must implement Logger interface")
	}
	return logger
}

// Schedule run periodical scheduler
func (g *Group) Schedule(ctx context.Context, middlewares ...Middleware) {
	// init ticker with user repeat duration
	ticker := time.NewTicker(g.d)
	// when context done ticker must be stopped
	defer ticker.Stop()
	// apply middlewares to all job at current moment
	g.applyMiddleware(middlewares...)
	// extract logger from context
	logger := g.getLoggerFromContext(ctx)
	// init chan for case when N parallel job possible at the moment
	var dataChan chan parallelData
	// flag for stop all job
	var exitChan = make(chan struct{})
	// for parallel N define worker logic
	if g.parallel > 0 {
		// define parallelData chan with g.parallel length
		dataChan = make(chan parallelData, g.parallel)
		// sleep one iteration before close channel
		defer func() {
			time.Sleep(g.d)
			close(dataChan)
		}()
		// make goroutines for parallel processing
		for i := 0; i < int(g.parallel); i++ {
			go func(x context.Context, dc chan parallelData, q <-chan struct{}) {
				for {
					select {
					case d := <-dc:
						// run job at provided time
						e := d.job.RunAt(x, d.t)
						if e != nil {
							d.logger.Println(e.Error())
						}
					case <-q:
						// close goroutine when context is done
						return
					}

				}
			}(ctx, dataChan, exitChan)
		}
	}
	// common scheduler life cycle
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for i := range g.jobs {
				job := g.jobs[i]
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
				} else if g.parallel > 0 {
					dataChan <- parallelData{
						logger: logger,
						job:    job,
						t:      now,
					}
				}
			}
		case <-ctx.Done():
			// quit from child goroutines
			close(exitChan)
			return
		}
	}
}

// AddJob add jobs
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
