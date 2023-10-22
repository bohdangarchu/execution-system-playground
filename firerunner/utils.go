package firerunner

import (
	"app/types"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

func GetUniqueDrive(id string) models.Drive {
	startTime := time.Now()
	path, err := CopyBaseRootfsWithIO(id)
	endTime := time.Now()
	fmt.Printf("Copying rootfs took: %s\n", endTime.Sub(startTime))
	if err != nil {
		panic(fmt.Sprintf("Failed to copy rootfs: %v", err))
	}
	return models.Drive{
		DriveID:      firecracker.String("1"),
		PathOnHost:   &path,
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(false),
	}
}

func getDefaultDrive() models.Drive {
	return models.Drive{
		DriveID:      firecracker.String("1"),
		PathOnHost:   firecracker.String("/home/bohdan/workspace/uni/thesis/worker/firecracker/rootfs.ext4"),
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(true),
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

func CopyBaseRootfsWithIO(id string) (string, error) {
	// copy rootfs.ext4 to /tmp/<id>-rootfs.ext4
	root_drive_path := "/home/bohdan/workspace/uni/thesis/worker/firecracker/rootfs.ext4"
	tmpDir := os.TempDir()
	sourceFileName := filepath.Base(root_drive_path)
	destinationPath := filepath.Join(tmpDir, id+"-"+sourceFileName)
	err := copyFileWithIO(root_drive_path, destinationPath)
	if err != nil {
		return "", err
	}
	return destinationPath, nil
}

func copyFileWithIO(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func getVMConfig(vmID string, cpuCount int64, memSizeMib int64, useDefaultDrive bool) firecracker.Config {
	// options for logging:
	// LogPath:           "/tmp/fc.log",
	// LogLevel:          "Debug",
	socket_path := GetSocketPath(vmID)
	var drive models.Drive
	if useDefaultDrive {
		drive = getDefaultDrive()
	} else {
		drive = GetUniqueDrive(vmID)
	}
	return firecracker.Config{
		SocketPath:        socket_path,
		KernelImagePath:   KERNEL_IMAGE_PATH,
		KernelArgs:        "console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw",
		Drives:            []models.Drive{drive},
		NetworkInterfaces: getCNINetworkInterfaces(),
		MachineCfg: models.MachineConfiguration{
			VcpuCount:   &cpuCount,
			CPUTemplate: models.CPUTemplate("C3"),
			MemSizeMib:  &memSizeMib,
		},
	}
}

func CheckVMHealth(vm *types.FirecrackerVM) bool {
	url := "http://" + vm.Ip.String() + ":8080/health"
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// check response status code
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func WaitUntilAvailable(vm *types.FirecrackerVM) {
	for {
		fmt.Println("Waiting for VM to be available...")
		if CheckVMHealth(vm) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
