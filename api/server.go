package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/workerrunner"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Run(option string, workers int) {
	if workers > 0 {
		fmt.Printf("Running with %d %s workers\n", workers, option)
		runInWorkerPool(option, workers)
	} else {
		fmt.Printf("Running with a new %s worker for each request\n", option)
		if option == "docker" {
			http.HandleFunc("/execute", handleRequestWithNewContainer)
		} else if option == "firecracker" {
			http.HandleFunc("/execute", handleRequestWithNewFirecrackerVM)
		} else {
			http.HandleFunc("/execute", handleRequestWithNewProcessWorker)
		}
		http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Stopping the server...")
			w.WriteHeader(http.StatusOK)
			os.Exit(0)
		})
	}
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func runInWorkerPool(option string, workers int) {
	var vmPool chan types.FirecrackerVM
	var containerPool chan types.DockerContainer
	var workerPool chan types.V8Worker
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
		startTime := time.Now()
		for i := 0; i < workers; i++ {
			vm, err := firerunner.StartVM(false)
			if err != nil {
				log.Fatalf("Failed to start VM: %v", err)
			}
			vmPool <- *vm
		}
		elapsed := time.Since(startTime)
		fmt.Printf("VM pool initialized in %s\n", elapsed)
		http.HandleFunc("/execute", getFirecrackerHandler(vmPool))
	} else {
		workerPool = make(chan types.V8Worker, workers)
		for i := 0; i < workers; i++ {
			worker := workerrunner.StartProcessWorker()
			workerPool <- *worker
		}
		http.HandleFunc("/execute", getWorkerHandler(workerPool))
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
				vm.StopVMandCleanUp()
			}
		} else {
			for i := 0; i < workers; i++ {
				worker := <-workerPool
				worker.Cmd.Process.Signal(os.Interrupt)
				// check if the socket file exists
				// if it does, remove it
				if _, err := os.Stat(worker.SocketPath); err == nil {
					os.Remove(worker.SocketPath)
				}
			}
		}
		w.WriteHeader(http.StatusOK)
		os.Exit(0)
	})
}
