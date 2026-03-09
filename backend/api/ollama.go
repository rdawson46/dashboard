package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/ollama/ollama/api"
	ollamaModel "github.com/ollama/ollama/types/model"
)

var (
    ollamaSemaphore chan struct{}
    semOnce         sync.Once
)

func InitOllamaSemaphore(maxConcurrent int) {
    semOnce.Do(func() {
        ollamaSemaphore = make(chan struct{}, maxConcurrent)
    })
}

func acquireSemaphore(ctx context.Context) error {
    if ollamaSemaphore == nil {
        return nil
    }
    select {
    case ollamaSemaphore <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func releaseSemaphore() {
    if ollamaSemaphore != nil {
        <-ollamaSemaphore
    }
}

type OllamaClient struct {
    client *api.Client
    logger *log.Logger
}


type StreamRequest struct {
	Query []api.Message `json:"messages"`
	Websearch bool `json:"webSearch"`
	Code bool `json:"code"`
	Model string `json:"model"`
	Username string `json:"username"`
    UserId string `json:"userId"`
    MessageId string `json:"messageId,omitempty"`
	UseRag bool `json:"useRag,omitempty"`
	FileIds []string `json:"fileIds,omitempty"`
}

func NewOllamaClient(_url string, logger *log.Logger) (*OllamaClient, error) {
    u, err := url.Parse(_url)
    if err != nil {
        return nil, err
    }

    if logger == nil {
        logger = log.New(os.Stderr)
        logger.SetOutput(io.Discard)
        logger.Info("No logger provided, using a silent one.")
    }

    client := api.NewClient(u, http.DefaultClient)
    return &OllamaClient{
        client: client,
        logger: logger.WithPrefix("ollama"),
    }, nil
}

// simple for now but will progress in complexity and will need to work with
func (oc OllamaClient) newRequest(query string, stream *bool) *api.ChatRequest {
	model := os.Getenv("DEFAULT_MODEL")

	if model == "" {
        oc.logger.Error("DEFAULT_MODEL environment variable not set")
		return nil
	}


    return &api.ChatRequest{
        Model: model,
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
        Model: model,
        Messages: totalMessages,
        Options: map[string]any{
            "temperature": 0.7,
            "top_k": 40,
            "top_p": 0.9,
        },
        Stream: &[]bool{stream}[0],
        Think: &[]bool{false}[0],
    }
}

func (oc OllamaClient) GetModelList(ctx context.Context) (*api.ListResponse, error) {
    oc.logger.Info("Fetching model list")
    resp, err := oc.client.List(ctx)
    if err != nil {
        oc.logger.Error("Failed to fetch model list", "error", err)
    }
    return resp, err
}

type ShowResponse struct {
    SystemInfo string
    ModelInfo map[string]any
    Capabilities  []ollamaModel.Capability
    ProjectedInfo map[string]any
}

func (oc OllamaClient) GetShow(ctx context.Context, model string) (*ShowResponse, error) {
    oc.logger.Info("Fetching model details", "model", model)
    req := &api.ShowRequest {
        Model: model,
    }

    res, err := oc.client.Show(ctx, req)

    if err != nil {
        oc.logger.Error("Failed to fetch model details", "model", model, "error", err)
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
    oc.logger.Info("Starting chat", "query", query)
    
    if err := acquireSemaphore(ctx); err != nil {
        oc.logger.Error("Failed to acquire semaphore", "error", err)
        return "", err
    }
    defer releaseSemaphore()

    req := oc.newRequest(query, &[]bool{false}[0])

	if req == nil {
        err := errors.New("Unable to make request")
        oc.logger.Error("Failed to create new chat request", "error", err)
		return "", err
	}

    var fullResponse string
    err := oc.client.Chat(ctx, req, func(resp api.ChatResponse) error {
        fullResponse += resp.Message.Content
        return nil
    })

    if err != nil {
        oc.logger.Error("Chat failed", "error", err)
        return "", err
    }

    oc.logger.Info("Chat finished successfully")
    return fullResponse, nil
}

func (oc OllamaClient) Stream(ctx context.Context, userReq StreamRequest, model string, msgChan chan any, errChan chan error) {
    defer close(msgChan)
    defer close(errChan)

    oc.logger.Info("Starting stream", "model", model, "user", userReq.Username)

    if err := acquireSemaphore(ctx); err != nil {
        oc.logger.Error("Failed to acquire semaphore", "error", err)
        errChan <- err
        return
    }
    defer releaseSemaphore()

    req := oc.newRequestWithMessages(userReq.Query, model, true)

	tools := api.Tools{}

	if userReq.Code {
        oc.logger.Info("Adding code execution tool")
		tools = append(tools, getCodeExecution())
	}

	if userReq.Websearch {
        oc.logger.Info("Adding web search tool")
		tools = append(tools, getWebSearchTool())
	}

	req.Tools = tools

	var handler api.ChatResponseFunc

	handler = func(resp api.ChatResponse) error {
        msgChan <- resp

		if len(resp.Message.ToolCalls) > 0 {
            oc.logger.Info("Tool call detected", "count", len(resp.Message.ToolCalls))
            toolResponses := oc.toolHandler(resp.Message)

            req.Messages = append(req.Messages, resp.Message)
            req.Messages = append(req.Messages, toolResponses...)

			for _, toolResp := range toolResponses {
				msgChan <- toolResp
			}

            oc.logger.Info("Re-querying model with tool responses")
            return oc.client.Chat(ctx, req, handler)
		}

        return nil
	}


	err := oc.client.Chat(ctx, req, handler)

    if err != nil {
        oc.logger.Error("Stream failed", "error", err)
        errChan <- err
    }
    oc.logger.Info("Stream finished")
}

// pass the message?
// return []tool response
func (oc OllamaClient) toolHandler(message api.Message) []api.Message {
    var toolResponse []api.Message

    for _, toolCall := range message.ToolCalls {
        oc.logger.Info("Processing tool call", "tool", toolCall.Function.Name)
        switch toolCall.Function.Name {
        case "web search":
            query, ok := toolCall.Function.Arguments["query"]
            if !ok {
                oc.logger.Warn("Tool 'web search' called without 'query' argument")
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
                oc.logger.Warn("Tool 'web search' argument 'query' is not a string")
                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: "Query param must be provided as a string",
                    },
                )
                continue
			}
            oc.logger.Info("Executing web search", "query", q)
			result, err := WebSearch(q)
			if err != nil {
                oc.logger.Error("Web search failed", "error", err)
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
            oc.logger.Info("Web search completed")

        case "python code execution":
            code, ok := toolCall.Function.Arguments["code"]
            if !ok {
                oc.logger.Warn("Tool 'python code execution' called without 'code' argument")
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
                oc.logger.Warn("Tool 'python code execution' argument 'code' is not a string")
                toolResponse = append(
                    toolResponse,
                    api.Message{
                        Role: "tool",
                        Content: "Code param must be provided as a string",
                    },
                )
                continue
            }
            oc.logger.Info("Executing python code", "code", c)
            result, err := ExecutePython(c)

            if err != nil {
                oc.logger.Error("Python execution failed", "error", err)
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

            oc.logger.Info("Python execution completed")
            content := fmt.Sprintf("Result: %s\nError: %s", result.Result, result.Error)
            toolResponse = append(
                toolResponse,
                api.Message{
                    Role: "tool",
                    Content: content,
                },
            )

        default:
            oc.logger.Warn("Unknown tool called", "tool", toolCall.Function.Name)
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
