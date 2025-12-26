package server

import (
	"context"
	"fmt"
	"os"

	ollama "github.com/ollama/ollama/api"
	api "github.com/rdawson46/dashboard/api"
	"github.com/rdawson46/dashboard/db"
)

const LIMIT_BUFFER = 0.8

var FilePrompt = `
The user has attached files for you to use in your response - they are attached below.
Accurately leverage the file content to answer their query and do not result to your own knowledge if the information is preseneted to you.


# Attached Files

%s

# User Query

%s
`

// rough estimation
func getTokenCount(x any) int {
	switch i := x.(type) {
	case string:
		return len(i) / 4

	case ollama.Message:
		return getTokenCount(i.Content)

	case []ollama.Message:
		sum := 0
		for _, m := range i {
			sum += getTokenCount(m.Content)
		}

		return sum
	default:
		return 0
	}
}


func makeFileString(files []db.File) string {
	var fileString = ""

	for _, f := range files {
		name := f.FileName
		content := f.Content

		fileString += fmt.Sprintf("## Name: %s\nContent:\n%s\n\n", name, content)
	}

	return fileString
}

func getContextLimit(ctx context.Context, model string) float64 {
	url := os.Getenv("OLLAMA_URL")
	
	if url == "" {
        return 0
	}

    oc, err := api.NewOllamaClient(url)

    if err != nil {
        return 0
    }

	res, err := oc.GetShow(ctx, model)

    if err != nil {
        return 0
    }

	count := res.ModelInfo["general.parameter_count"]

	if val, ok := count.(float64); ok {
		return val
	}
	
	return 0
}

func addFileToMessages(ctx context.Context, model string, messages []ollama.Message, files[]db.File) ([]ollama.Message, bool){
	fileString := makeFileString(files)
	m := messages[len(messages)-1]

	temp := fmt.Sprintf(FilePrompt, fileString, m.Content)

	estimatedTokens := getTokenCount(messages[:len(messages) -1]) + getTokenCount(temp)
	limit := getContextLimit(ctx, model) * LIMIT_BUFFER

	if float64(estimatedTokens) > limit {
		return messages, false
	}

	messages[len(messages)-1].Content = temp

	return messages, true
}

