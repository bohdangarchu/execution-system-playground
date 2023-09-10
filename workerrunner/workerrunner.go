package workerrunner

import (
	"app/types"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/xid"
)

func StartV8Worker() *types.V8Worker {
	// generate random id
	id := xid.New().String()
	socketPath := fmt.Sprintf("/tmp/worker-%s.sock", id)
	println("socket path: ", socketPath)
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
	manager := createDefaultCgroup()
	// add the pid to the cgroup
	err := manager.AddProc(uint64(pid))
	if execErr != nil {
		println("error: ", err.Error())
	}
	return &types.V8Worker{
		Id:             id,
		SocketPath:     socketPath,
		ExecutablePath: workerPath,
		Pid:            pid,
		Cmd:            cmd,
	}
}
