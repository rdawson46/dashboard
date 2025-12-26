package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	ollama "github.com/ollama/ollama/api"
	api "github.com/rdawson46/dashboard/api"
	"github.com/rdawson46/dashboard/db"
)

func SetSSE(w http.ResponseWriter) (http.Flusher, bool){
	// headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control", "*")

	flusher, ok := w.(http.Flusher)

	return flusher, ok
}

func GetMessageID(s *Server, r *http.Request, chatReq api.StreamRequest, user *User_jwt) (*db.CreateMessage, error) {
    // check for message id
	var messageCreation *db.CreateMessage
    if chatReq.MessageId == "" {
        s.logger.Info(
            "Creating new chat for user", 
            "user", user.Username,
        )

		temp, err := s.db.CreateMessage(r.Context(), user.ID, chatReq.Model, chatReq.Query)
		messageCreation = &temp

        if err != nil {
            s.logger.Error(
                "Failed to Create Chat",
                "model", chatReq.Model,
                "remote", r.RemoteAddr,
                "user", user.Username,
            )
            return nil, err
        }

        s.logger.Info(
            "New Chat ID",
            "User", user.Username,
            "Chat Id", messageCreation.Id,
        )
    } else {
		messageCreation = &db.CreateMessage{
			Id: chatReq.MessageId,
			New: false,
		}
    }

	return messageCreation, nil
}

func SendSSEMessage(s *Server, flusher http.Flusher, w http.ResponseWriter, data any, messageType string) error {
	b, err := json.Marshal(map[string]any{
		"type": messageType,
		"data": data,
	})

	// NOTE: might have to move this out due to logger line nums
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, "failed to marshal resp", http.StatusInternalServerError)
		return err
	}

	fmt.Fprintf(w, "data: %s\n\n", b)
	flusher.Flush()
	return nil
}

func SendSimpleSEE(s *Server, flusher http.Flusher, w http.ResponseWriter, message string) {
	fmt.Fprintf(w, "data: %s\n\n", message)
	flusher.Flush()
}

func Streamer(s *Server, w http.ResponseWriter, flusher http.Flusher, r *http.Request, chatReq api.StreamRequest, messageId string, user *User_jwt) (string, error) {
	url := os.Getenv("OLLAMA_URL")

	if url == "" {
		s.logger.Error("Ollama URL not set")
		http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
		return "", errors.New("No ollama url")
	}

	ctx := r.Context()

	oc, err := api.NewOllamaClient(url)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", err
	}

	msgChan := make(chan any)
	errChan := make(chan error)

	s.logger.Info(
		"Starting stream",
		"model", chatReq.Model,
		"remote", r.RemoteAddr,
		"user", user.Username,
		"chatId", messageId,
	)

	go oc.Stream(ctx, chatReq, chatReq.Model, msgChan, errChan)

	token_count := 0
    chatResponse := ""

OuterLoop:
	for {
		select {
		case resp, ok := <-msgChan:
			if !ok {
				break OuterLoop
			}

            switch resp := resp.(type) {
            case ollama.ChatResponse:
                switch resp.Message.Role {
                case "assistant":
                    chatResponse += resp.Message.Content
                case "tool":
                    chatReq.Query = append(
                        chatReq.Query,
                        resp.Message,
                    )
                }
                token_count = resp.EvalCount;

				err = SendSSEMessage(s, flusher, w, resp, "response")
				if err != nil { return "", err }

            case ollama.Message:
				err = SendSSEMessage(s, flusher, w, resp, "message")
				if err != nil { return "", err }

            default:
                s.logger.Warn("Stream sent unknown data type through channel")
            }

		case err, ok := <-errChan:
			if !ok {
				break OuterLoop
			}

			s.logger.Error(
				"ErrChan received message",
				"Err", err.Error(),
			)

			err = SendSSEMessage(s, flusher, w, map[string]any{"error": "broke"}, "error")
			return "", err
		}
	}

	s.logger.Info(
		"Streaming finished",
		"token_count", token_count,
		"remote", r.RemoteAddr,
		"user", user.Username,
		"chatId", messageId,
	)

	SendSimpleSEE(s, flusher, w, `{"done": true}`)

	return chatResponse, nil
}

func SaveMessage(s *Server, r *http.Request, messageId, chatResponse string, chatReq *api.StreamRequest, user *User_jwt) {
    final := ollama.Message{
    	Role:      "assistant",
    	Content:   chatResponse,
    }

    chatReq.Query = append(chatReq.Query, final)

    success, err := s.db.AddMessage(r.Context(), messageId, user.ID, chatReq.Model, chatReq.Query)

    if err != nil {
        s.logger.Error("Failed to add messages", "Err", err.Error())
		return
    }

    if !success {
        s.logger.Error("Failed to add messages", "Err", "Failed to Add Message")
		return
    }

	s.logger.Info(
		"Chat successfully updated",
		"user", user.Username,
        "chat id", messageId,
		"chatId", messageId,
	)
}
