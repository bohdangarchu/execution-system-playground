package main

import "app/performance"

func main() {
	// firerunner.RunFirecracker()

	performance.TimeDockerStartup()

	// code := "function square(a) { return a*a; } console.log(square(55));"
	// output, err := v8runner.ExecuteJsWithConsoleOutput(code)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("out: ", output)

	// testPerformance()
	// performance.TestPerformance()
}
