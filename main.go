package main

import (
	"app/api"
	"app/firerunner"
	"app/types"
	"app/utils"
	"flag"
	"fmt"
)

func main() {
	// performance.MeasureStartupTimes()
	runServer()
}

func runServer() {
	pathPtr := flag.String("path", "config.json", "Path to the config")
	flag.Parse()
	config := utils.LoadConfig(*pathPtr)
	api.Run(&config)
}

func runVM() {
	vm, err := firerunner.StartVM(true, &types.FirecrackerConfig{
		MemSizeMib: 128,
		CPUQuota:   200000,
		CPUPeriod:  1000000,
	}, true)
	defer vm.StopVMandCleanUp()
	if err != nil {
		panic(fmt.Sprintf("Failed to start VM: %v", err))
	}
	res, err := firerunner.RunSubmissionInsideVM(vm, utils.JsonSubmission)
	if err != nil {
		panic(fmt.Sprintf("Failed to run submission inside VM: %v", err))
	}
	fmt.Println(res)
}
