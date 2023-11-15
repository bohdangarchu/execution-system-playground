package firerunner

import (
	"app/types"
	"app/utils"
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
const FILESYSTEM_IMAGE_PATH = "../worker/firecracker/rootfs.ext4"

func RunSubmissionInsideVM(vm *types.FirecrackerVM, jsonSubmission string) (string, error) {
	return executeJSONSubmissionInVM(
		vm.Ip.String(),
		jsonSubmission,
	)
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

func StartVM(useDefaultDrive bool, config *types.FirecrackerConfig, debug bool) (*types.FirecrackerVM, error) {
	debugLevel := log.ErrorLevel
	if debug {
		debugLevel = log.InfoLevel
	}
	logger := &log.Logger{
		Out:       os.Stdout,
		Formatter: new(log.TextFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     debugLevel,
	}
	vmID := xid.New().String()
	fcCfg := getVMConfig(vmID, config, useDefaultDrive)
	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(logger)),
	}
	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	builder := firecracker.
		VMCommandBuilder{}.
		WithBin(FIRECRACKER_BIN_PATH).
		WithSocketPath(fcCfg.SocketPath)
	if debug {
		builder = builder.
			WithStdin(os.Stdin).
			WithStdout(os.Stdout).
			WithStderr(os.Stderr)
	}
	cmd := builder.Build(ctx)
	machineOpts = append(machineOpts, firecracker.WithProcessRunner(cmd))
	vm, err := firecracker.NewMachine(vmmCtx, fcCfg, machineOpts...)

	if err != nil {
		panic("Failed creating machine" + err.Error())
	}
	if err := vm.Start(vmmCtx); err != nil {
		panic("Failed to start machine" + err.Error())
	}
	pid, err := vm.PID()
	if err != nil {
		fmt.Printf("Failed to get PID: %v", err)
	}
	fmt.Printf("VM PID: %d\n", pid)
	manager := utils.CreateCPUCgroup("firecracker-"+vmID+".slice", int64(config.CPUQuota), uint64(config.CPUPeriod))
	err = manager.AddProc(uint64(pid))
	if err != nil {
		fmt.Println("error adding process to the cgroup: ", err.Error())
	}

	stopVMandCleanUp := func() error {
		if debug {
			fmt.Printf("Stopping VM: %s\n", vmID)
		}
		err := vm.StopVMM()
		if err != nil {
			fmt.Printf("Failed to stop VM: %v\n", err)
		}
		err = utils.RemoveFileIfExists(fcCfg.SocketPath)
		if !useDefaultDrive {
			os.Remove(*vm.Cfg.Drives[0].PathOnHost)
		}
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
	vm, err := StartVM(true, &types.FirecrackerConfig{
		MemSizeMib: 128,
		CPUQuota:   200000,
		CPUPeriod:  1000000,
	}, true)
	executionTime := time.Since(startTime)
	if err != nil {
		panic(fmt.Sprintf("Failed to start VM: %v", err))
	}
	fmt.Printf("VM started in %s\n", executionTime)
	fmt.Printf("VM IP: %s\n", vm.Ip.String())

	time.Sleep(1 * time.Second)
	vm.StopVMandCleanUp()
	fmt.Println("Start machine was happy")
}
