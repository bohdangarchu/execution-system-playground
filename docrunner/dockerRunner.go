package docrunner

import (
	"app/types"
	"context"
	"fmt"
	"io/ioutil"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func RunExecutionServerInDocker() error {
	dockerContainer, err := StartExecutionServerInDocker()

	if err != nil {
		return err
	}

	// Wait for the container to stop
	statusCh, errCh := dockerContainer.Cli.ContainerWait(
		dockerContainer.Ctx,
		dockerContainer.ContainerId,
		container.WaitConditionNotRunning,
	)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	// Retrieve the logs of the container
	out, err := dockerContainer.Cli.ContainerLogs(
		dockerContainer.Ctx, dockerContainer.ContainerId,
		dockertypes.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	bytes, err := ioutil.ReadAll(out)
	fmt.Println("Execution Server Logs:")
	fmt.Println(string(bytes))
	return err
}

func StartExecutionServerInDocker() (*types.DockerContainer, error) {
	// starts a docker container with the image "execution-server"
	// and exposes port 8080

	// Create a background context
	ctx := context.Background()

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close() // Close the Docker client when function returns
	if err != nil {
		return nil, err
	}

	// Create a container with a volume mount to write the script file inside the container
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
						HostPort: "8080",
					},
				},
			},
		}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, dockertypes.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return &types.DockerContainer{
		ContainerId: resp.ID,
		Cli:         cli,
		Ctx:         ctx,
	}, nil
}

func RunJsInDocker(jsCode string) (string, error) {
	// creates a js file inside docker using shell

	// Create a background context
	ctx := context.Background()

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close() // Close the Docker client when function returns
	if err != nil {
		return "", err
	}

	// Prepare the JavaScript code as a script file
	scriptFile := "/app/script.js"

	// Create a container with a volume mount to write the script file inside the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "node",
		Cmd: []string{
			"sh", "-c", fmt.Sprintf(
				"mkdir -p /app && echo \"%s\" > %s && node %s", jsCode, scriptFile, scriptFile)},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, dockertypes.ContainerStartOptions{}); err != nil {
		return "", err
	}

	// Wait for the container to stop
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	// Retrieve the logs of the container
	out, err := cli.ContainerLogs(ctx, resp.ID, dockertypes.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	bytes, err := ioutil.ReadAll(out)
	return string(bytes), err
}
