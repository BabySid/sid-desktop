package base

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
	// name -> id
	jobs map[string]cron.EntryID
}

var GlobalScheduler *Scheduler

type ScheJob interface {
	Run()
}

func NewScheduler() *Scheduler {
	if GlobalScheduler != nil {
		return GlobalScheduler
	}

	GlobalScheduler = &Scheduler{
		c:    cron.New(cron.WithSeconds()),
		jobs: make(map[string]cron.EntryID),
	}
	return GlobalScheduler
}

func (s *Scheduler) AddJob(name string, spec string, cmd ScheJob) error {
	if _, ok := s.jobs[name]; ok {
		return fmt.Errorf("%s exist", name)
	}

	id, e := s.c.AddJob(spec, cmd)
	if e != nil {
		return e
	}

	s.jobs[name] = id
	return nil
}

func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) Stop() {
	s.c.Stop()
}
