package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	ollamaModel "github.com/ollama/ollama/types/model"
)

type OllamaClient struct {
    client *api.Client
}

type StreamRequest struct {
	Query []api.Message `json:"messages"`
	Websearch bool `json:"webSearch"`
	Code bool `json:"code"`
	Model string `json:"model"`
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
					"query": struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
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
			Description: "execute python code dynamically",
			Parameters: struct{Type string "json:\"type\""; Defs any "json:\"$defs,omitempty\""; Items any "json:\"items,omitempty\""; Required []string "json:\"required\""; Properties map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""} "json:\"properties\""}{
				Type: "object",
				Properties: map[string]struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
					"code": struct{Type api.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
						Type: []string{"string"},
						Description: "code that you produced",
					},
				},
				Required: []string{"code"},
			},
		},
	}
}

func NewOllamaClient(_url string) (*OllamaClient, error) {
    u, err := url.Parse(_url)
    if err != nil {
        return nil, err
    }

    client := api.NewClient(u, http.DefaultClient)
    return &OllamaClient{
        client: client,
    }, nil
}

// simple for now but will progress in complexity and will need to work with
func (oc OllamaClient) newRequest(query string, stream *bool) *api.ChatRequest {
    return &api.ChatRequest{
        // Model: "qwen3:1.7b",
        Model: "gemma3:1b",
        Messages: []api.Message{
            {
                Role: "system",
                Content: "You are a helpful assistant that answers in a concise manner",
            },
            {
                Role: "user",
                Content: query,
            },
        },
        Options: map[string]any{
            "temperature": 0.7,
            "top_k": 40,
            "top_p": 0.9,
        },
        Stream: stream,
    }
}

func (oc OllamaClient) newRequestWithMessages(messages []api.Message, model string, stream bool) *api.ChatRequest {
    totalMessages := append([]api.Message{
        {
            Role: "system",
            Content: "You are a helpful assistant that answers in a concise manner",
        },
    }, messages...)


    return &api.ChatRequest{
        // Model: "qwen3:1.7b",
        Model: model,
        Messages: totalMessages,
        Options: map[string]any{
            "temperature": 0.7,
            "top_k": 40,
            "top_p": 0.9,
        },
        Stream: &[]bool{stream}[0],
    }
}

func (oc OllamaClient) GetModelList(ctx context.Context) (*api.ListResponse, error) {
    resp, err := oc.client.List(ctx)
    return resp, err
}

type ShowResponse struct {
    SystemInfo string
    ModelInfo map[string]any
    Capabilities  []ollamaModel.Capability
    ProjectedInfo map[string]any
}

func (oc OllamaClient) GetShow(ctx context.Context, model string) (*ShowResponse, error) {
    req := &api.ShowRequest {
        Model: model,
    }
    // TODO: simplify response for UI, remove extra fields that won't be used
    res, err := oc.client.Show(ctx, req)

    if err != nil {
        return nil, err
    }

    r := &ShowResponse {
        SystemInfo: res.System,
        ModelInfo: res.ModelInfo,
        Capabilities: res.Capabilities,
        ProjectedInfo: res.ProjectorInfo,
    }

    return r, err
}

func (oc OllamaClient) Chat(ctx context.Context, query string) (string, error) {
    req := oc.newRequest(query, &[]bool{false}[0])

    var fullResponse string
    err := oc.client.Chat(ctx, req, func(resp api.ChatResponse) error {
        fullResponse += resp.Message.Content
        return nil
    })

    if err != nil {
        return "", err
    }

    return fullResponse, nil
}

func (oc OllamaClient) Stream(ctx context.Context, userReq StreamRequest, model string, msgChan chan api.ChatResponse, errChan chan error) {
    defer close(msgChan)
    defer close(errChan)

    req := oc.newRequestWithMessages(userReq.Query, model, true)

	tools := api.Tools{}

	if userReq.Code {
		tools = append(tools, getCodeExecution())
	}

	if userReq.Websearch {
		tools = append(tools, getWebSearchTool())
	}

	req.Tools = tools

	var handler api.ChatResponseFunc

	handler = func(resp api.ChatResponse) error {
        msgChan <- resp

		if len(resp.Message.ToolCalls) > 0 {
			for _, toolCall := range resp.Message.ToolCalls {
				switch toolCall.Function.Name {
				case "web search":
					// will be query
					_, ok := toolCall.Function.Arguments["query"]

					if !ok {
						// TODO: continue
						break
					}

				case "python code execution":
					fmt.Println("\n\nCalling Code\n\n")
					code, ok := toolCall.Function.Arguments["code"]

					if !ok {
						// TODO: handle
						break
					}

					c, ok := code.(string)

					if !ok {
						// TODO: handle
						break
					}

					result, err := ExecutePython(c)

					if err != nil {
						// TODO: handle
						fmt.Printf("\nError: %s\n", err.Error())
						fmt.Printf("\nCode: %s\n", c)

						break
					}

					content := fmt.Sprintf("Result: %s\nError: %s", result.Result, result.Error)

					req.Messages = append(
						req.Messages,
						resp.Message,
						api.Message{
							Role: "tool",
							Content: content,
						},
					)

					return oc.client.Chat(ctx, req, handler)
				default:
					fmt.Println("can not find tool")
				}
			}
		}

        return nil
	}


	err := oc.client.Chat(ctx, req, handler)

    if err != nil {
        errChan <- err
    }
}
