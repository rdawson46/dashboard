package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rdawson46/dashboard/db"
)

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {
	/*
	IDEA: upload a file to be parsed

		Either:
			* parse plain text, code etc
			* pdf to vision model for cheap OCR
		
		store in DB
		return the file ID

	GOAL: user based file-storage
		* allows for system wide RAG on user's files - down the road
		* allow for direct selection
	
	CURRENT REQUIREMENTS:
		* get file from user, store it, return the file uuid in from storage
	*/

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Max upload size: 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.logger.Error("Failed to parse multipart form", "error", err)
		http.Error(w, "Failed to parse multipart form", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		s.logger.Error("Error retrieving file from form-data", "error", err)
		http.Error(w, "Error retrieving file from form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	s.logger.Info("File uploaded", "fileName", handler.Filename, "fileSize", handler.Size, "mimeType", handler.Header.Get("Content-Type"), "user", user.Username)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("Error reading uploaded file", "error", err)
		http.Error(w, "Error reading uploaded file", http.StatusInternalServerError)
		return
	}
	
	fileContent := string(fileBytes)
	contentType := handler.Header.Get("Content-Type")

	if contentType == "application/pdf" {
		// TODO: Implement OCR with a vision model for PDFs
		http.Error(w, "Not available", http.StatusNotImplemented)
		return
	}

	fileID, err := s.db.SaveFile(r.Context(), user.ID, handler.Filename, contentType, fileContent)
	if err != nil {
		s.logger.Error("Failed to save file to DB", "error", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	
	s.logger.Info("File processed successfully", "fileID", fileID, "user", user.Username)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"fileID": fileID})
}

func (s *Server) getFile(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	fileId := r.URL.Query().Get("fileId")
	if fileId == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	file, err := s.db.GetFile(r.Context(), fileId, user.ID)
	if err != nil {
		if err == db.ErrFileNotFound {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		s.logger.Error("Failed to get file from DB", "error", err, "fileId", fileId, "userId", user.ID)
		http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.FileName))
	_, err = w.Write([]byte(file.Content))
	if err != nil {
		s.logger.Error("Failed to write file content to response", "error", err, "fileId", fileId)
		// It might be too late to send an HTTP error code if headers are already sent,
		// but logging is still important.
	}
}

func (s *Server) getFileList(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	
	totalItems, err := s.db.GetFilesCount(r.Context(), user.ID)
	if err != nil {
		s.logger.Error("Failed to get files count", "error", err, "userId", user.ID)
		http.Error(w, "Failed to get file count", http.StatusInternalServerError)
		return
	}

	files, err := s.db.GetFiles(r.Context(), user.ID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get files", "error", err, "userId", user.ID)
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}
	
	response := map[string]any{
		"files":      files,
		"totalItems": totalItems,
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func (s *Server) deleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	
	fileId := r.URL.Query().Get("fileId")
	if fileId == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	err := s.db.DeleteFile(r.Context(), fileId, user.ID)
	if err != nil {
		s.logger.Error("Failed to delete file", "error", err, "fileId", fileId, "userId", user.ID)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
