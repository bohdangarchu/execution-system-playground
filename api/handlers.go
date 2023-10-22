package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/v8runner"
	"app/workerrunner"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rs/xid"
)

func getFirecrackerHandler(vmPool chan types.FirecrackerVM) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get json string from request body
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()
		// get a VM from the pool
		vm := <-vmPool
		job := types.Job{
			Submission: jsonSubmission,
			JobId:      xid.New().String(),
		}
		result := types.JobResult{
			JobId: job.JobId,
		}
		fmt.Printf("VM %s Running job: %s", vm.VmmID, job.Submission)
		result.Result, result.Err = firerunner.RunSubmissionInsideVM(&vm, job.Submission)
		// push the VM back to the pool
		vmPool <- vm
		if result.Err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", result.Err.Error()), http.StatusBadRequest)
			return
		}
		responseJSON := []byte(result.Result)
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
		if workerrunner.IsWorkerRunning(&worker) {
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
		defer docrunner.KillContainerAndGetLogs(container)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to start docker container: %v", err), http.StatusInternalServerError)
			return
		}
		time.Sleep(50 * time.Millisecond)
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

func getFirecrackerHandlerWithNewVM(config *types.Config, drivePool chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// starts a new VM for each request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonSubmission := buf.String()

		// get a drive from the pool
		drivePath := <-drivePool
		drive := firerunner.GetDriveFromPath(drivePath)
		vm, err := firerunner.StartVMWithDrive(config.Firecracker, drive)
		defer vm.StopVMandCleanUp()
		if err != nil {
			log.Fatalf("Failed to start VM: %v", err)
		}
		fmt.Printf("VM %s running job: %s", vm.VmmID, jsonSubmission)
		result, err := firerunner.RunSubmissionInsideVM(vm, jsonSubmission)
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
		defer workerrunner.KillWorker(worker)
		fmt.Printf("Worker %s running job: %s", worker.Id, jsonSubmission)
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

func getV8Handler(isolatePool chan types.V8Isolate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var functionSubmission types.FunctionSubmission
		err := json.NewDecoder(r.Body).Decode(&functionSubmission)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
			log.Println(fmt.Sprintf("failed to parse request body: %v", err))
			return
		}
		iso := <-isolatePool
		// Execute the JavaScript code
		outputArray, err := v8runner.RunSubmissionOnIsolate(iso.Isolate, functionSubmission)
		// push the isolate back to the pool
		isolatePool <- iso
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute submission: %v", err), http.StatusInternalServerError)
			return
		}
		// Convert the result to JSON
		responseJSON, err := json.Marshal(outputArray)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to convert result to JSON: %v", err), http.StatusInternalServerError)
			return
		}
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
	_, _ = w.Write(responseJSON)
}

func handleRequestWithV8(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", err))
		return
	}
	// Execute the JavaScript code
	outputArray, err := v8runner.RunSubmission(functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute submission: %v", err), http.StatusInternalServerError)
		return
	}
	// Convert the result to JSON
	responseJSON, err := json.Marshal(outputArray)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to convert result to JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}
