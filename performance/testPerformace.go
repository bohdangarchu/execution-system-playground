package performance

import (
	"app/docrunner"
	"app/types"
	"app/v8runner"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var jsonSubmission = `
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

type StringFunction func(string) (string, error)

func TestPerformance() {
	// compares v8 with v8 in docker
	outputV8, errV8, timeV8 := ExecuteWithTime(jsonSubmission, executeJSONSubmissionUsingV8)
	if errV8 != nil {
		panic(errV8)
	}
	fmt.Println("output v8: ", outputV8)
	fmt.Println("time v8: ", timeV8)
	outputDocker, errDocker, timeDocker := ExecuteWithTime(jsonSubmission, ExecuteJSONSubmissionUsingDocker)
	if errDocker != nil {
		panic(errDocker)
	}
	fmt.Println("output docker: ", outputDocker)
	fmt.Println("time docker: ", timeDocker)
}

func ExecuteWithTime(input string, fn StringFunction) (string, error, time.Duration) {
	startTime := time.Now()
	output, err := fn(input)
	executionTime := time.Since(startTime)
	return output, err, executionTime
}

func TimeDockerStartupAndSubmission() error {
	startTime := time.Now()
	// time the execution
	dockerContainer, err := docrunner.StartExecutionServerInDocker()
	if err != nil {
		return err
	}
	defer killContainerAndGetLogs(dockerContainer)
	executionTime := time.Since(startTime)
	fmt.Println("Execution Server started in Docker in: ", executionTime)

	// sometimes the docker container is not ready to receive requests
	// time.Sleep(50 * time.Millisecond)

	outputDocker, errDocker, timeDocker := ExecuteWithTime(jsonSubmission, ExecuteJSONSubmissionUsingDocker)
	if errDocker != nil {
		panic(errDocker)
	}
	fmt.Println("output docker: ", outputDocker)
	fmt.Println("execution time docker: ", timeDocker)

	return err
}

func ExecuteJSONSubmissionUsingDocker(jsonSubmission string) (string, error) {
	url := "http://localhost:8080"

	// Create a request body as a bytes.Buffer
	requestBody := bytes.NewBuffer([]byte(jsonSubmission))

	// Make the POST request
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response status code: %d", resp.StatusCode)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(responseBody), nil
}

func executeJSONSubmissionUsingV8(jsonSubmission string) (string, error) {
	var functionSubmission types.FunctionSubmission
	err := json.Unmarshal([]byte(jsonSubmission), &functionSubmission)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json: %v", err)
	}

	// Execute the JavaScript code
	outputArray := v8runner.RunFunctionWithInputs(functionSubmission)

	// Convert the result to JSON
	responseJSON, err := json.Marshal(outputArray)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %v", err)
	}
	return string(responseJSON), nil
}

func killContainerAndGetLogs(dockerContainer *types.DockerContainer) {
	// kill the container
	err := docrunner.KillDockerContainer(dockerContainer)
	if err != nil {
		panic(err)
	}

	// Retrieve the logs of the container
	logs, err := docrunner.RetrieveLogsFromDockerContainer(dockerContainer)
	if err != nil {
		panic(err)
	}
	fmt.Println("logs: ", logs)
}
