package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/workerrunner"
	"fmt"
	"time"
)

func monitorContainerHealth(containerPool chan types.DockerContainer, config *types.Config) {
	for {
		container := <-containerPool
		fmt.Printf("health check, container %s\n", container.Port)
		healthy := docrunner.CheckContainerHealth(&container)
		if healthy {
			containerPool <- container
		} else {
			fmt.Printf("container %s is not healthy, killing it\n", container.Port)
			docrunner.CleanUp(&container, true)
			newContainer, err := docrunner.StartExecutionServerInDocker(
				container.Port,
				int64(config.Docker.MaxMemSize),
				int64(config.Docker.NanoCPUs),
			)
			if err != nil {
				fmt.Printf("failed to start docker container: %v\n", err)
				continue
			}
			containerPool <- *newContainer
		}
		time.Sleep(1 * time.Second)
	}
}

func monitorVMHealth(vmPool chan types.FirecrackerVM, config *types.Config) {
	for {
		vm := <-vmPool
		fmt.Printf("health check, vm %s\n", vm.Ip.String())
		healthy := firerunner.CheckVMHealth(&vm)
		if healthy {
			vmPool <- vm
		} else {
			fmt.Printf("vm %s is not healthy, killing it\n", vm.VmmID)
			vm.StopVMandCleanUp()
			newVM, err := firerunner.StartVM(true, config.Firecracker)
			if err != nil {
				fmt.Printf("failed to start vm: %v\n", err)
				continue
			}
			vmPool <- *newVM
		}
		time.Sleep(1 * time.Second)
	}
}

func monitorV8Worker(workerPool chan types.V8Worker, config *types.ProcessIsolationConfig) {
	for {
		worker := <-workerPool
		fmt.Printf("health check, worker %s\n", worker.Id)
		healthy := workerrunner.CheckWorkerHealth(&worker)
		if healthy {
			workerPool <- worker
		} else {
			fmt.Printf("worker %s is not healthy, killing it\n", worker.Id)
			workerrunner.KillWorker(&worker)
			newWorker := workerrunner.StartProcessWorker(
				config.CgroupMaxMem,
				config.CgroupMaxCPU,
			)
			workerPool <- *newWorker
		}
		time.Sleep(1 * time.Second)
	}
}
