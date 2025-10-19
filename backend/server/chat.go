package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	ollama "github.com/ollama/ollama/api"
	api "github.com/rdawson46/dashboard/api"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
    type request struct {
        Query string `json:"query"`
    }

    var chatReq request
    err := json.NewDecoder(r.Body).Decode(&chatReq)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	url := os.Getenv("OLLAMA_URL")
	
	if url == "" {
        http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
        return
	}

    ctx := r.Context()

    oc, err := api.NewOllamaClient(url)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    res, err := oc.Chat(ctx, chatReq.Query)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    response := map[string]string{"respones": res}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// will work on to replace the chatHandler
// TODO: add in user ID as well and along with username
func (s *Server) streamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var chatReq api.StreamRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := userFromContext(r.Context())

	// user not found
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
		return
	}

	if user.Username != chatReq.Username {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Invalid username"})
		return
	}

    // check for message id
    var messageId string
    if chatReq.MessageId == "" {
        s.logger.Info(
            "Creating new chat for user", 
            "user", user.Username,
        )

        id, err := s.db.CreateMessage(r.Context(), user.ID, chatReq.Query)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"Error": "Unable to Create Chat"})
            s.logger.Error(
                "Failed to Create Chat",
                "model", chatReq.Model,
                "remote", r.RemoteAddr,
                "user", user.Username,
            )
            return
        }

        messageId = id

        s.logger.Info(
            "New Chat ID",
            "User", user.Username,
            "Chat Id", id,
        )
    } else {
        messageId = chatReq.MessageId
    }

	// headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control", "*")

	flusher, ok := w.(http.Flusher)

	if !ok {
		s.logger.Error("Failed to get flusher")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	idMessage := map[string]any{
		"type": "Message ID",
		"messageId": messageId,
	}

	b, err := json.Marshal(idMessage)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, "failed to marshal resp", http.StatusInternalServerError)
		return
	}


	s.logger.Info(
		"Setting chat Id",
		"User", user.Username,
		"Chat Id", messageId,
	)

	fmt.Fprintf(w, "data: %s\n\n", b)
	flusher.Flush()

	s.logger.Info(
		"Chat Id set",
		"User", user.Username,
		"Chat Id", messageId,
	)

	url := os.Getenv("OLLAMA_URL")

	if url == "" {
		s.logger.Error("Ollama URL not set")
		http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	oc, err := api.NewOllamaClient(url)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	msgChan := make(chan any)
	errChan := make(chan error)


	s.logger.Info(
		"Setting model preference",
		"model", chatReq.Model,
		"user", user.Username,
	)

	s.db.SetPerferredModel(ctx, user.ID, chatReq.Model)

	s.logger.Info(
		"Starting stream",
		"model", chatReq.Model,
		"remote", r.RemoteAddr,
		"user", user.Username,
		"chatId", messageId,
	)

	go oc.Stream(ctx, chatReq, chatReq.Model, msgChan, errChan)


token_count := 0
    chatRespone := ""
	OuterLoop:
	for {
		select {
		case resp, ok := <-msgChan:
			if !ok {
				break OuterLoop
			}

            var b []byte
            var err error
            switch resp := resp.(type) {
            case ollama.ChatResponse:
                switch resp.Message.Role {
                case "assistant":
                    chatRespone += resp.Message.Content
                case "tool":
                    chatReq.Query = append(
                        chatReq.Query,
                        resp.Message,
                    )
                }
                token_count = resp.EvalCount;

                b, err = json.Marshal(map[string]any{
                    "type": "response",
                    "data": resp,
                })

                if err != nil {
                    s.logger.Error(err.Error())
                    http.Error(w, "failed to marshal resp", http.StatusInternalServerError)
                    return
                }

            case ollama.Message:
                b, err = json.Marshal(map[string]any{
                    "type": "message",
                    "data": resp,
                })

            default:
                s.logger.Warn("Stream sent unknown data type through channel")
            }

			fmt.Fprintf(w, "data: %s\n\n", b)
			flusher.Flush()
		case err, ok := <-errChan:
			if !ok {
				break OuterLoop
			}

			// HACK: replace error message
			s.logger.Error(err.Error())

			// TODO: superfluous error
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Fprintf(w, "data: %s\n\n", b)
			flusher.Flush()
			return
		}
	}

	s.logger.Info(
		"Streaming finished",
		"token_count", token_count,
		"remote", r.RemoteAddr,
		"user", user.Username,
		"chatId", messageId,
	)

	fmt.Fprintf(w, "data: %s\n\n", `{"done": true}`)
	flusher.Flush()

    final := ollama.Message{
    	Role:      "assistant",
    	Content:   chatRespone,
    }

    chatReq.Query = append(
        chatReq.Query,
        final,
    )

    success, err := s.db.AddMessage(r.Context(), messageId, user.ID, chatReq.Query)

    if err != nil {
        s.logger.Error(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if !success {
        s.logger.Error("Failed to add messages")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	s.logger.Info(
		"Chat successfully updated",
		"user", user.Username,
        "chat id", messageId,
		"chatId", messageId,
	)
}

func (s *Server) chatDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	user, ok := userFromContext(ctx)

	if !ok {
		s.logger.Error(
			"User not in context",
			"path", r.URL.Path,
		)
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	s.logger.Info(
		"Pulling user history",
		"userId", user.ID,
		"username", user.Username,
	)

	history, err := s.db.GetDescriptions(ctx, user.ID, 10, 0)

	if err != nil {
		s.logger.Error(
			"Failed to get descriptions",
			"error", err.Error(),
			"path", r.URL.Path,
		)

		http.Error(w, "Unable to get chat history", http.StatusInternalServerError)
		return
	}

	s.logger.Info(
		"History pulled",
		"length", len(history),
		"path", r.URL.Path,
	)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(history)
}
