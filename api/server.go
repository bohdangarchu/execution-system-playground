package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Executor interface {
	RunFunctionSubmission(functionSubmission types.FunctionSubmission) []string
}

func handleRequestWithFirecracker(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
		return
	}

	jsonSubmission, err := json.Marshal(functionSubmission)
	responseString := firerunner.RunSubmissionInsideVM(string(jsonSubmission))
	responseJSON := []byte(responseString)

	// Convert the result to JSON
	// responseJSON, err := json.Marshal(outputArray)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to convert result to JSON: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

func Run(option string) {
	if option == "docker" {
		http.HandleFunc("/", handleRequestWithDocker)
	} else if option == "firecracker" {
		http.HandleFunc("/", handleRequestWithFirecracker)
	} else {
		http.HandleFunc("/", handleRequestWithV8)
	}
	log.Println("Listening on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleRequestWithDocker(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
		return
	}

	jsonSubmission, err := json.Marshal(functionSubmission)
	responseString, err := docrunner.RunSubmissionInsideDocker(string(jsonSubmission))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
		return
	}
	responseJSON := []byte(responseString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func handleRequestWithV8(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
		return
	}

	jsonSubmission, err := json.Marshal(functionSubmission)
	responseString := firerunner.RunSubmissionInsideVM(string(jsonSubmission))
	responseJSON := []byte(responseString)

	// Convert the result to JSON
	// responseJSON, err := json.Marshal(outputArray)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to convert result to JSON: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}
