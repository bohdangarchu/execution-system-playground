package api

import (
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

func handleRequest(w http.ResponseWriter, r *http.Request) {
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

func Run() {
	http.HandleFunc("/", handleRequest)
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
