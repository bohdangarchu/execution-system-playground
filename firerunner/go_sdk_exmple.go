package firerunner

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

const FIRECRACKER_BIN_PATH = "/home/bohdan/software/firecracker/build/cargo_target/x86_64-unknown-linux-musl/debug/firecracker"
const KERNEL_IMAGE_PATH = "/home/bohdan/workspace/assets/hello-vmlinux.bin"

func RunSubmissionInsideVM(jsonSubmission string) string {
	startTime := time.Now()
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
	executionTime := time.Since(startTime)
	log.Printf("VM started in: %s", executionTime)

	result, err := executeJSONSubmissionInVM(
		vm.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP.String(),
		jsonSubmission,
	)
	if err != nil {
		log.Printf("Failed to execute JSON submission in VM: %v", err)
	}

	// time.Sleep(2 * time.Second)
	vm.StopVMM()
	log.Printf("Start machine was happy")
	return result
}

func getCNINetworkInterfaces() []firecracker.NetworkInterface {
	return []firecracker.NetworkInterface{{
		// Use CNI to get dynamic IP
		CNIConfiguration: &firecracker.CNIConfiguration{
			NetworkName: "fcnet",
			IfName:      "veth0",
		},
	}}
}

func getStaticNetworkInterfaces() []firecracker.NetworkInterface {
	return []firecracker.NetworkInterface{
		{
			StaticConfiguration: &firecracker.StaticNetworkConfiguration{
				MacAddress:  "2e:d5:b6:27:e8:8a",
				HostDevName: "tap0",
			},
		},
	}
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
		return "", fmt.Errorf("bad response status code: %d", resp.StatusCode)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(responseBody), nil
}

func getDrives() []models.Drive {
	root_drive_path := "/home/bohdan/workspace/uni/thesis/codebench-reference-project/agent/rootfs.ext4"
	return []models.Drive{
		{
			DriveID:      firecracker.String("1"),
			PathOnHost:   &root_drive_path,
			IsRootDevice: firecracker.Bool(true),
			IsReadOnly:   firecracker.Bool(false),
		},
	}
}

func getVMConfig(vmID string) firecracker.Config {
	socket_path := GetSocketPath(vmID)
	var cpu_count int64 = 1
	var mem_size_mib int64 = 512
	drives := getDrives()
	return firecracker.Config{
		SocketPath:        socket_path,
		KernelImagePath:   KERNEL_IMAGE_PATH,
		KernelArgs:        "console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw",
		Drives:            drives,
		NetworkInterfaces: getCNINetworkInterfaces(),
		MachineCfg: models.MachineConfiguration{
			VcpuCount:   &cpu_count,
			CPUTemplate: models.CPUTemplate("C3"),
			MemSizeMib:  &mem_size_mib,
		},
	}
}
