package types

type FunctionSubmission struct {
	FunctionName string     `json:"functionName"`
	Code         string     `json:"code"`
	TestCases    []TestCase `json:"testCases"`
}

type TestCase struct {
	InputArray []Argument `json:"input"`
}

type Argument struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type TestResult struct {
	TestCase     TestCase        `json:"testCase"`
	ActualOutput ExecutionOutput `json:"actualOutput"`
}

type ExecutionOutput struct {
	Output Argument `json:"output"`
	Error  string   `json:"error"`
}
