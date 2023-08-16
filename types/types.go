package types

import (
	"context"
	"net"

	"github.com/docker/docker/client"

	"github.com/firecracker-microvm/firecracker-go-sdk"
)

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
	Id         string   `json:"id"`
	InputArray []string `json:"input"`
}

type TestResult struct {
	TestCase     TestCase        `json:"testCase"`
	ActualOutput ExecutionOutput `json:"actualOutput"`
}

type Argument struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type ExecutionOutput struct {
	Output string `json:"output"`
	Error  string `json:"error"`
	Logs   string `json:"logs"`
}

type Job struct {
	Submission string
	JobId      string
}

type JobResult struct {
	JobId  string
	Result string
	Err    error
}

type FirecrackerVM struct {
	VmmCtx           context.Context
	VmmID            string
	Machine          *firecracker.Machine
	Ip               net.IP
	StopVMandCleanUp func(vm *firecracker.Machine, vmID string) error
}

type DockerContainer struct {
	ContainerId string
	Port        string
	Cli         *client.Client
	Ctx         context.Context
}
