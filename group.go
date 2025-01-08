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
func (g *Group) Schedule(logger Logger, ctx context.Context) {
	ticker := time.NewTicker(g.d)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for _, job := range g.jobs {
				e := job.RunAt(now)
				if e != nil {
					logger.Println(e.Error())
				}
			}
		case <-ctx.Done():
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
