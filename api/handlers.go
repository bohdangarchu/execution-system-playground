package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/workerrunner"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func getFirecrackerHandler(vmPool chan types.FirecrackerVM) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get json string from request body
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()
		// get a VM from the pool
		vm := <-vmPool
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
		result, err := workerrunner.SendJsonToUnixSocket(worker.SocketPath, jsonSubmission)
		if workerrunner.CheckWorkerHealth(&worker) {
			workerPool <- worker
		} else {
			fmt.Println("worker ", worker.Pid, " is not running, starting a new one")
			newWorker := workerrunner.StartProcessWorker(
				config,
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
		// startTime := time.Now()
		// starts a new container for each request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		container, err := docrunner.StartExecutionServerInDocker(
			// with 0 docker will pick an available port
			"0",
			config.Docker,
		)
		defer docrunner.CleanUp(container, false)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to start docker container: %v", err), http.StatusInternalServerError)
			return
		}
		docrunner.WaitUntilAvailable(container)
		// workerIsReady := time.Now()
		result, err := docrunner.SendJSONSubmissionToDocker(container.Port, jsonSubmission)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
			return
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
		// executed := time.Now()
		// done := time.Now()
		// fmt.Printf(
		// 	"took %s to start container, %s to run the submission and %s to clean up\n",
		// 	workerIsReady.Sub(startTime), executed.Sub(workerIsReady), done.Sub(executed),
		// )
	}
}

func getFirecrackerHandlerWithNewVM(config *types.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// starts a new VM for each request
		// startTime := time.Now()
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		vm, err := firerunner.StartVM(true, config.Firecracker, false)
		defer vm.StopVMandCleanUp()
		firerunner.WaitUntilAvailable(vm)
		// workerIsReady := time.Now()
		if err != nil {
			fmt.Printf("failed to start vm: %v\n", err)
		}
		result, err := firerunner.RunSubmissionInsideVM(vm, jsonSubmission)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
			return
		}
		responseJSON := []byte(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
		// executed := time.Now()
		// done := time.Now()
		// fmt.Printf(
		// 	"took %s to start VM, %s to run the submission and %s to clean up\n",
		// 	workerIsReady.Sub(startTime), executed.Sub(workerIsReady), done.Sub(executed),
		// )
	}
}

func getWorkerHandlerWithNewWorker(config *types.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// starts a new worker for each request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		worker := workerrunner.StartProcessWorker(
			config.ProcessIsolation,
		)
		defer worker.CleanUp()
		workerrunner.WaitUntilAvailable(worker)
		result, err := workerrunner.SendJsonToUnixSocket(worker.SocketPath, jsonSubmission)
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

func handleRequestWithDocker(w http.ResponseWriter, r *http.Request) {
	// not used
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		fmt.Println(fmt.Sprintf("failed to parse request body: %v", err))
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
	// not used
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
