package workerrunner

import (
	"app/types"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/containerd/cgroups/v2/cgroup2"
	"github.com/rs/xid"
)

func SendJsonToUnixSocket(socketPath string, jsonSubmission string) (string, error) {
	// Create a custom HTTP client with a Unix domain socket transport
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(proto, addr string) (conn net.Conn, err error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
	// Define the URL for the Unix domain socket (it's not a real URL)
	url := "http://localhost/execute"
	// Create a request body as a bytes.Buffer
	requestBody := bytes.NewBuffer([]byte(jsonSubmission))

	// Create an HTTP GET request
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
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

func StartV8Worker() *types.V8Worker {
	// generate random id
	id := xid.New().String()
	socketPath := fmt.Sprintf("/tmp/worker-%s.sock", id)
	// start the worker with the socket path
	workerPath := "../worker/main"
	cmd := exec.Command(workerPath, "--socket-path", socketPath)
	// print stdout
	cmd.Stdout = os.Stdout

	execErr := cmd.Start()
	if execErr != nil {
		println("error: ", execErr.Error())
	}
	pid := cmd.Process.Pid
	println("pid of the worker: ", pid)
	// add the pid to the cgroup
	manager, err := cgroup2.LoadSystemd("/", "my-cgroup-abc.slice")
	err = manager.AddProc(uint64(pid))
	if execErr != nil {
		println("error: ", err.Error())
	}
	// cmd.Wait()
	return &types.V8Worker{
		SocketPath:     socketPath,
		ExecutablePath: workerPath,
		Pid:            pid,
		Cmd:            cmd,
	}
}
