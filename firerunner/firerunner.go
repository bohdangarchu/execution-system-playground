package firerunner

import (
	"app/types"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

const FIRECRACKER_BIN_PATH = "/home/bohdan/software/firecracker/build/cargo_target/x86_64-unknown-linux-musl/debug/firecracker"
const KERNEL_IMAGE_PATH = "/home/bohdan/workspace/assets/hello-vmlinux.bin"

func RunSubmissionInsideVM(vm *types.FirecrackerVM, jsonSubmission string) (string, error) {
	return executeJSONSubmissionInVM(
		vm.Ip.String(),
		jsonSubmission,
	)
}

func StartVMandRunSubmission(jsonSubmission string) string {
	startTimeStamp := time.Now()
	logger := log.New()
	vmID := xid.New().String()
	fcCfg := getVMConfig(vmID)
	defer RemoveSocket(vmID)
	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(logger)),
	}
	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()
	cmd := firecracker.VMCommandBuilder{}.
		WithBin(FIRECRACKER_BIN_PATH).
		WithSocketPath(fcCfg.SocketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)
	machineOpts = append(machineOpts, firecracker.WithProcessRunner(cmd))
	vm, err := firecracker.NewMachine(vmmCtx, fcCfg, machineOpts...)

	if err != nil {
		log.Fatalf("Failed creating machine: %s", err)
	}
	if err := vm.Start(vmmCtx); err != nil {
		log.Fatalf("Failed to start machine: %v", err)
	}
	bootTime := time.Since(startTimeStamp)
	bootTimeStamp := time.Now()
	log.Printf("VM started at: %v", bootTimeStamp)
	log.Printf("VM started in: %s", &bootTime)

	// for some reason takes >1s
	result, err := executeJSONSubmissionInVM(
		vm.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP.String(),
		jsonSubmission,
	)
	if err != nil {
		log.Printf("Failed to execute JSON submission in VM: %v", err)
	}
	executionTime := time.Since(bootTimeStamp)
	log.Printf("Submission executed in: %s", executionTime)

	// time.Sleep(30 * time.Second)
	vm.StopVMM()
	log.Printf("Start machine was happy")
	return result
}

func executeJSONSubmissionInVM(ip string, jsonSubmission string) (string, error) {
	url := "http://" + ip + ":8080/execute"

	// Create a request body as a bytes.Buffer
	requestBody := bytes.NewBuffer([]byte(jsonSubmission))

	// Make the POST request
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// resp body to string
		respBody, _ := ioutil.ReadAll(resp.Body)
		return "",
			fmt.Errorf("bad response status code: %d with error: %s", resp.StatusCode, respBody)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(responseBody), nil
}

func StartVM() (*types.FirecrackerVM, error) {
	logger := log.New()
	vmID := xid.New().String()
	fcCfg := getVMConfig(vmID)
	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(logger)),
	}
	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	cmd := firecracker.VMCommandBuilder{}.
		WithBin(FIRECRACKER_BIN_PATH).
		WithSocketPath(fcCfg.SocketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)
	machineOpts = append(machineOpts, firecracker.WithProcessRunner(cmd))
	vm, err := firecracker.NewMachine(vmmCtx, fcCfg, machineOpts...)

	if err != nil {
		log.Fatalf("Failed creating machine: %s", err)
	}
	if err := vm.Start(vmmCtx); err != nil {
		log.Fatalf("Failed to start machine: %v", err)
	}

	// todo dont need to pass the parameters
	stopVMandCleanUp := func(vm *firecracker.Machine, vmID string) error {
		fmt.Println("stoppping VM " + vmID)
		vm.StopVMM()
		RemoveSocket(vmID)
		os.Remove(*vm.Cfg.Drives[0].PathOnHost)
		vmmCancel()
		return nil
	}
	return &types.FirecrackerVM{
		VmmCtx:           vmmCtx,
		VmmID:            vmID,
		Machine:          vm,
		Ip:               vm.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP,
		StopVMandCleanUp: stopVMandCleanUp,
	}, nil
}

func RunStandaloneVM() {
	startTime := time.Now()
	vm, err := StartVM()
	executionTime := time.Since(startTime)
	if err != nil {
		log.Fatalf("Failed to start VM: %v", err)
	}
	log.Printf("VM started in: %s", executionTime)
	log.Printf("ip address: %s", vm.Machine.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP.String())

	time.Sleep(1 * time.Second)
	vm.StopVMandCleanUp(vm.Machine, vm.VmmID)
	log.Printf("Start machine was happy")
}
