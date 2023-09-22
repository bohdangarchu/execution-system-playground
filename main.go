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
	pathPtr := flag.String("path", "config.json", "Path to the config")
	flag.Parse()
	config := LoadConfig(*pathPtr)
	runServer(&config)
}

func runServer(config *types.Config) {
	api.Run(config.Isolation, config.Workers)
}

func runVM() {
	vm, err := firerunner.StartVM(true)
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
