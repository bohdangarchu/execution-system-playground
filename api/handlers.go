package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/workerrunner"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getFirecrackerHandler(vmPool chan types.FirecrackerVM) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get json string from request body
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()
		// get a VM from the pool
		vm := <-vmPool
		fmt.Printf("VM %s Running job: %s", vm.VmmID, jsonSubmission)
		result, err := firerunner.RunSubmissionInsideVM(&vm, jsonSubmission)
		// push the VM back to the pool
		vmPool <- vm
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err.Error()), http.StatusBadRequest)
			return
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func getDockerHandler(containerPool chan types.DockerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		container := <-containerPool
		fmt.Printf("Container %s running job: %s", container.ContainerId, jsonSubmission)
		result, err := docrunner.SendJSONSubmissionToDocker(container.Port, jsonSubmission)
		containerPool <- container
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
			return
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func getWorkerHandler(workerPool chan types.V8Worker, config *types.ProcessIsolationConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		worker := <-workerPool
		fmt.Printf("Worker %s running job: %s", worker.Id, jsonSubmission)
		result, err := workerrunner.SendJsonToUnixSocket(worker.SocketPath, jsonSubmission)
		if workerrunner.CheckWorkerHealth(&worker) {
			workerPool <- worker
		} else {
			println("worker ", worker.Pid, " is not running, starting a new one")
			newWorker := workerrunner.StartProcessWorker(
				config.CgroupMaxMem,
				config.CgroupMaxCPU,
			)
			workerPool <- *newWorker
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute submission: %v", err), http.StatusInternalServerError)
			return
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func getDockerHandlerWithNewContainer(config *types.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// starts a new container for each request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		container, err := docrunner.StartExecutionServerInDocker(
			// with 0 docker will pick an available port
			"0",
			int64(config.Docker.MaxMemSize),
			int64(config.Docker.NanoCPUs),
		)
		defer docrunner.KillContainerAndGetLogs(container, false)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to start docker container: %v", err), http.StatusInternalServerError)
			return
		}
		docrunner.WaitUntilAvailable(container)
		result, err := docrunner.SendJSONSubmissionToDocker(container.Port, jsonSubmission)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
			return
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func getFirecrackerHandlerWithNewVM(config *types.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// starts a new VM for each request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		vm, err := firerunner.StartVM(true, config.Firecracker)
		defer vm.StopVMandCleanUp()
		firerunner.WaitUntilAvailable(vm)
		if err != nil {
			log.Fatalf("Failed to start VM: %v", err)
		}
		fmt.Printf("VM %s running job: %s", vm.VmmID, jsonSubmission)
		result, err := firerunner.RunSubmissionInsideVM(vm, jsonSubmission)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
			return
		}
		vm.StopVMandCleanUp()
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func getWorkerHandlerWithNewWorker(config *types.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// starts a new worker for each request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		worker := workerrunner.StartProcessWorker(
			config.ProcessIsolation.CgroupMaxMem,
			config.ProcessIsolation.CgroupMaxCPU,
		)
		workerrunner.WaitUntilAvailable(worker)
		fmt.Printf("Worker %s running job: %s", worker.Id, jsonSubmission)
		result, err := workerrunner.SendJsonToUnixSocket(worker.SocketPath, jsonSubmission)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute submission: %v", err), http.StatusInternalServerError)
			return
		}
		worker.Cmd.Process.Signal(os.Interrupt)
		if _, err := os.Stat(worker.SocketPath); err == nil {
			os.Remove(worker.SocketPath)
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func handleRequestWithDocker(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", err))
		return
	}

	jsonSubmission, err := json.Marshal(functionSubmission)
	responseString, err := docrunner.StartContainerAndRunSubmission(string(jsonSubmission))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
		return
	}
	responseJSON := []byte(responseString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func handleRequestWithFirecracker(w http.ResponseWriter, r *http.Request) {
	// get json string from request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	jsonSubmission := buf.String()

	responseString := firerunner.StartVMandRunSubmission(jsonSubmission)
	responseJSON := []byte(responseString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
