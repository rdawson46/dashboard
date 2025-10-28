package api

import (
	"github.com/ollama/ollama/api"
)
type ToolResult struct {
	Result string `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

func getWebSearchTool() api.Tool {
	return api.Tool{
		Type: "string",
		Function: api.ToolFunction {
			Name: "web search",
			Description: "tool for searching the web with the user's query",
			Parameters: struct{Type string "json:\"type\""; Defs any "json:\"$defs,omitempty\""; Items any "json:\"items,omitempty\""; Required []string "json:\"required\""; Properties map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""} "json:\"properties\""}{
				Type: "object",
				Properties: map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
					"query": {
						Type: []string{"string"},
						Description: "query to help answer the user's question",
					},
				},
				Required: []string{"query"},
			},
		},
	}
}

func getCodeExecution() api.Tool {
	return api.Tool{
		Type: "string",
		Function: api.ToolFunction {
			Name: "python code execution",
			Description: "Execute python code in a sandboxed environment. You can use this tool to perform calculations, manipulate strings, and other simple tasks. Make sure to print all results to view them. For example, to add two numbers, you can write: `print(1 + 2)`.",
			Parameters: struct{Type string "json:\"type\""; Defs any "json:\"$defs,omitempty\""; Items any "json:\"items,omitempty\""; Required []string "json:\"required\""; Properties map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""} "json:\"properties\""}{
				Type: "object",
				Properties: map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
					"code": {
						Type: []string{"string"},
						Description: "code that you produced",
					},
				},
				Required: []string{"code"},
			},
		},
	}
}
