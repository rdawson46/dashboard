package server

import (
	"encoding/json"
	"fmt"
	"net/http"
    "strconv"
	"os"

	api "github.com/rdawson46/dashboard/api"
)

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
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Status": "failed", "Message": "Row failed to delete"})
        return
    }

    if !ok {
        s.logger.Error(
            "Did not delete row",
            "userId", delReq.UserId,
        )
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Status": "failed", "Message": "Row not deleted"})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"Status": "ok", "Message": ""})
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Status": "failed", "Message": "Failed to get messages"})
        return
	}

    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
	return
}
