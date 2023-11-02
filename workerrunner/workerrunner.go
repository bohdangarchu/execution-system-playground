package workerrunner

import (
	"app/types"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/xid"
)

const WORKER_PATH = "/home/bohdan/workspace/uni/thesis/worker/main"

func StartProcessWorker(config *types.ProcessIsolationConfig) *types.V8Worker {
	id := xid.New().String()
	socketPath := fmt.Sprintf("/tmp/worker-%s.sock", id)
	println("socket path: ", socketPath)
	// start the worker with the socket path
	cmd := exec.Command(WORKER_PATH, "--socket-path", socketPath)
	// print stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	execErr := cmd.Start()
	if execErr != nil {
		fmt.Println("error starting the worker: ", execErr.Error())
	}
	// need to wait for the process to finish
	// otherwise it will become a zombie process
	go func() {
		cmd.Wait()
		if _, err := os.Stat(socketPath); err == nil {
			os.Remove(socketPath)
		}
	}()
	pid := cmd.Process.Pid
	fmt.Println("pid of the worker: ", pid)
	manager := getCgroup(id, int64(config.MaxMemSize), int64(config.CPUQuota), uint64(config.CPUPeriod))
	// add the pid to the cgroup
	err := manager.AddProc(uint64(pid))
	if err != nil {
		println("error adding process to the cgroup: ", err.Error())
	}
	return &types.V8Worker{
		Id:             id,
		SocketPath:     socketPath,
		ExecutablePath: WORKER_PATH,
		Pid:            pid,
		Cmd:            cmd,
	}
}
