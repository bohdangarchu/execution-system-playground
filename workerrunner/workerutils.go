package workerrunner

import (
	"app/types"
	"app/utils"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"syscall"
	"time"

	"github.com/containerd/cgroups/v2/cgroup2"
)

const CGROUP_NAME = "worker3.slice"

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
		fmt.Printf("Response: %v", resp)
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

func getCgroup(id string, maxMem, cpuQuota int64, cpuPeriod uint64) *cgroup2.Manager {
	name := "mycgroup-" + id + ".slice"
	fmt.Printf("cgroup name: %s\n", name)
	return utils.CreateCgroup(name, maxMem, cpuQuota, cpuPeriod)
}

func IsProcessRunning(pid int) bool {
	// Check if the process exists by sending signal 0 to the given PID
	err := syscall.Kill(pid, 0)
	if err == nil || err == syscall.EPERM {
		return true
	}
	return false
}

func IsWorkerRunning(worker *types.ProcessWorker) bool {
	finished := make(chan error, 1)
	go func() {
		err := worker.Cmd.Wait()
		finished <- err
	}()
	select {
	case <-finished:
		return false
	case <-time.After(10 * time.Millisecond):
		return true
	}
}

func KillWorker(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

func CheckWorkerHealth(worker *types.ProcessWorker) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(proto, addr string) (conn net.Conn, err error) {
				return net.Dial("unix", worker.SocketPath)
			},
		},
	}
	url := "http://localhost/health"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func WaitUntilAvailable(worker *types.ProcessWorker) {
	for {
		if CheckWorkerHealth(worker) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
