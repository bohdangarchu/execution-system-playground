package performance

import (
	"app/api"
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/xid"
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

type StringFunction func(string, string) (string, error)

func EndToEndExecutionTime() {
	// is not up to date
	fmt.Println("Measuring firecracker execution time...")
	api.Run(&types.Config{})
	startTime := time.Now()
	out, err := SendSubmissionToUrl(jsonSubmission, "http://localhost:8081")
	if err != nil {
		panic(err)
	}
	executionTime := time.Since(startTime)
	fmt.Println("output: ", out)
	fmt.Println("execution time: ", executionTime)

	fmt.Println("Measuring docker execution time...")
	api.Run(&types.Config{})
	startTime = time.Now()
	out, err = SendSubmissionToUrl(jsonSubmission, "http://localhost:8081")
	if err != nil {
		panic(err)
	}
	executionTime = time.Since(startTime)
	fmt.Println("output: ", out)
	fmt.Println("execution time: ", executionTime)
}

func ExecuteWithTime(input1 string, input2 string, fn StringFunction) (string, error, time.Duration) {
	startTime := time.Now()
	output, err := fn(input1, input2)
	executionTime := time.Since(startTime)
	return output, err, executionTime
}

func TimeDockerStartupAndSubmission() error {
	// old
	startTime := time.Now()
	config := &types.DockerConfig{
		MaxMemSize: 10000000,
		CPUQuota:   125000,
		CPUPeriod:  1000000,
	}
	// time the execution
	dockerContainer, err := docrunner.StartExecutionServerInDocker("8080", config)
	if err != nil {
		return err
	}
	defer killContainerAndGetLogs(dockerContainer)
	executionTime := time.Since(startTime)
	fmt.Println("Execution Server started in Docker in: ", executionTime)

	// sometimes the docker container is not ready to receive requests
	// time.Sleep(50 * time.Millisecond)

	outputDocker, errDocker, timeDocker := ExecuteWithTime(jsonSubmission, "", SendSubmissionToUrl)
	if errDocker != nil {
		panic(errDocker)
	}
	fmt.Println("output docker: ", outputDocker)
	fmt.Println("execution time docker: ", timeDocker)

	return err
}

func SendSubmissionToUrl(jsonSubmission string, url string) (string, error) {
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

func TestCopying() {
	for i := 0; i < 15; i++ {
		id := xid.New().String()
		go firerunner.GetUniqueDrive(id)
	}
	time.Sleep(60 * time.Second)
}
