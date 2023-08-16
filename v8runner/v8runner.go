package v8runner

import (
	"bytes"
	"fmt"
	"strings"

	"app/types"

	v8console "go.kuoruan.net/v8go-polyfills/console"
	v8 "rogchap.com/v8go"
)

func RunFunctionWithInputs(submission types.FunctionSubmission) []types.TestResult {
	fmt.Println("function to be executed: \n", submission.Code)
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
	results := make([]types.TestResult, len(submission.TestCases))
	for i, testCase := range submission.TestCases {
		// TODO add error handling
		values := make([]v8.Valuer, len(testCase.InputArray))
		for i, input := range testCase.InputArray {
			value, _ := v8.NewValue(iso, input.Value)
			values[i] = value
		}
		val, err := function.Call(ctx.Global(), values...)
		if err != nil {
			// If an error occurs, create an ExecutionOutput object with the error message
			results[i] = types.TestResult{
				TestCase: testCase,
				ActualOutput: types.ExecutionOutput{
					Output: types.Argument{},
					Error:  err.Error(),
				},
			}
			continue
		}
		// If no error occurs, create an ExecutionOutput object with the actual output value
		results[i] = types.TestResult{
			TestCase: testCase,
			ActualOutput: types.ExecutionOutput{
				Output: v8ValueToArgument(*val),
				Error:  "",
			},
		}
	}
	return results
}

func ExecuteJsWithConsoleOutput(code string) (string, error) {
	// returns logs
	ctx := v8.NewContext()

	var buf bytes.Buffer

	if err := v8console.InjectTo(ctx, v8console.WithOutput(&buf)); err != nil {
		return "", err
	}

	_, err := ctx.RunScript(code, "main.js")
	if err != nil {
		return "", err
	}

	logs := buf.String()

	return logs, nil

}

func RunFunctionWithInputsManual(submission types.FunctionSubmissionOld, inOutArray []types.InputOutput) {
	// calls the function manually

	functionCode := fmt.Sprintf(
		`function %s (%s) {
	%s
}`, submission.FunctionName, submission.ParameterString, submission.CodeSubmission)

	fmt.Println("function to be tested: \n", functionCode)
	// creates a new V8 context with a new Isolate aka VM
	ctx := v8.NewContext()
	// executes a script on the global context
	_, err := ctx.RunScript(functionCode, "math.js")
	if err != nil {
		panic(err)
	}

	for _, inOut := range inOutArray {
		params := strings.Join(inOut.Input, ", ")
		val, err := ctx.RunScript(
			fmt.Sprintf("%s(%s);", submission.FunctionName, params),
			"main.js",
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("function execution result: %s\n", val)
		fmt.Printf("expected output: %s\n", inOut.ExpectedOutput)
	}
}

func ExecuteJavaScript(code string) (string, error) {
	// for every console.log the logMessages variable gets updated
	// and is returned in the end
	// v8console ("go.kuoruan.net/v8go-polyfills/console") should be used instead
	iso := v8.NewIsolate()

	ctx := v8.NewContext(iso)

	ctx.RunScript("var logMessages = [];", "main.js")
	ctx.RunScript("console.log = function() { logMessages.push.apply(logMessages, arguments); };", "main.js")

	_, err := ctx.RunScript(code, "main.js")

	if err != nil {
		return "", err
	}

	output, err := ctx.RunScript("logMessages", "main.js")

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
