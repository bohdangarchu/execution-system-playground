package main

import (
	"app/firerunner"
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

func main() {
	// api.Run(firecracker)

	firerunner.RunStandaloneVM()

	// err := performance.TimeDockerStartupAndSubmission()
	// if err != nil {
	// 	panic(err)
	// }

	firerunner.RunFirecrackerVM()

}
