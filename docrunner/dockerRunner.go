package docrunner

import (
	"app/types"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func StartExecutionServerInDocker(port string, config *types.DockerConfig) (*types.DockerContainer, error) {
	// starts a docker container with the image "execution-server"
	fmt.Println("Starting docker container...")
	// Create a background context
	ctx := context.Background()
	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close()
	if err != nil {
		return nil, err
	}
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "execution-server",
			ExposedPorts: nat.PortSet{
				"8080/tcp": struct{}{},
			},
			Tty: false,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8080/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: port,
					},
				},
			},
			Resources: container.Resources{
				Memory:    int64(config.MaxMemSize),
				CPUQuota:  int64(config.CPUQuota),
				CPUPeriod: int64(config.CPUPeriod),
			},
		}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, dockertypes.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	container, _ := cli.ContainerInspect(ctx, resp.ID)
	realPort := container.NetworkSettings.Ports["8080/tcp"][0].HostPort
	return &types.DockerContainer{
		ContainerId: resp.ID,
		Port:        realPort,
		Cli:         cli,
		Ctx:         ctx,
	}, nil
}

func SendJSONSubmissionToDocker(port string, jsonSubmission string) (string, error) {
	url := "http://localhost:" + port + "/execute"

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
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad response status code: %d, resp: %s", resp.StatusCode, body)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(responseBody), nil
}

func CleanUp(dockerContainer *types.DockerContainer, debug bool) {
	// kill the container
	err := KillDockerContainer(dockerContainer)
	if err != nil {
		fmt.Println("failed to kill container: ", err)
	}
	if debug {
		// Retrieve the logs of the container
		logs, err := RetrieveLogsFromDockerContainer(dockerContainer)
		if err != nil {
			fmt.Println("failed to retrieve logs: ", err)
		}
		fmt.Printf("Logs from container with port %s\n", dockerContainer.Port)
		fmt.Println(logs)
		fmt.Println("--------------------------------------------------")
	}
}

func RetrieveLogsFromDockerContainer(dockerContainer *types.DockerContainer) (string, error) {
	// Retrieve the logs of the container
	out, err := dockerContainer.Cli.ContainerLogs(
		dockerContainer.Ctx, dockerContainer.ContainerId,
		dockertypes.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	bytes, err := ioutil.ReadAll(out)
	return string(bytes), err
}

func KillDockerContainer(dockerContainer *types.DockerContainer) error {
	return dockerContainer.Cli.ContainerKill(dockerContainer.Ctx, dockerContainer.ContainerId, "SIGKILL")
}

func StartContainerAndRunSubmission(jsonSubmission string) (string, error) {
	config := &types.DockerConfig{
		MaxMemSize: 10000000,
		CPUQuota:   125000,
		CPUPeriod:  1000000,
	}
	dockerContainer, err := StartExecutionServerInDocker("8080", config)
	if err != nil {
		return "", err
	}
	defer CleanUp(dockerContainer, true)
	WaitUntilAvailable(dockerContainer)

	res, err := SendJSONSubmissionToDocker("8080", jsonSubmission)
	if err != nil {
		return "", err
	}
	return res, nil
}
