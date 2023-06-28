package performance

import (
	"app/types"
	"app/v8runner"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type StringFunction func(string) string

func TestPerformance() {
	// compares v8 with v8 in docker
	jsonSubmission := `
	{
		"functionName": "addTwoNumbers",
		"code": "function addTwoNumbers(a, b) {\n  return a + b;\n}",
		"testCases": [
		  {
			"input": [
			  {
				"value": 3,
				"type": "number"
			  },
			  {
				"value": -10,
				"type": "number"
			  }
			]
		  }
		]
	  }
	`
	outputV8, timeV8 := ExecuteWithTime(jsonSubmission, executeJSONSubmissionUsingV8)
	fmt.Println("output v8: ", outputV8)
	fmt.Println("time v8: ", timeV8)
	outputDocker, timeDocker := ExecuteWithTime(jsonSubmission, executeJSONSubmissionUsingDocker)
	fmt.Println("output docker: ", outputDocker)
	fmt.Println("time docker: ", timeDocker)

}

func ExecuteWithTime(input string, fn StringFunction) (string, time.Duration) {
	startTime := time.Now()
	output := fn(input)
	executionTime := time.Since(startTime)
	return output, executionTime
}

func executeJSONSubmissionUsingDocker(jsonSubmission string) string {
	url := "http://localhost:8080"

	// Create a request body as a bytes.Buffer
	requestBody := bytes.NewBuffer([]byte(jsonSubmission))

	// Make the POST request
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		panic(fmt.Sprintf("Error making POST request: %v", err))
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Received non-OK response: %v", resp.Status))
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Error reading response body: %v", err))
	}

	return string(responseBody)
}

func executeJSONSubmissionUsingV8(jsonSubmission string) string {
	var functionSubmission types.FunctionSubmission
	err := json.Unmarshal([]byte(jsonSubmission), &functionSubmission)
	if err != nil {
		panic(fmt.Sprintf("Error decoding JSON: %v", err))
	}

	// Execute the JavaScript code
	outputArray := v8runner.RunFunctionWithInputs(functionSubmission)

	// Convert the result to JSON
	responseJSON, err := json.Marshal(outputArray)
	if err != nil {
		panic(fmt.Sprintf("failed to convert result to JSON: %v", err))
	}
	return string(responseJSON)
}
