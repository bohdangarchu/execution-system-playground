package v8runner

import (
	"bytes"
	"fmt"
	"time"

	"app/types"

	v8console "go.kuoruan.net/v8go-polyfills/console"
	v8 "rogchap.com/v8go"
)

func RunSubmission(submission types.FunctionSubmission) ([]types.TestResult, error) {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	return RunSubmissionOnIsolate(iso, submission)
}

func RunSubmissionOnIsolate(iso *v8.Isolate, submission types.FunctionSubmission) ([]types.TestResult, error) {
	ctx := v8.NewContext(iso)
	_, err := ctx.RunScript(submission.Code, "main.js")
	if err != nil {
		return []types.TestResult{}, err
	}
	fnVal, err := ctx.Global().Get(submission.FunctionName)
	if err != nil {
		return []types.TestResult{}, err
	}
	function, err := fnVal.AsFunction()
	if err != nil {
		return []types.TestResult{}, err
	}
	results := make([]types.TestResult, len(submission.TestCases))
	for i, testCase := range submission.TestCases {
		results[i] = runTestCase(ctx, function, testCase)
	}
	return results, nil
}

func jsonToV8Values(ctx *v8.Context, arguments []string) ([]v8.Valuer, error) {
	values := make([]v8.Valuer, len(arguments))
	for i, input := range arguments {
		value, err := v8.JSONParse(ctx, input)
		if err != nil {
			return values, err
		}
		values[i] = value
	}
	return values, nil
}

func runTestCase(ctx *v8.Context, fun *v8.Function, testCase types.TestCase) types.TestResult {
	values, err := jsonToV8Values(ctx, testCase.InputArray)
	if err != nil {
		return types.TestResult{
			TestCase: testCase,
			ActualOutput: types.ExecutionOutput{
				Output: "",
				Error:  err.Error(),
				Logs:   "",
			},
		}
	}
	var buf bytes.Buffer
	if err := v8console.InjectTo(ctx, v8console.WithOutput(&buf)); err != nil {
		return types.TestResult{
			TestCase: testCase,
			ActualOutput: types.ExecutionOutput{
				Output: "",
				Error:  err.Error(),
				Logs:   "",
			},
		}
	}
	val, err := callFunctionWithTimeout(ctx, fun, values, 1000*time.Millisecond)
	logs := buf.String()
	if err != nil {
		return types.TestResult{
			TestCase: testCase,
			ActualOutput: types.ExecutionOutput{
				Output: "",
				Error:  err.Error(),
				Logs:   logs,
			},
		}
	}
	jsonValue, err := v8.JSONStringify(ctx, val)
	if err != nil {
		return types.TestResult{
			TestCase: testCase,
			ActualOutput: types.ExecutionOutput{
				Output: jsonValue,
				Error:  err.Error(),
				Logs:   logs,
			},
		}
	}
	return types.TestResult{
		TestCase: testCase,
		ActualOutput: types.ExecutionOutput{
			Output: jsonValue,
			Error:  "",
			Logs:   logs,
		},
	}
}

func callFunctionWithTimeout(ctx *v8.Context, fun *v8.Function, values []v8.Valuer, timeout time.Duration) (*v8.Value, error) {
	vals := make(chan *v8.Value, 1)
	errs := make(chan error, 1)
	memoryErr := make(chan error, 1)
	go func() {
		val, err := fun.Call(ctx.Global(), values...)
		if err != nil {
			errs <- err
			return
		}
		vals <- val
	}()
	// 20 MB memory limit
	go monitorMemoryUsage(ctx, memoryErr, 20000000)
	select {
	case val := <-vals:
		return val, nil
	case err := <-errs:
		return nil, err
	case <-memoryErr:
		ctx.Isolate().TerminateExecution()
		return nil, fmt.Errorf("memory usage exceeded")
	case <-time.After(timeout):
		ctx.Isolate().TerminateExecution()
		// will get a termination error back from the running script
		<-errs
		return nil, fmt.Errorf("execution timed out after %d ms", timeout.Milliseconds())
	}
}

func ExecuteJsWithConsoleOutput(code string) (string, error) {
	// not used
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
