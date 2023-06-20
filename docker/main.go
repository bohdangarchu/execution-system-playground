package main

import (
	"docker/executor"
	"docker/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the request JSON body
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
		return
	}

	// Execute the JavaScript code
	outputArray := executor.RunFunctionWithInputs(functionSubmission)

	results := make([]types.TestResult, len(functionSubmission.TestCases))
	status := ""

	for i, testCase := range functionSubmission.TestCases {
		status = "Fail"
		if outputArray[i].Value == testCase.ExpectedOutput {
			status = "Pass"
		}
		results[i] = types.TestResult{
			TestCase:     testCase,
			ActualOutput: outputArray[i],
			Status:       status,
		}
	}

	// Convert the result to JSON
	responseJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to convert result to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
