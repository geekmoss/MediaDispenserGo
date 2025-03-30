package main

import (
	"MediaDispenserGo/config"
	"MediaDispenserGo/handlers"
	"MediaDispenserGo/middlewares"
	"MediaDispenserGo/s3client"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Načtení konfigurace
	config.LoadConfig()

	// Inicializace S3 klienta
	s3client.Initialize()

	// Nastavení směrovače
	router := mux.NewRouter()

	// Definice endpointů
	router.HandleFunc("/u", middlewares.Authenticate(handlers.UploadHandler, middlewares.OperationUpload)).Methods("POST")
	router.HandleFunc("/g/{name:.*}", middlewares.Authenticate(handlers.DownloadHandler, middlewares.OperationGet)).Methods("GET")
	router.HandleFunc("/d/{name:.*}", middlewares.Authenticate(handlers.DeleteHandler, middlewares.OperationDelete)).Methods("DELETE")

	// Spuštění serveru
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
