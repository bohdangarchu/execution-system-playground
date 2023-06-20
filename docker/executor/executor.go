package executor

import (
	"docker/types"
	"fmt"
	"strings"

	v8 "rogchap.com/v8go"
)

func RunFunctionWithInputs(submission types.FunctionSubmission) []types.ExecutionOutput {

	fmt.Println("function to be tested: \n", submission.Code)
	// creates a new V8 context with a new Isolate aka VM
	ctx := v8.NewContext()
	// executes a script on the global context
	_, err := ctx.RunScript(submission.Code, "main.js")
	if err != nil {
		panic(err)
	}

	results := make([]types.ExecutionOutput, len(submission.TestCases))
	for i, testCase := range submission.TestCases {
		params := strings.Join(testCase.InputArray, ", ")
		val, err := ctx.RunScript(
			fmt.Sprintf("%s(%s);", submission.FunctionName, params),
			"main.js",
		)
		if err != nil {
			// If an error occurs, create an ExecutionOutput object with the error message
			results[i] = types.ExecutionOutput{
				Value: "",
				Error: err.Error(),
			}
			continue
		}

		// If no error occurs, create an ExecutionOutput object with the actual output value
		results[i] = types.ExecutionOutput{
			Value: fmt.Sprintf("%v", val),
			Error: "",
		}

	}
	return results
}
