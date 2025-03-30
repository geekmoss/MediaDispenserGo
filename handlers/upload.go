package handlers

import (
	"MediaDispenserGo/config"
	"MediaDispenserGo/middlewares"
	"MediaDispenserGo/plugins"
	"MediaDispenserGo/s3client"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Načítáme aktivní pluginy dynamicky
	var activePlugins []plugins.FileProcessor
	for _, pluginName := range config.AppConfig.Plugins {
		plugin, err := plugins.GetPluginInstance(pluginName, middlewares.OperationUpload)
		if err != nil {
			if err.Error() == "plugin are not allowed in this operation" {
				continue
			}
			log.Printf("Failed to initialize plugin %s: %v", pluginName, err)
			continue
		}
		activePlugins = append(activePlugins, plugin)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	name := filepath.Base(header.Filename)[:len(header.Filename)-len(ext)]
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	var processedFile bytes.Buffer
	_, err = processedFile.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	if len(activePlugins) > 0 {
		for _, plugin := range activePlugins {
			processedData, pluginMimeType, pluginFileName, pluginExtension, err := plugin.Process(processedFile, mimeType, name, ext)
			if err != nil {
				log.Println("Plugin error ", plugin.GetName(), ":", err)
				continue
			}
			processedFile = *bytes.NewBuffer(processedData)
			mimeType = pluginMimeType
			ext = pluginExtension
			name = pluginFileName
		}
	}

	options := minio.PutObjectOptions{
		ContentType: mimeType,
	}

	newFilename := fmt.Sprintf("%s%s", name, ext)

	token := r.Header.Get("Authorization")
	objectName := newFilename
	if config.AppConfig.Tokens != nil {
		for _, t := range config.AppConfig.Tokens {
			if t.Token == token && t.Isolation != "" {
				objectName = fmt.Sprintf("%s/%s", t.Isolation, newFilename)
				break
			}
		}
	}

	_, err = s3client.Client.PutObject(
		r.Context(),
		config.AppConfig.S3.Bucket,
		objectName,
		bytes.NewReader(processedFile.Bytes()),
		int64(processedFile.Len()),
		options,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJson{
			File:    objectName,
			Message: "File uploaded successfully",
		})
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, objectName)
	}
}
