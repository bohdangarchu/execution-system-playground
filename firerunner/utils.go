package firerunner

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

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
	// options for logging:
	// LogPath:           "/tmp/fc.log",
	// LogLevel:          "Debug",
	socket_path := GetSocketPath(vmID)
	var cpu_count int64 = 1
	var mem_size_mib int64 = 100
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
