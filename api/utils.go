package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/workerrunner"
	"fmt"
	"time"
)

const INTERVAL = 1 * time.Second

func monitorContainerHealth(containerPool chan types.DockerContainer, config *types.Config) {
	for {
		container := <-containerPool
		healthy := docrunner.CheckContainerHealth(&container)
		if healthy {
			containerPool <- container
		} else {
			fmt.Printf("container %s is not healthy, killing it\n", container.Port)
			docrunner.CleanUp(&container, true)
			newContainer, err := docrunner.StartExecutionServerInDocker(
				container.Port,
				config.Docker,
			)
			if err != nil {
				fmt.Printf("failed to start docker container: %v\n", err)
				continue
			}
			containerPool <- *newContainer
		}
		time.Sleep(INTERVAL)
	}
}

func monitorVMHealth(vmPool chan types.FirecrackerVM, config *types.Config, useDefaultDrive bool, debug bool) {
	for {
		vm := <-vmPool
		healthy := firerunner.CheckVMHealth(&vm)
		if healthy {
			vmPool <- vm
		} else {
			fmt.Printf("vm %s is not healthy, killing it\n", vm.VmmID)
			vm.StopVMandCleanUp()
			newVM, err := firerunner.StartVM(useDefaultDrive, config.Firecracker, debug)
			if err != nil {
				fmt.Printf("failed to start vm: %v\n", err)
				continue
			}
			vmPool <- *newVM
		}
		time.Sleep(INTERVAL)
	}
}

func monitorV8Worker(workerPool chan types.V8Worker, config *types.ProcessIsolationConfig) {
	for {
		worker := <-workerPool
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
		time.Sleep(INTERVAL)
	}
}
