package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
    client *api.Client
    ctx context.Context
}

func NewOllamaClient(_url string, ctx context.Context) (*OllamaClient, error) {
    u, err := url.Parse(_url)
    if err != nil {
        return nil, err
    }

    client := api.NewClient(u, http.DefaultClient)
    return &OllamaClient{
        client: client,
        ctx: ctx,
    }, nil
}

func (oc OllamaClient) Chat(query string) (string, error) {
    req := &api.ChatRequest{
        Model: "qwen3:1.7b",
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
    }

    var fullResponse string
    err := oc.client.Chat(oc.ctx, req, func(resp api.ChatResponse) error {
        fullResponse += resp.Message.Content
        return nil
    })

    if err != nil {
        return "", err
    }


    return fullResponse, nil
}
