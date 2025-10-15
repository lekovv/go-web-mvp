package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/lekovv/go-web-mvp/service"
)

type Job interface {
	Start(ctx context.Context)
}

type JobManager struct {
	jobs []Job
	wg   sync.WaitGroup
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make([]Job, 0),
	}
}

func (jm *JobManager) AddJob(job Job) {
	jm.jobs = append(jm.jobs, job)
}

func (jm *JobManager) StartAll(ctx context.Context) {
	for _, job := range jm.jobs {
		jm.wg.Add(1)
		go func(j Job) {
			defer jm.wg.Done()
			j.Start(ctx)
		}(job)
	}
}

func (jm *JobManager) Wait() {
	jm.wg.Wait()
}

func CreateJobs(authService service.AuthServiceInterface) *JobManager {
	manager := NewJobManager()

	cleanupJob := NewTokenBlacklistCleaner(authService, 24*time.Hour)
	manager.AddJob(cleanupJob)

	return manager
}
