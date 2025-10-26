package jobs

type taskType string

const (
	llm taskType = "LLM"
	tool taskType = "Tool"
)

func (t taskType) isValidTask() bool {
	switch t {
	case llm, tool:
		return true
	default:
		return false
	}
}

type Task struct {
	T      taskType `json:"task_type"`
	Query  string `json:"query,omitempty"`
	Code   string `json:"code,omitempty"`
}

func NewTask(t string) *Task {
	tt := taskType(t)

	if !tt.isValidTask() {
		return nil
	}

	return &Task{
		T: tt,
	}
}

