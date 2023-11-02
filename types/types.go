package types

import (
	"context"
	"net"
	"os/exec"

	"github.com/docker/docker/client"

	"github.com/firecracker-microvm/firecracker-go-sdk"
)

type FunctionSubmission struct {
	FunctionName string     `json:"functionName"`
	Language     string     `json:"language"`
	Code         string     `json:"code"`
	TestCases    []TestCase `json:"testCases"`
}

type TestCase struct {
	Id         string   `json:"id"`
	InputArray []string `json:"input"`
}

type TestResult struct {
	TestCase TestCase        `json:"testCase"`
	Result   ExecutionOutput `json:"result"`
}

type ExecutionOutput struct {
	Output string `json:"output"`
	Error  string `json:"error"`
	Logs   string `json:"logs"`
}

type Response struct {
	Results []TestResult `json:"results"`
	Error   string       `json:"error"`
}

type FirecrackerVM struct {
	VmmCtx           context.Context
	VmmID            string
	Machine          *firecracker.Machine
	Ip               net.IP
	StopVMandCleanUp func() error
}

type DockerContainer struct {
	ContainerId string
	Port        string
	Cli         *client.Client
	Ctx         context.Context
}

type V8Worker struct {
	Id             string
	SocketPath     string
	ExecutablePath string
	Pid            int
	Cmd            *exec.Cmd
	CleanUp        func() error
}
