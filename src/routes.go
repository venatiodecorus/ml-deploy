package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleRequests() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.HandleFunc("/health", health)
	http.HandleFunc("/docker", dockerHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

type APIResponse struct {
	Output string `json:"output"`
}

func dockerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Define a struct to hold the JSON data
	type RequestData struct {
		Instructions string `json:"instructions"`
	}

	// Unmarshal the JSON data into the struct
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Access the query value from the JSON data
	query := requestData.Instructions

	// Use the query value as needed
	resp := dockerRequest(query)

	ret := APIResponse{Output: resp}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}