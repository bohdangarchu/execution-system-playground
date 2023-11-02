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

func CreateCgroup(name string, maxMem, cpuQuota int64, cpuPeriod uint64) *cgroup2.Manager {
	// TODO delete cgroup after usage
	resources := &cgroup2.Resources{
		Memory: &cgroup2.Memory{
			Max: &maxMem,
		},
		CPU: &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(&cpuQuota, &cpuPeriod),
		},
	}
	manager, err := cgroup2.NewSystemd("/", name, -1, resources)
	if err != nil {
		println("error creating a cgroup: ", err.Error())
	}
	return manager
}

func CreateCPUCgroup(name string, cpuQuota int64, cpuPeriod uint64) *cgroup2.Manager {
	resources := &cgroup2.Resources{
		CPU: &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(&cpuQuota, &cpuPeriod),
		},
	}
	manager, err := cgroup2.NewSystemd("/", name, -1, resources)
	if err != nil {
		println("error creating a cgroup: ", err.Error())
	}
	return manager
}

func getCgroup(id string, maxMem, cpuQuota int64, cpuPeriod uint64) *cgroup2.Manager {
	name := "mycgroup-" + id + ".slice"
	fmt.Printf("cgroup name: %s\n", name)
	return CreateCgroup(name, maxMem, cpuQuota, cpuPeriod)
}

func IsProcessRunning(pid int) bool {
	// Check if the process exists by sending signal 0 to the given PID
	err := syscall.Kill(pid, 0)
	if err == nil || err == syscall.EPERM {
		return true
	}
	return false
}

func IsWorkerRunning(worker *types.V8Worker) bool {
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

func KillWorker(worker *types.V8Worker) {
	worker.Cmd.Process.Signal(os.Interrupt)
	// check if the socket file exists
	// if it does, remove it
	if _, err := os.Stat(worker.SocketPath); err == nil {
		os.Remove(worker.SocketPath)
	}
}

func CheckWorkerHealth(worker *types.V8Worker) bool {
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

func WaitUntilAvailable(worker *types.V8Worker) {
	for {
		if CheckWorkerHealth(worker) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
