package main

import (
	"app/docrunner"
	"fmt"
)

func main() {
	// firerunner.RunFirecracker()

	err := docrunner.RunExecutionServerInDocker()

	fmt.Println("err: ", err)

	// code := "function square(a) { return a*a; } console.log(square(55));"
	// output, err := v8runner.ExecuteJsWithConsoleOutput(code)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("out: ", output)

	// testPerformance()
	// performance.TestPerformance()
}
