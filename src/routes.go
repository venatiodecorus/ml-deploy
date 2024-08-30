package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/venatiodecorus/ml-deploy/src/frontend"
	"github.com/venatiodecorus/ml-deploy/src/utils"
)

func handleRequests() {
	// Deprecated? Except /health maybe
	http.HandleFunc("/health", health)
	http.HandleFunc("/dockergen", dockerHandler)
	http.HandleFunc("/terraform", terraformHandler)
	http.HandleFunc("/deploy", deployHandler)
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Serve pages & components
	frontend.RegisterRoutes()
	// User input
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/destroy", destroyHandler)
	// Start 'er up
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	type RequestData struct {
		Instructions string `json:"instructions"`
	}
	var reqData RequestData
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	var plan string
	plan, err = update(reqData.Instructions)
	if err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	ret := APIResponse{Output: plan}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}

func destroyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Failed to read request body", http.StatusBadRequest)
	// 	return
	// }
	// defer r.Body.Close()

	// type RequestData struct {
	// 	Instructions string `json:"instructions"`
	// }
	// var reqData RequestData
	// err = json.Unmarshal(body, &reqData)
	// if err != nil {
	// 	http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
	// 	return
	// }

	err := utils.Destroy()
	if !err {
		http.Error(w, "Failed to destroy", http.StatusInternalServerError)
		return
	}

	ret := APIResponse{Output: "Infra destroyed"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
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
	// TODO Fix return type and response
	docker(query)

	ret := APIResponse{Output: "success"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}

// func dockerListHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	list,err := dockerList();
// 	if err != nil {
// 		http.Error(w, "Failed to list docker images", http.StatusInternalServerError)
// 		return
// 	}

// 	// ret := APIResponse{Output: "success"}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(list)
// }

func terraformHandler(w http.ResponseWriter, r *http.Request) {
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
	resp, err := terraformRequest(query)

	ret := APIResponse{Output: resp}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}

// Run the entire deployment process based on the instructions provided
func deployHandler(w http.ResponseWriter, r *http.Request) {
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
	deploy := hetznerDeploy(query)

	// TODO Proper error type here?
	if !deploy {
		http.Error(w, "Failed to deploy", http.StatusInternalServerError)
		return
	}

	ret := APIResponse{Output: "Deployment successful"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}