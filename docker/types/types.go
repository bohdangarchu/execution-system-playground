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
	Status       string          `json:"status"`
}

type ExecutionResult struct {
	Results []TestResult `json:"results"`
}

type ExecutionOutput struct {
	Value Argument `json:"value"`
	Error string   `json:"error"`
}
