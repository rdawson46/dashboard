package server

// TODO: add logging every where
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
    "strconv"

	ollama "github.com/ollama/ollama/api"
	api "github.com/rdawson46/dashboard/api"
)

// ========== ROUTING STRUCTS ==========
type HealthResponse struct {
    Status string `json:"status"`
}

// =====================================


// ========== ROUTING FUNCTIONS ==========
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/health", http.StatusTemporaryRedirect)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{Status: "ok"}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

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
    var messageIdString string
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

        messageIdString = strconv.FormatInt(id, 10)

        s.logger.Info(
            "New Chat ID",
            "User", user.Username,
            "Chat Id", id,
        )
    } else {
        messageIdString = chatReq.MessageId
    }

    messageId, err := strconv.ParseInt(messageIdString, 10, 64)

    if err != nil {
        s.logger.Errorf("Invalid id: %s\nError: %s", messageIdString, err.Error())
        http.Error(w, "Failed to indentify chat", http.StatusInternalServerError)
        return
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


	// TODO: fix this logging
	s.logger.Info("Setting chat Id")
	fmt.Fprintf(w, "data: %s\n\n", b)
	flusher.Flush()
	s.logger.Info("Chat Id set")

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

	msgChan := make(chan ollama.ChatResponse)
	errChan := make(chan error)

	s.logger.Info(
		"Starting stream",
		"model", chatReq.Model,
		"remote", r.RemoteAddr,
		"user", user.Username,
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

            switch resp.Message.Role{
            case "assistant":
                chatRespone += resp.Message.Content
            case "tool":
                chatReq.Query = append(
                    chatReq.Query,
                    resp.Message,
                )
            }

			// encode resp
			b, err := json.Marshal(resp)

			if err != nil {
				s.logger.Error(err.Error())
				http.Error(w, "failed to marshal resp", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", b)
			flusher.Flush()
			token_count = resp.EvalCount;
		case err, ok := <-errChan:
			if !ok {
				break OuterLoop
			}

			// HACK: replace error message
			s.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.logger.Info(
		"Streaming finished",
		"token_count", token_count,
		"remote", r.RemoteAddr,
		"user", user.Username,
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
	)
}

func modelListHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

    resp, err := oc.GetModelList(ctx)

    if err != nil {
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resp.Models)
}

func modelShowHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    type request struct {
        Model string `json:"model"`
    }

    var showReq request
    err := json.NewDecoder(r.Body).Decode(&showReq)

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

    resp, err := oc.GetShow(ctx, showReq.Model)

    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting model details for: %s", showReq.Model), http.StatusInternalServerError)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resp)
}

// TODO: have to get the user ID from content, requires DB setup
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


	// TEMP: placeholder until db is set up
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


func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		http.Error(w, "Improper values", http.StatusUnauthorized)
		s.logger.Error(
			"Login failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	ctx := r.Context()
	user_db, err := s.db.SignInUser(ctx, username, password)

	if err != nil {
		http.Error(w, "Couldn't login", http.StatusUnauthorized)
		s.logger.Error(
			"Login failed: No user found",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	// TODO: make and set token in jwt
	var user User_jwt

	user.Username = user_db.Name
	user.ID = user_db.ID

	token, err := s.jwt_manager.GenerateToken(&user)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		s.logger.Error(
			"Failed to generate user token",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	s.jwt_manager.SetTokenCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s.jwt_manager.ClearTokenCookie(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// TODO: modify for cookies and figure where to trigger
func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := s.jwt_manager.GetTokenFromCookie(r)

	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized: No session cookie found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusUnauthorized)
		return
	}

	newToken, err := s.jwt_manager.RefreshToken(token)

	if err != nil {
		http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	s.jwt_manager.SetTokenCookie(w, newToken)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) verifyHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())

	if !ok {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
} 

// =======================================

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	confirm := strings.TrimSpace(r.FormValue("confirm"))

	if username == "" || password == "" || confirm == "" {
		http.Error(w, "Improper values", http.StatusUnauthorized)
		s.logger.Error(
			"User registration failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	if password != confirm {
		http.Error(w, "password and confirmation do not match", http.StatusUnauthorized)
		s.logger.Error(
			"User registration failed: password != confirmation",
			"path", r.URL.Path,
		)
		return
	}

	s.logger.Info("User information pulled",
		"path", r.URL.Path,
	)

	db_user, err := s.db.CreateUser(r.Context(), username, password)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusUnauthorized)
		s.logger.Error(
			"User registration failed: Failed to create in db",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	user := &User_jwt{} // username and ID
	user.Username = db_user.Name
	user.ID = db_user.ID

	token, err := s.jwt_manager.GenerateToken(user)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		s.logger.Error(
			"Failed to generate user token",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	s.jwt_manager.SetTokenCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}


func (s *Server) deleteChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	type deleteReqStruct struct  {
		ChatId string `json:"chatId"`
		UserId string `json:"userId"`
	}

	var delReq deleteReqStruct
	err := json.NewDecoder(r.Body).Decode(&delReq)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := userFromContext(r.Context())

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
		return
	}

	userId, err := strconv.ParseInt(delReq.UserId, 10, 64)

    if err != nil {
        s.logger.Error(
			"Invalid id",
			"userId", delReq.UserId,
			"error", err.Error(),
		)
        http.Error(w, "Failed to indentify chat", http.StatusInternalServerError)
        return
    }

	if user.ID != userId {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Invalid username"})
		return
	}

	chatId, err := strconv.ParseInt(delReq.ChatId, 10, 64)

    if err != nil {
        s.logger.Error(
			"Invalid id",
			"userId", delReq.UserId,
			"error", err.Error(),
		)
        http.Error(w, "Failed to indentify chat", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")

	ok, err = s.db.DeleteMessage(r.Context(), chatId, userId)

	if err != nil {
        s.logger.Error(
			"Failed to delete row",
			"userId", delReq.UserId,
			"error", err.Error(),
		)
		json.NewEncoder(w).Encode(map[string]string{"Status": "failed", "Message": "Row failed to delete"})
		w.WriteHeader(http.StatusInternalServerError)
        return
	}

	if !ok {
        s.logger.Error(
			"Did not delete row",
			"userId", delReq.UserId,
		)
		json.NewEncoder(w).Encode(map[string]string{"Status": "failed", "Message": "Row not deleted"})
		w.WriteHeader(http.StatusInternalServerError)
        return
	}

	json.NewEncoder(w).Encode(map[string]string{"Status": "ok", "Message": ""})
    w.WriteHeader(http.StatusOK)
	return
}


func (s *Server) getChatHandler(w http.ResponseWriter, r *http.Request) {
	/*

	- get user ID and the chat ID
	- check association
	- return the chat history

	*/

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	type getChatStruct struct {
		ChatId string `json:"chatId"`
		UserId string `json:"userId"`
	}

	var getReq getChatStruct
	err := json.NewDecoder(r.Body).Decode(&getReq)


	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := userFromContext(r.Context())


	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
		return
	}

	userId, err := strconv.ParseInt(getReq.UserId, 10, 64)

    if err != nil {
        s.logger.Errorf("Invalid id: %s\nError: %s", getReq.UserId, err.Error())
        http.Error(w, "Failed to indentify chat", http.StatusInternalServerError)
        return
    }

	if user.ID != userId {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Invalid username"})
		return
	}

	chatId, err := strconv.ParseInt(getReq.ChatId, 10, 64)

    if err != nil {
        s.logger.Errorf("Invalid id: %s\nError: %s", getReq.UserId, err.Error())
        http.Error(w, "Failed to indentify chat", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")

	messages, err := s.db.GetMessage(r.Context(), chatId)

	if err != nil {
		s.logger.Error(
			"Failed to load Messages:",
			"userId", getReq.UserId,
			"error", err.Error(),
		)
		json.NewEncoder(w).Encode(map[string]string{"Status": "failed", "Message": "Failed to get messages"})
		w.WriteHeader(http.StatusInternalServerError)
        return
	}

	json.NewEncoder(w).Encode(messages)
    w.WriteHeader(http.StatusOK)
	return
}


// =======================================



// ========== HANDLER FUNCTIONS ==========
// TODO: set up and api handler

func addRoutes(h *http.ServeMux, s *Server) {
    // basic routes
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)

    // user status routes
    h.HandleFunc("/api/login", s.loginHandler)
    h.HandleFunc("/api/logout", s.logoutHandler)

	// TODO: will have to set up as an api and requires DB
    h.HandleFunc("/api/reresh", s.refreshHandler)

    h.HandleFunc("/api/register", s.register)

    h.HandleFunc("/api/chat", s.jwt_manager.AuthApiMiddleware(s.rateLimitMiddleware(chatHandler)))

    h.HandleFunc("/api/stream", s.jwt_manager.AuthApiMiddleware(s.streamHandler))
    h.HandleFunc("/api/modelList", s.jwt_manager.AuthApiMiddleware(modelListHandler))
    h.HandleFunc("/api/modelInfo", s.jwt_manager.AuthApiMiddleware(modelShowHandler))

	h.HandleFunc("/api/chatDescription", s.jwt_manager.AuthApiMiddleware(s.chatDescriptionHandler))

	h.HandleFunc("/api/me", s.jwt_manager.AuthApiMiddleware(s.verifyHandler))

	h.HandleFunc("/api/messages", s.jwt_manager.AuthApiMiddleware(s.getChatHandler))
	h.HandleFunc("/api/deleteMessages", s.jwt_manager.AuthApiMiddleware(s.deleteChatHandler))
}

// =======================================
