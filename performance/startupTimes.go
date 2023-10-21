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
	vm, err := firerunner.StartVM(false, config)
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
	measureAvgFirecrackerStartupTime()
}
