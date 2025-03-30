package handlers

import (
	"MediaDispenserGo/config"
	"MediaDispenserGo/s3client"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if config.AppConfig.Dispensing == "none" {
		http.Error(w, "Download disabled", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	fileName := vars["name"]

	object, err := s3client.Client.GetObject(r.Context(), config.AppConfig.S3.Bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer object.Close()

	// Získání vlastností objektu, jako je MIME typ
	objectInfo, err := object.Stat()
	if err != nil {
		http.Error(w, "Object not found or is corrupted", http.StatusInternalServerError)
		return
	}

	// Nastavení správných hlaviček odpovědi
	w.Header().Set("Content-Type", objectInfo.ContentType)
	io.Copy(w, object)
}
