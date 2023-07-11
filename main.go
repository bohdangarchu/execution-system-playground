package main

import (
	"app/firerunner"
)

func main() {
	firerunner.StartVM()

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
