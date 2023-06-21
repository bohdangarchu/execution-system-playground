package types

type InputOutput struct {
	Input          []string
	ExpectedOutput string
}

type FunctionSubmissionOld struct {
	FunctionName    string
	ParameterString string
	CodeSubmission  string
}

// use the types below for the rest api interface

type FunctionSubmission struct {
	FunctionName string     `json:"functionName"`
	Code         string     `json:"code"`
	TestCases    []TestCase `json:"testCases"`
}

type TestCase struct {
	InputArray     []string `json:"input"`
	ExpectedOutput string   `json:"expectedOutput"`
}

type Value struct {
	Value string `json:"value"`
	Type  string `json:"type"`
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
	Value string `json:"value"`
	Error string `json:"error"`
}
