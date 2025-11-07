package jobs

import (
	"errors"
	"net/url"
	"strings"
)

type Status string

const (
	StatusRunning Status = "Running"
	StatusPending Status = "Pending"
	StatusCompleted Status = "Completed"
	StatusFailed Status = "Failed"
)


func (s Status) isValidStatus() bool {
	switch s {
	case StatusRunning, StatusPending, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

type Job struct {
	Name   string `json:"name"`
	Result string `json:"result"`
	Time   string `json:"time,omitempty"`
	Freq   string `json:"freq,omitempty"`
	// Tasks  []Task `json:"tasks"` 
	Task   *Task `json:"task"`
	Model  string `json:"model"`
	Id	   string `json:"id,omitempty"`
	Status Status `json:"status"`
}

type Jobs []*Job


func NewJob(task *Task, name, freq string) *Job {
	return &Job{
		Name: name,
		Freq: freq,
		Status: StatusPending,
		Task: task,
	}
}

func handleLlmTask(task *Task) string {
	return ""
}

func handleToolTask(task *Task) string {
	return ""
}


func (j *Job) Run() {
	var r string
	switch j.Task.T {
	case llm:
		r = handleLlmTask(j.Task)
	case tool:
		r = handleToolTask(j.Task)
	default:
	}

	j.Result = r
}

func (j *Job) FillIn(form url.Values) error {
	switch j.Task.T {
	case llm:
		query := strings.TrimSpace(form.Get("query"))

		if query == "" {
			return errors.New("Invalid query")
		}

		j.Task.Query = query
		return nil
	case tool:
		return errors.New("Not implemented")
	default:
		return errors.New("Invalid Job Type")
	}
}
