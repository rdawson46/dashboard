package jobs

type taskType int

const (
	llm taskType = iota
	tool
)


type Task struct {
	T      taskType `json:"task_type"`
	Query  string `json:"query"`
	Params map[string]string `json:"parameters,omitempty"`
}

type Job struct {
	Result []string `json:"result"`
	Time   string `json:"time,omitempty"`
	Freq   string `json:"freq,omitempty"`
	Tasks  []Task `json:"tasks"` 
	Model  string `json:"model"`
	Id	   string `json:"id"`
	Status string `json:"status"`
}

type Jobs []Job

func handleLlmTask(task Task) string {
	return ""
}

func handleToolTask(task Task) string {
	return ""
}


func (j *Job) Run() {
	results := make([]string, 0)

	for _, t := range j.Tasks {
		var r string
		switch t.T {
		case llm:
			r = handleLlmTask(t)
		case tool:
			r = handleLlmTask(t)
		default:
		}

		results = append(results, r)
	}

	j.Result = results
}
