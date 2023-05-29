package main

import (
	"app/docrunner"
	"app/firerunner"
	"app/v8runner"
	"fmt"
	"time"
)

type inputOutput struct {
	input          []string
	expectedOutput string
}

type functionSubmission struct {
	name             string
	code             string
	inputOutputPairs []inputOutput
}

func main() {
	firerunner.RunFirecracker()
}

func testPerformance() {
	jsSubmissions := []string{
		`function square(a) { return a*a; } console.log(square(55));`,
		`function factorial(n) { if (n === 0) return 1; return n * factorial(n - 1); } console.log(factorial(10));`,
		`function fibonacci(n) { if (n <= 1) return n; return fibonacci(n - 1) + fibonacci(n - 2); } console.log(fibonacci(10));`,
	}

	fmt.Println("===== Performance Test Results =====")

	for i, jsSubmission := range jsSubmissions {
		fmt.Printf("\nTest Case %d:\n", i+1)
		fmt.Println("JavaScript Code:\n", jsSubmission)

		// Measure execution time for v8runner
		startTimeV8 := time.Now()
		_, outputV8, errV8 := v8runner.ExecuteJavaScript(jsSubmission)
		if errV8 != nil {
			fmt.Println("v8runner Error:", errV8)
		}
		elapsedTimeV8 := time.Since(startTimeV8)
		fmt.Println("v8runner Output:", outputV8)
		fmt.Println("v8runner Execution Time:", elapsedTimeV8)

		// Measure execution time for docrunner
		startTimeDoc := time.Now()
		outDoc, errDoc := docrunner.RunJsInDocker(jsSubmission)
		if errDoc != nil {
			fmt.Println("docrunner Error:", errDoc)
		}
		elapsedTimeDoc := time.Since(startTimeDoc)
		fmt.Println("docrunner Output:", outDoc)
		fmt.Println("docrunner Execution Time:", elapsedTimeDoc)
	}
}

func testSecurity() {
	jsSubmissions := []string{
		`function exploitFileAccess() {
			var fs = require('fs');
			var content = fs.readFileSync('/etc/passwd');
			console.log(content);
		}
		exploitFileAccess();`,
		// Add more test cases here...
	}

	fmt.Println("===== Security Test Results =====")

	for i, jsSubmission := range jsSubmissions {
		fmt.Printf("\nTest Case %d:\n", i+1)
		fmt.Println("JavaScript Code:\n", jsSubmission)

		// Security test using v8runner
		fmt.Println("=== v8runner ===")
		startTimeV8 := time.Now()
		_, outputV8, errV8 := v8runner.ExecuteJavaScript(jsSubmission)
		if errV8 != nil {
			fmt.Println("v8runner Error:", errV8)
		}
		elapsedTimeV8 := time.Since(startTimeV8)
		fmt.Println("v8runner Output:", outputV8)
		fmt.Println("v8runner Execution Time:", elapsedTimeV8)

		// Security test using docrunner
		fmt.Println("=== docrunner ===")
		startTimeDoc := time.Now()
		outDoc, errDoc := docrunner.RunJsInDocker(jsSubmission)
		if errDoc != nil {
			fmt.Println("docrunner Error:", errDoc)
		}
		elapsedTimeDoc := time.Since(startTimeDoc)
		fmt.Println("docrunner Output:", outDoc)
		fmt.Println("docrunner Execution Time:", elapsedTimeDoc)
	}
}
