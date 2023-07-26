package main

import "app/api"

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

func main() {
	api.Run(firecracker)

	// firerunner.RunStandaloneVM()

	// test StartVM and vm.StopVMandCleanUp(vm.Machine, vm.VmmID)
	// vm, err := firerunner.StartVM()
	// if err != nil {
	// 	log.Fatalf("Failed to start VM: %v", err)
	// }
	// res, err := firerunner.RunSubmissionInsideVM(vm, jsonSubmission)
	// if err != nil {
	// 	log.Fatalf("Failed to run submission inside VM: %v", err)
	// }
	// log.Println(res)
	// vm.StopVMandCleanUp(vm.Machine, vm.VmmID)

}
