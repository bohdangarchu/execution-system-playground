package main

import (
	"app/api"
	"app/firerunner"
	"app/types"
	"flag"
	"log"
)

var allowedImplValues = []string{"docker", "firecracker", "process"}

func main() {
	// performance.MeasureStartupTimes()
	runServer()
}

func runServer() {
	pathPtr := flag.String("path", "config.json", "Path to the config")
	flag.Parse()
	config := LoadConfig(*pathPtr)
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
		log.Fatalf("Failed to start VM: %v", err)
	}
	res, err := firerunner.RunSubmissionInsideVM(vm, jsonSubmission)
	if err != nil {
		log.Fatalf("Failed to run submission inside VM: %v", err)
	}
	log.Println(res)
}
