package main

import (
	"app/api"
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
	api.Run("docker")

	// err := performance.TimeDockerStartupAndSubmission()
	// if err != nil {
	// 	panic(err)
	// }

	// code := "function square(a) { return a*a; } console.log(square(55));"
	// output, err := v8runner.ExecuteJsWithConsoleOutput(code)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("out: ", output)

	// testPerformance()
	// performance.TestPerformance()
}
