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
	for i, job := range s.Jobs.Jobs {
		rs[len(s.Jobs.Jobs)-1-i] = job
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
	for i, task := range s.Tasks {
		rs[i] = task
	}
	return rs
}

type JobInstanceWrapper struct {
	JobInstances *sodor.JobInstances
}

func NewJobInstanceWrapper(ins *sodor.JobInstances) *JobInstanceWrapper {
	return &JobInstanceWrapper{
		JobInstances: ins,
	}
}

func (s *JobInstanceWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.JobInstances.JobInstances), len(s.JobInstances.JobInstances))
	for i, ins := range s.JobInstances.JobInstances {
		rs[len(s.JobInstances.JobInstances)-1-i] = ins
	}
	return rs
}

type TaskInstanceWrapper struct {
	TaskInstances *sodor.TaskInstances
}

func NewTaskInstanceWrapper(ins *sodor.TaskInstances) *TaskInstanceWrapper {
	return &TaskInstanceWrapper{
		TaskInstances: ins,
	}
}

func (s *TaskInstanceWrapper) AsInterfaceArray(taskIDs ...int32) []interface{} {
	rs := make([]interface{}, 0, len(s.TaskInstances.TaskInstances))

	for i := len(s.TaskInstances.TaskInstances) - 1; i >= 0; i-- {
		for _, id := range taskIDs {
			if s.TaskInstances.TaskInstances[i].TaskId == id {
				rs = append(rs, s.TaskInstances.TaskInstances[i])
				break
			}
		}
	}

	return rs
}
