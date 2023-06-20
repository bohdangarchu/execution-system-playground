package firerunner

import (
	"context"
	"os"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	log "github.com/sirupsen/logrus"
)

func RunFirecracker() {
	firecracker_bin_path := "/home/bohdan/software/firecracker/build/cargo_target/x86_64-unknown-linux-musl/debug/firecracker"
	rootDrivePath := "/home/bohdan/workspace/uni/thesis/firecracker/hello-rootfs.ext4"
	kernelImagePath := "/home/bohdan/workspace/uni/thesis/firecracker/hello-vmlinux.bin"
	socket_path := "/home/bohdan/workspace/uni/thesis/firecracker/ficracker.socket"
	var cpu_count int64 = 1
	var mem_size_mib int64 = 512
	logger := log.New()
	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()
	devices := []models.Drive{}
	rootDrive := models.Drive{
		DriveID:      firecracker.String("1"),
		PathOnHost:   &rootDrivePath,
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(false),
	}
	devices = append(devices, rootDrive)
	fcCfg := firecracker.Config{
		SocketPath:      socket_path,
		KernelImagePath: kernelImagePath,
		KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off",
		Drives:          devices,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:   &cpu_count,
			CPUTemplate: models.CPUTemplate("C3"),
			MemSizeMib:  &mem_size_mib,
		},
	}
	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(logger)),
	}
	cmd := firecracker.VMCommandBuilder{}.
		WithBin(firecracker_bin_path).
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
	defer vm.StopVMM()

	// wait for the VMM to exit
	if err := vm.Wait(vmmCtx); err != nil {
		log.Fatalf("Wait returned an error %s", err)
	}
	log.Printf("Start machine was happy")
}