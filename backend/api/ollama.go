package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	ollamaModel "github.com/ollama/ollama/types/model"
)

type OllamaClient struct {
    client *api.Client
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

func (oc OllamaClient) Stream(ctx context.Context, messages []api.Message, model string, msgChan chan api.ChatResponse, errChan chan error) {
    defer close(msgChan)
    defer close(errChan)

    req := oc.newRequestWithMessages(messages, model, true)

    err := oc.client.Chat(ctx, req, func(resp api.ChatResponse) error {
        msgChan <- resp
        return nil
    })

    if err != nil {
        errChan <- err
    }
}
