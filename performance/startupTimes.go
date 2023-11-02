package performance

import (
	"app/firerunner"
	"app/types"
	"fmt"
	"time"
)

func measureFirecrackerStartupTime() time.Duration {
	fmt.Println("Measuring firecracker startup time...")
	config := &types.FirecrackerConfig{
		CPUCount:   1,
		MemSizeMib: 128,
	}
	startTime := time.Now()
	vm, err := firerunner.StartVM(false, config, false)
	if err != nil {
		fmt.Printf("Failed to start VM: %v\n", err)
	}
	startupTime := time.Since(startTime)
	fmt.Println("startup time: ", startupTime)
	fmt.Println("IP: ", vm.Ip)
	vm.StopVMandCleanUp()
	return startupTime
}

func measureAvgFirecrackerStartupTime() {
	avgStartupTime := 0 * time.Millisecond
	n := 5
	for i := 0; i < n; i++ {
		avgStartupTime += measureFirecrackerStartupTime()
	}
	fmt.Printf("Average firecracker startup time: %d milliseconds\n", int(avgStartupTime.Milliseconds())/n)
}

func MeasureStartupTimes() {
	// measureFirecrackerStartupTime()
}

// func MeasureDockerStartupTimes() {
// 	n := 1
// 	containerChannel := make(chan *types.DockerContainer, n)
// 	config := &types.DockerConfig{
// 		MaxMemSize: 10000000,
// 		// 10 ^ 9 = 1 CPU
// 		NanoCPUs: 1000000000,
// 	}
// 	startTime := time.Now()
// 	for i := 0; i < n; i++ {
// 		go func() {
// 			container, err := docrunner.StartExecutionServerInDocker(
// 				// with 0 docker will pick an available port
// 				"0",
// 				int64(config.MaxMemSize),
// 				int64(config.NanoCPUs),
// 			)
// 			containerChannel <- container
// 			if err != nil {
// 				fmt.Printf("Failed to start docker container: %v\n", err)
// 			}
// 		}()
// 	}
// 	for i := 0; i < n; i++ {
// 		container := <-containerChannel
// 		docrunner.WaitUntilAvailable(container)
// 	}
// 	startupTime := time.Since(startTime)
// 	fmt.Println("startup time: ", startupTime)
// }
