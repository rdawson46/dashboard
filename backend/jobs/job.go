package jobs

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"


	api "github.com/rdawson46/dashboard/api"
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

type taskResultType struct {
	value string
	status Status
}

func handleLLMTask(task *Task) taskResultType {
	query := task.Query

	if query == "" {
		return taskResultType{
			value: "Empty query",
			status: StatusFailed,
		}
	}

	url := os.Getenv("OLLAMA_URL")
	
	if url == "" {
		return taskResultType{
			value: "Server Issue",
			status: StatusFailed,
		}
	}

    oc, err := api.NewOllamaClient(url, nil)

	if err != nil {
		return taskResultType{
			value: "Server error",
			status: StatusFailed,
		}
	}

	ctx := context.Background()

	res, err := oc.Chat(ctx, query)

	if err != nil {
		return taskResultType{
			value: "Server error",
			status: StatusFailed,
		}
	}

	return taskResultType{
		value: res,
		status: StatusCompleted,
	}
}

func handleToolTask(task *Task) taskResultType {
	return taskResultType{

	}
}

/*

  1. Update status to running
  2. run
  3. catch value or error
  4. save value
  5. update status to failed or done

*/

func (j *Job) Run() taskResultType {
	var r taskResultType
	switch j.Task.T {
	case llm:
		r = handleLLMTask(j.Task)
	case tool:
		r = handleToolTask(j.Task)
	default:
	}

	j.Result = r.value

	return r
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
