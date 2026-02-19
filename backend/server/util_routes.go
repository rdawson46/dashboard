package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	api "github.com/rdawson46/dashboard/api"
	ollama "github.com/ollama/ollama/api"
)

func (s *Server) modelListHandler(w http.ResponseWriter, r *http.Request) {
	type modelResponse struct {
		Models []ollama.ListModelResponse `json:"models"`
		Preference string `json:"preference"`
	}

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	user, ok := userFromContext(r.Context())

	if !ok {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
        return
	}

	url := os.Getenv("OLLAMA_URL")
	
	if url == "" {
        http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
        return
	}

    ctx := r.Context()

    oc, err := api.NewOllamaClient(url, s.logger)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    resp, err := oc.GetModelList(ctx)

    if err != nil {
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        return 
    }

	modelPreference, err := s.db.GetPerferredModel(ctx, user.ID)

	if err != nil {
		s.logger.Error("Could not load model preference", "error", err.Error())
		modelPreference = ""
	}

	finalRes := modelResponse{
		Models: resp.Models,
		Preference: modelPreference,
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(finalRes)
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

    oc, err := api.NewOllamaClient(url, nil)

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

    if user.ID != delReq.UserId {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"Error": "Invalid username"})
        return
    }

    w.Header().Set("Content-Type", "application/json")

    ok, err = s.db.DeleteMessage(r.Context(), delReq.ChatId, delReq.UserId)

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

	if user.ID != getReq.UserId {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Invalid username"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	messages, err := s.db.GetMessage(r.Context(), getReq.UserId, getReq.ChatId)

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
