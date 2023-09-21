package firerunner

import (
	"io/ioutil"
	"os"
	"path/filepath"

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

func getUniqueDrive(id string) models.Drive {
	path, err := CopyBaseRootfs(id)
	if err != nil {
		panic(err)
	}
	return models.Drive{
		DriveID:      firecracker.String("1"),
		PathOnHost:   &path,
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(false),
	}
}

func CopyBaseRootfs(id string) (string, error) {
	// copy rootfs.ext4 to /tmp/<id>-rootfs.ext4
	root_drive_path := "/home/bohdan/workspace/uni/thesis/worker/firecracker/rootfs.ext4"
	// Read the contents of the source file
	data, err := ioutil.ReadFile(root_drive_path)
	if err != nil {
		return "", err
	}

	tmpDir := os.TempDir()

	// Get the filename from the source path
	sourceFileName := filepath.Base(root_drive_path)

	// Create the destination file path in the temporary directory
	destinationPath := filepath.Join(tmpDir, id+"-"+sourceFileName)

	// Write the data to the destination file
	err = ioutil.WriteFile(destinationPath, data, 0644)
	if err != nil {
		return "", err
	}
	return destinationPath, nil
}

func getVMConfig(vmID string) firecracker.Config {
	// options for logging:
	// LogPath:           "/tmp/fc.log",
	// LogLevel:          "Debug",
	socket_path := GetSocketPath(vmID)
	var cpu_count int64 = 1
	var mem_size_mib int64 = 100
	drive := getUniqueDrive(vmID)
	return firecracker.Config{
		SocketPath:        socket_path,
		KernelImagePath:   KERNEL_IMAGE_PATH,
		KernelArgs:        "console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw",
		Drives:            []models.Drive{drive},
		NetworkInterfaces: getCNINetworkInterfaces(),
		MachineCfg: models.MachineConfiguration{
			VcpuCount:   &cpu_count,
			CPUTemplate: models.CPUTemplate("C3"),
			MemSizeMib:  &mem_size_mib,
		},
	}
}
