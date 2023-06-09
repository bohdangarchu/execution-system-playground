package docrunner

import (
	"app/types"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

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
	// to not wait for the container to exit gracefully
	noWaitTimeout := 0
	return dockerContainer.Cli.ContainerStop(
		dockerContainer.Ctx,
		dockerContainer.ContainerId,
		container.StopOptions{Timeout: &noWaitTimeout},
	)
}

func StartExecutionServerInDocker() (*types.DockerContainer, error) {
	// starts a docker container with the image "execution-server"
	// and exposes port 8080

	// Create a background context
	ctx := context.Background()

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// TODO cli is used later so maybe don't defer close
	defer cli.Close() // Close the Docker client when function returns
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

	waitForContainerRunning(cli, resp.ID)

	return &types.DockerContainer{
		ContainerId: resp.ID,
		Cli:         cli,
		Ctx:         ctx,
	}, nil
}

func waitForContainerRunning(cli *client.Client, containerID string) {
	ctx := context.Background()
	filterArgs := filters.NewArgs()
	filterArgs.Add("id", containerID)
	filterArgs.Add("status", "running")
	options := dockertypes.ContainerListOptions{Filters: filterArgs}

	for {
		containers, err := cli.ContainerList(ctx, options)
		if err != nil {
			panic(err)
		}

		if len(containers) > 0 {
			break
		}

		// Sleep for a short duration before checking again
		time.Sleep(100 * time.Millisecond)
	}
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
