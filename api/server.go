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

func Run(config *types.Config) {
	if config.Workers > 0 {
		fmt.Printf("Running with %d %s workers\n", config.Workers, config.Isolation)
		runInWorkerPool(config)
	} else {
		fmt.Printf("Running with a new %s worker for each request\n", config.Isolation)
		if config.Isolation == "docker" {
			http.HandleFunc("/execute", getDockerHandlerWithNewContainer(config))
		} else if config.Isolation == "firecracker" {
			http.HandleFunc("/execute", getFirecrackerHandlerWithNewVM(config))
		} else {
			http.HandleFunc("/execute", getWorkerHandlerWithNewWorker(config))
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

func runInWorkerPool(config *types.Config) {
	var vmPool chan types.FirecrackerVM
	var containerPool chan types.DockerContainer
	var workerPool chan types.V8Worker
	if config.Isolation == "docker" {
		containerPool = make(chan types.DockerContainer, config.Workers)
		for i := 0; i < config.Workers; i++ {
			container, err := docrunner.StartExecutionServerInDocker(
				// with 0 docker will pick an available port
				"0",
				int64(config.Docker.MaxMemSize),
				int64(config.Docker.NanoCPUs),
			)
			if err != nil {
				panic(err)
			}
			docrunner.WaitUntilAvailable(container)
			containerPool <- *container
		}
		go monitorContainerHealth(containerPool, config)
		http.HandleFunc("/execute", getDockerHandler(containerPool))
	} else if config.Isolation == "firecracker" {
		vmPool = make(chan types.FirecrackerVM, config.Workers)
		startTime := time.Now()
		for i := 0; i < config.Workers; i++ {
			// use a unique drive for every VM
			vm, err := firerunner.StartVM(true, config.Firecracker, true)
			firerunner.WaitUntilAvailable(vm)
			if err != nil {
				log.Fatalf("Failed to start VM: %v", err)
			}
			vmPool <- *vm
		}
		elapsed := time.Since(startTime)
		fmt.Printf("VM pool initialized in %s\n", elapsed)
		go monitorVMHealth(vmPool, config, true, true)
		http.HandleFunc("/execute", getFirecrackerHandler(vmPool))
	} else {
		workerPool = make(chan types.V8Worker, config.Workers)
		for i := 0; i < config.Workers; i++ {
			worker := workerrunner.StartProcessWorker(
				config.ProcessIsolation.CgroupMaxMem,
				config.ProcessIsolation.CgroupMaxCPU,
			)
			workerrunner.WaitUntilAvailable(worker)
			workerPool <- *worker
		}
		go monitorV8Worker(workerPool, config.ProcessIsolation)
		http.HandleFunc("/execute", getWorkerHandler(workerPool, config.ProcessIsolation))
	}
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Stopping the server...")
		if config.Isolation == "docker" && containerPool != nil {
			for i := 0; i < config.Workers; i++ {
				container := <-containerPool
				docrunner.CleanUp(&container, false)
			}
		} else if config.Isolation == "firecracker" && vmPool != nil {
			for i := 0; i < config.Workers; i++ {
				vm := <-vmPool
				vm.StopVMandCleanUp()
			}
		} else {
			for i := 0; i < config.Workers; i++ {
				worker := <-workerPool
				workerrunner.KillWorker(&worker)
			}
		}
		w.WriteHeader(http.StatusOK)
		os.Exit(0)
	})
}
