package middlewares

import (
	"MediaDispenserGo/config"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Operation string

const (
	OperationList   Operation = "list"
	OperationGet    Operation = "get"
	OperationUpload Operation = "upload"
	OperationDelete Operation = "delete"
)

// ResponseJson je pomocná funkce pro odpověď ve formátu JSON
func ResponseJson(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// Authenticate ověřuje token pouze tam, kde je to nutné
func Authenticate(next http.HandlerFunc, operation Operation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Zjištění odpovídajícího pravidla na základě prefixu
		dispensingMode := config.AppConfig.Dispensing

		// Pokud se jedná o root, tedy bez prefixu
		if strings.Count(r.URL.Path, "/") == 2 && config.AppConfig.DispensingRoot != "" {
			dispensingMode = config.AppConfig.DispensingRoot
		} else {
			// Hledá se dle prefixu zda je třeba přepsat pravidlo
			for _, prefixRule := range config.AppConfig.DispensingPrefix {
				if strings.HasPrefix(r.URL.Path, "/g/"+prefixRule.Prefix) {
					dispensingMode = prefixRule.Mode
					break
				}
			}
		}
		log.Println(r.URL.Path)
		// Zohlednění typu výdeje (dispensing)
		if strings.HasPrefix(r.URL.Path, "/g") {
			switch dispensingMode {
			case config.DispensingPublic:
				// Public režim - token není potřeba
				next.ServeHTTP(w, r)
				return
			case config.DispensingNone:
				// Zákaz výdeje - žádný požadavek nebude povolen
				writeError(w, r, http.StatusForbidden, "Download disabled")
				return
			}
		}

		// Pro všechny jiné případy ověřujeme autorizaci
		token := r.Header.Get("Authorization")
		if token == "" {
			writeError(w, r, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Kontrola, zda token existuje a není omezen
		for _, t := range config.AppConfig.Tokens {
			if t.Token == token {
				// Disallow kontrola
				for _, dis := range t.Disallow {
					if Operation(dis) == operation {
						writeError(w, r, http.StatusForbidden, "Permission denied")
						return
					}
				}
				// Pokračujeme, protože token je platný
				next.ServeHTTP(w, r)
				return
			}
		}

		// Pokud není token platný, vracíme Unauthorized
		writeError(w, r, http.StatusUnauthorized, "Invalid token")
	}
}

// writeError je pomocná funkce pro zápis odpovědi dle hlavičky Accept
func writeError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		ResponseJson(w, statusCode, message)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		w.Write([]byte(message))
	}
}
