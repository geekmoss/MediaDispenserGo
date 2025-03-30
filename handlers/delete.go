package handlers

import (
	"MediaDispenserGo/config"
	"MediaDispenserGo/s3client"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["name"]

	// Získání tokenu z hlaviček požadavku
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized: Token required", http.StatusUnauthorized)
		return
	}

	// Validace tokenu a získání izolace
	var allowedPrefix string
	isAllowed := false
	for _, t := range config.AppConfig.Tokens {
		if t.Token == token {
			allowedPrefix = t.Isolation
			isAllowed = true
			break
		}
	}

	// Pokud token není platný (není nalezen v konfiguraci)
	if !isAllowed {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	// Pokud je token izolován, zkontroluj prefix souboru
	if allowedPrefix != "" && !strings.HasPrefix(fileName, allowedPrefix+"/") {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	// Pokračuj v mazání souboru
	err := s3client.Client.RemoveObject(r.Context(), config.AppConfig.S3.Bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Kontrola hlavičky Accept
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJson{
			File:    fileName,
			Message: "File deleted successfully",
		})
	} else {
		// Defaultní textová odpověď
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "File deleted successfully")
	}
}
