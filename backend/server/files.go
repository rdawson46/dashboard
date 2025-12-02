package server

import (
	"encoding/json"
	"io"
	"net/http"
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
	if err := r.ParseMultipartForm(2 << 20); err != nil {
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

}

func (s *Server) getFileList(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"files": []int{},
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func (s *Server) deleteFile(w http.ResponseWriter, r *http.Request) {

}
