package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/v8runner"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	firecracker = iota
	docker
	v8
)

func handleRequestWithFirecracker(w http.ResponseWriter, r *http.Request) {
	// get json string from request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	jsonSubmission := buf.String()

	// var functionSubmission types.FunctionSubmission
	// err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
	// 	log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
	// 	return
	// }
	// jsonSubmission, err := json.Marshal(functionSubmission)
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

func Run(option int) {
	if option == docker {
		http.HandleFunc("/", handleRequestWithDocker)
	} else if option == firecracker {
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

	outputArray := v8runner.RunFunctionWithInputs(functionSubmission)

	// Convert the result to JSON
	responseJSON, err := json.Marshal(outputArray)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to convert result to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}
