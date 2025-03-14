package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type URLMapping struct {
	LongURL string `json:"long_url"`
}

var urlStore = make(map[string]string) // In-memory store

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generates a random string of length 6 - Vishal did this change
func generateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

// Handler to create a shortened URL
func createShortURL(w http.ResponseWriter, r *http.Request) {
	var request URLMapping
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Generate a unique short code
	shortCode := generateShortCode()
	urlStore[shortCode] = request.LongURL

	response := map[string]string{"short_url": fmt.Sprintf("http://localhost:8080/%s", shortCode)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler to redirect to original URL
func redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["code"]

	if originalURL, ok := urlStore[shortCode]; ok {
		http.Redirect(w, r, originalURL, http.StatusFound)
	} else {
		http.Error(w, "URL not found", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/shorten", createShortURL).Methods("POST")
	r.HandleFunc("/{code}", redirectURL).Methods("GET")

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
