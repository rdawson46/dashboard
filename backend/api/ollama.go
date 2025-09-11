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
	Username string `json:"username"`
    UserId string `json:"userId"`
    MessageId string `json:"messageId,omitempty"`
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
			Description: "execute python code dynamically",
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

// func (oc OllamaClient) Stream(ctx context.Context, userReq StreamRequest, model string, msgChan chan api.ChatResponse, errChan chan error) {
func (oc OllamaClient) Stream(ctx context.Context, userReq StreamRequest, model string, msgChan chan any, errChan chan error) {
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
            toolResponses := toolHandler(resp.Message)

            req.Messages = append(req.Messages, resp.Message)
            req.Messages = append(req.Messages, toolResponses...)

			for _, toolResp := range toolResponses {
				msgChan <- toolResp
			}

            return oc.client.Chat(ctx, req, handler)
		}

        return nil
	}


	err := oc.client.Chat(ctx, req, handler)

    if err != nil {
        errChan <- err
    }
}

// pass the message?
// return []tool response
func toolHandler(message api.Message) []api.Message {
    var toolResponse []api.Message

    for _, toolCall := range message.ToolCalls {
        switch toolCall.Function.Name {
        case "web search":
            // will be query
            query, ok := toolCall.Function.Arguments["query"]

            if !ok {
                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: "Query parameter is required",
                    },
                )
                continue
            }

			q, ok := query.(string)

			if !ok {
                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: "Query param must be provided as a string",
                    },
                )
                continue
			}

			result, err := WebSearch(q)

			if err != nil {
                content := fmt.Sprintf("\nError: %s\n", err.Error())

                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: content,
                    },
                )
                continue
			}

            content := fmt.Sprintf("Result: %s\nError: %s", result.Result, result.Error)

            toolResponse = append(
                toolResponse,
                api.Message{
                    Role: "tool",
                    Content: content,
                },
            )

        case "python code execution":
            code, ok := toolCall.Function.Arguments["code"]



            if !ok {
                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: "Code parameter is required",
                    },
                )
                continue
            }

            c, ok := code.(string)

            if !ok {
                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: "Code param must be provided as a string",
                    },
                )
                continue
            }

            result, err := ExecutePython(c)

            if err != nil {
                content := fmt.Sprintf("\nError: %s\n", err.Error())

                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: content,
                    },
                )
                continue
            }

            fmt.Println("code worked")
            content := fmt.Sprintf("Result: %s\nError: %s", result.Result, result.Error)

            toolResponse = append(
                toolResponse,
                api.Message{
                    Role: "tool",
                    Content: content,
                },
            )

        default:
            toolResponse = append(
                toolResponse,
                api.Message{
                    Role: "tool",
                    Content: fmt.Sprintf("Failed to find tool: %s", toolCall.Function.Name),
                },
            )
        }
    }
    
    return toolResponse
}
