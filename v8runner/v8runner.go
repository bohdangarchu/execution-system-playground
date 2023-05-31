package v8runner

import (
	"bytes"
	"fmt"
	"strings"

	"app/types"

	v8console "go.kuoruan.net/v8go-polyfills/console"
	v8 "rogchap.com/v8go"
)

func ExecuteJavaScript(code string) (string, error) {
	// for every console.log the logMessages variable gets updated
	// and is returned in the end
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

func ExecuteJsWithConsoleOutput(code string) (string, error) {
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

func RunFunctionWithInputs(submission types.FunctionSubmission, inOutArray []types.InputOutput) {
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
