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
			Description: "Search the internet for real-time information, news, documentation, or facts that are not available in your training data. Use this tool when the user's query requires current events or specific details from the web. The tool returns a list of relevant titles, URLs, and descriptions.",
			Parameters: struct{Type string "json:\"type\""; Defs any "json:\"$defs,omitempty\""; Items any "json:\"items,omitempty\""; Required []string "json:\"required\""; Properties map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""} "json:\"properties\""}{
				Type: "object",
				Properties: map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
					"query": {
						Type: []string{"string"},
						Description: "The search query to use for finding information.",
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
			Description: "Execute Python code in a sandboxed environment. Use this tool for complex calculations, data analysis, string manipulations, or any task that requires logical processing. Important: Only output printed using `print()` will be captured and returned. Ensure you print the final result or any intermediate steps you want to see. For example, `x = 5; y = 10; print(x + y)` will return '15'.",
			Parameters: struct{Type string "json:\"type\""; Defs any "json:\"$defs,omitempty\""; Items any "json:\"items,omitempty\""; Required []string "json:\"required\""; Properties map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""} "json:\"properties\""}{
				Type: "object",
				Properties: map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
					"code": {
						Type: []string{"string"},
						Description: "The Python code to execute. Remember to print the output.",
					},
				},
				Required: []string{"code"},
			},
		},
	}
}
