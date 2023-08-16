package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Run(option string, workers int) {
	var vmPool chan types.FirecrackerVM
	var containerPool chan types.DockerContainer
	if option == "docker" {
		containerPool = make(chan types.DockerContainer, workers)
		port := 8081
		for i := 0; i < workers; i++ {
			container, err := docrunner.StartExecutionServerInDocker(fmt.Sprintf("%d", port))
			if err != nil {
				log.Fatalf("Failed to start docker container: %v", err)
			}
			containerPool <- *container
			port++
		}
		http.HandleFunc("/execute", getDockerHandler(containerPool))
	} else if option == "firecracker" {
		vmPool = make(chan types.FirecrackerVM, workers)
		for i := 0; i < workers; i++ {
			vm, err := firerunner.StartVM()
			if err != nil {
				log.Fatalf("Failed to start VM: %v", err)
			}
			vmPool <- *vm
		}
		fmt.Println("VM pool initialized")
		http.HandleFunc("/execute", getFirecrackerHandler(vmPool))
	} else {
		http.HandleFunc("/execute", handleRequestWithV8)
	}
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Stopping the server...")
		if option == "docker" && containerPool != nil {
			for i := 0; i < workers; i++ {
				container := <-containerPool
				docrunner.KillContainerAndGetLogs(&container)
			}
		} else if option == "firecracker" && vmPool != nil {
			for i := 0; i < workers; i++ {
				vm := <-vmPool
				vm.StopVMandCleanUp(vm.Machine, vm.VmmID)
			}
		}
		w.WriteHeader(http.StatusOK)
		os.Exit(0)
	})
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
