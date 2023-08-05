package main

import (
	"app/api"
	"app/firerunner"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
)

var jsonSubmission = `
{
	"functionName": "addTwoNumbers",
	"code": "function addTwoNumbers(a, b) {\n  return a + b;\n}",
	"testCases": [
	  {
		"input": [
		  {
			"value": 3,
			"type": "number"
		  },
		  {
			"value": -10,
			"type": "number"
		  }
		]
	  }
	]
  }
`

const (
	firecracker = iota
	docker
	v8
)

type implValue string

var allowedImplValues = []string{"docker", "firecracker", "v8"}

func (iv *implValue) String() string {
	return string(*iv)
}

func (iv *implValue) Set(value string) error {
	for _, allowedValue := range allowedImplValues {
		if value == allowedValue {
			*iv = implValue(value)
			return nil
		}
	}
	return errors.New("invalid value for --impl flag")
}

func main() {
	var impl implValue
	var workers int
	// Define flags
	flag.Var(&impl, "impl", fmt.Sprintf("Choose from: %s", strings.Join(allowedImplValues, ", ")))
	flag.IntVar(&workers, "workers", 1, "Number of workers (int)")
	// Parse the command line flags
	flag.Parse()
	// Print the values
	fmt.Println("impl:", impl)
	fmt.Println("workers:", workers)
	api.Run(impl.String(), workers)
}

func runVM() {
	vm, err := firerunner.StartVM()
	defer vm.StopVMandCleanUp(vm.Machine, vm.VmmID)
	if err != nil {
		log.Fatalf("Failed to start VM: %v", err)
	}
	res, err := firerunner.RunSubmissionInsideVM(vm, jsonSubmission)
	if err != nil {
		log.Fatalf("Failed to run submission inside VM: %v", err)
	}
	log.Println(res)
}
