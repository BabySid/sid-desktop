package common

import (
	"fmt"
	"github.com/BabySid/proto/sodor"
	"github.com/sahilm/fuzzy"
)

type JobsWrapper struct {
	Jobs *sodor.Jobs
}

func NewJobsWrapper(jobs *sodor.Jobs) *JobsWrapper {
	return &JobsWrapper{
		Jobs: jobs,
	}
}

func (s *JobsWrapper) Find(name string) *JobsWrapper {
	matches := fuzzy.FindFrom(name, s)

	rs := sodor.Jobs{}
	rs.Jobs = make([]*sodor.Job, 0)
	for _, match := range matches {
		rs.Jobs = append(rs.Jobs, s.Jobs.Jobs[match.Index])
	}

	return NewJobsWrapper(&rs)
}

func (s *JobsWrapper) String(i int) string {
	return fmt.Sprintf("%s", s.Jobs.Jobs[i].Name)
}

func (s *JobsWrapper) Len() int {
	return len(s.Jobs.Jobs)
}

func (s *JobsWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.Jobs.Jobs), len(s.Jobs.Jobs))
	for i := range s.Jobs.Jobs {
		rs[len(s.Jobs.Jobs)-1-i] = s.Jobs.Jobs[i]
	}
	return rs
}

type JobTasksWrapper struct {
	Tasks []*sodor.Task
}

func NewJobTasksWrapper(job *sodor.Job) *JobTasksWrapper {
	return &JobTasksWrapper{
		Tasks: job.Tasks,
	}
}

func (s *JobTasksWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.Tasks), len(s.Tasks))
	for i := range s.Tasks {
		rs[i] = s.Tasks[i]
	}
	return rs
}
