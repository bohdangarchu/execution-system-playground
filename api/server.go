package api

import (
	"app/docrunner"
	"app/firerunner"
	"app/types"
	"app/v8runner"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/xid"
)

const (
	firecracker = iota
	docker
	v8
)

func Run(option int) {
	workers := 2
	vmPool := make(chan types.FirecrackerVM, workers)
	if option == docker {
		http.HandleFunc("/", handleRequestWithDocker)
	} else if option == firecracker {

		for i := 0; i < workers; i++ {
			vm, err := firerunner.StartVM()
			if err != nil {
				log.Fatalf("Failed to start VM: %v", err)
			}
			// doesnt work
			// defer vm.StopVMandCleanUp(vm.Machine, vm.VmmID)
			vmPool <- *vm
		}
		fmt.Println("VM pool initialized")
		http.HandleFunc("/", getFirecrackerHandler(vmPool))
	} else {
		http.HandleFunc("/", handleRequestWithV8)
	}
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Killing the server...")
		// stop all the VMs
		if option == firecracker {
			for i := 0; i < workers; i++ {
				vm := <-vmPool
				vm.StopVMandCleanUp(vm.Machine, vm.VmmID)
			}
		}
		w.WriteHeader(http.StatusOK)
		os.Exit(0)
	})
	log.Println("Listening on :8081...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}

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
		if result.Err != nil {
			http.Error(w, fmt.Sprintf("failed to execute the submission: %v", result.Err.Error()), http.StatusBadRequest)
			return
		}
		// push the VM back to the pool
		vmPool <- vm
		responseJSON := []byte(result.Result)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseJSON)
	}
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

func handleRequestWithDocker(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
		return
	}

	jsonSubmission, err := json.Marshal(functionSubmission)
	responseString, err := docrunner.RunSubmissionInsideDocker(string(jsonSubmission))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute the submission: %v", err), http.StatusBadRequest)
		return
	}
	responseJSON := []byte(responseString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func handleRequestWithV8(w http.ResponseWriter, r *http.Request) {
	var functionSubmission types.FunctionSubmission
	err := json.NewDecoder(r.Body).Decode(&functionSubmission)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		log.Println(fmt.Sprintf("failed to parse request body: %v", r.Body))
		return
	}

	outputArray := v8runner.RunFunctionWithInputs(functionSubmission)

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
