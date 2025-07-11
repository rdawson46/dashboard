package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
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

func (oc OllamaClient) newRequestWithMessages(messages []api.Message, stream bool) *api.ChatRequest {
    totalMessages := append([]api.Message{
        {
            Role: "system",
            Content: "You are a helpful assistant that answers in a concise manner",
        },
    }, messages...)


    return &api.ChatRequest{
        // Model: "qwen3:1.7b",
        Model: "gemma3:1b",
        Messages: totalMessages,
        Options: map[string]any{
            "temperature": 0.7,
            "top_k": 40,
            "top_p": 0.9,
        },
        Stream: &[]bool{stream}[0],
    }
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

func (oc OllamaClient) Stream(ctx context.Context, query string, msgChan chan api.ChatResponse, errChan chan error) {
    defer close(msgChan)
    defer close(errChan)

    req := oc.newRequest(query, &[]bool{true}[0])

    err := oc.client.Chat(ctx, req, func(resp api.ChatResponse) error {
        msgChan <- resp
        return nil
    })

    if err != nil {
        errChan <- err
    }
}

func (oc OllamaClient) Stream2(ctx context.Context, messages []api.Message, msgChan chan api.ChatResponse, errChan chan error) {
    defer close(msgChan)
    defer close(errChan)

    req := oc.newRequestWithMessages(messages, true)

    err := oc.client.Chat(ctx, req, func(resp api.ChatResponse) error {
        msgChan <- resp
        return nil
    })

    if err != nil {
        errChan <- err
    }
}


/*
func main() {
    url := "http://10.0.2.2:11434"
    client, err := NewOllamaClient(url)

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    ctx := context.Background()

    msgChan := make(chan api.ChatResponse)
    errChan := make(chan error)
    go client.Stream(ctx, "why is the sky blue?", msgChan, errChan)

    fmt.Println("Starting")

    for {
        select {
        case resp, ok := <- msgChan:
            if !ok {
                fmt.Println("done")
                return
            }
            fmt.Println(resp.Message.Content)
        case err, ok := <-errChan:
            if !ok {
                fmt.Println("done")
                return
            }
            fmt.Printf("error occurred: %v\n", err)
            return
        case <-time.After(30 * time.Second):
            fmt.Println("taken 30 seconds")
        }
    }
}
*/
