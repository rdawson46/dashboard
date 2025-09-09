package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	ollama "github.com/ollama/ollama/api"
)

func __main() {
    err := godotenv.Load()

    if err != nil {
        fmt.Println(err)
        return
    }

	_url := os.Getenv("OLLAMA_URL")

    u, err := url.Parse(_url)
    if err != nil {
        fmt.Println(err)
        return
    }

    client := ollama.NewClient(u, http.DefaultClient)

    query := "What is 2 + 2? Use the addition tool."

    request := &ollama.ChatRequest{
        // Model: "qwen3:1.7b",
        Model: "qwen3:0.6b",
        // Model: "gemma3:1b",
        Messages: []ollama.Message{
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
        Stream: &[]bool{true}[0],
        Tools: []ollama.Tool{
            ollama.Tool{
                Type: "string",
                Function: ollama.ToolFunction {
                    Name: "addition",
                    Description: "used to add two numbers together",
                    Parameters: struct{Type string "json:\"type\""; Defs any "json:\"$defs,omitempty\""; Items any "json:\"items,omitempty\""; Required []string "json:\"required\""; Properties map[string]struct{Type ollama.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""} "json:\"properties\""}{
                        Type: "object",
                        Properties: map[string]struct{Type ollama.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
                            "a": struct{Type ollama.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
                                Type: []string{"int"},
                                Description: "first number",
                            },
                            "b": struct{Type ollama.PropertyType "json:\"type\""; Items any "json:\"items,omitempty\""; Description string "json:\"description\""; Enum []any "json:\"enum,omitempty\""}{
                                Type: []string{"int"},
                                Description: "second number",
                            },
                        },
                        Required: []string{"a", "b"},
                    },
                },
            },
        },
    }

    fmt.Printf("%+v\n\n", request)


    ctx := context.Background()

    err = client.Chat(ctx, request, func(resp ollama.ChatResponse) error {
        fmt.Printf("%+v\n\n", resp)

        if len(resp.Message.ToolCalls) > 0 {
            for _, toolCall := range resp.Message.ToolCalls {
                if toolCall.Function.Name == "addition" {
                    args := toolCall.Function.Arguments

                    a := args["a"].(float64)
                    b := args["b"].(float64)

                    result := a + b

                    fmt.Println(result)



                    toolCallReq := &ollama.ChatRequest{
                        Model: request.Model,
                        Messages: append(
                            request.Messages,
                            resp.Message,
                            ollama.Message{
                                Role: "tool",
                                Content: fmt.Sprintf("result of addition is: %d", result),
                                ToolCalls: []ollama.ToolCall{
                                    {
                                        Function: ollama.ToolCallFunction{
                                            Index: toolCall.Function.Index,
                                            Name: toolCall.Function.Name,
                                            Arguments: toolCall.Function.Arguments,
                                        },
                                    },
                                },
                            },
                        ),
                        Stream: &[]bool{true}[0],
                    }

                    return client.Chat(ctx, toolCallReq, func(finalResp api.ChatResponse) error {
                        fmt.Printf("%+v\n\n", finalResp.Message.Content)
                        return nil
                    })
                }
            }
        }
        fmt.Println("here")

        return nil
    })

    if err != nil {
        fmt.Println(err)
        return
    }
}
