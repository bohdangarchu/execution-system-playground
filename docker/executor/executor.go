package executor

import (
	"docker/types"
	"fmt"

	v8 "rogchap.com/v8go"
)

func RunFunctionWithInputs(submission types.FunctionSubmission) []types.ExecutionOutput {

	fmt.Println("function to be tested: \n", submission.Code)
	// creates a new V8 context with a new Isolate aka VM
	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)
	// executes a script on the global context
	_, err := ctx.RunScript(submission.Code, "main.js")
	if err != nil {
		// TODO add error handling
		panic(err)
	}

	fnVal, err := ctx.Global().Get(submission.FunctionName)

	if err != nil {
		// TODO add error handling
		panic(err)
	}

	function, err := fnVal.AsFunction()

	if err != nil {
		// TODO add error handling
		panic(err)
	}

	results := make([]types.ExecutionOutput, len(submission.TestCases))
	for i, testCase := range submission.TestCases {
		// TODO add error handling
		values := make([]v8.Valuer, len(testCase.InputArray))
		for i, input := range testCase.InputArray {
			value, _ := v8.NewValue(iso, input)
			values[i] = value
		}

		val, err := function.Call(ctx.Global(), values...)
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
			Value: val.String(),
			Error: "",
		}

	}
	return results
}
