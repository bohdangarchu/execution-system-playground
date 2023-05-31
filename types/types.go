package types

type InputOutput struct {
	Input          []string
	ExpectedOutput string
}

type FunctionSubmission struct {
	FunctionName    string
	ParameterString string
	CodeSubmission  string
}
