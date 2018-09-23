package services

import (
	"github.com/bamzi/jobrunner"
	"github.com/lytnin/lytnin"
)

// JobRunner service provides scheduler for the application
type JobRunner struct {
}

// Info returns information about the scheduler store
func (s *JobRunner) Info() interface{} {
	return "job runner"
}

// Init initializes the key/value service and registers it with the application
func (s *JobRunner) Init(a *lytnin.Application) {
	jobrunner.Start()
	a.AddService("jobrunner", s)
}

// Close releases any resources used by the service
func (s *JobRunner) Close() {
	jobrunner.Stop()
}
