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
		MemSizeMib: 128,
		CPUQuota:   200000,
		CPUPeriod:  1000000,
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
