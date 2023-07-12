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

func main() {
	// result := firerunner.RunSubmissionInsideVM(jsonSubmission)
	// fmt.Println("result: ", result)

	// err := performance.TimeDockerStartupAndSubmission()
	// if err != nil {
	// 	panic(err)
	// }

	firerunner.RunFirecrackerVM()

}
