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
	path, err := CopyBaseRootfsWithIO(id)
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
		PathOnHost:   firecracker.String(FILESYSTEM_IMAGE_PATH),
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(true),
	}
}

func CopyBaseRootfs(id string) (string, error) {
	// copy rootfs.ext4 to /tmp/<id>-rootfs.ext4
	// Read the contents of the source file
	data, err := ioutil.ReadFile(FILESYSTEM_IMAGE_PATH)
	if err != nil {
		return "", err
	}

	tmpDir := os.TempDir()

	// Get the filename from the source path
	sourceFileName := filepath.Base(FILESYSTEM_IMAGE_PATH)

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
	tmpDir := os.TempDir()
	sourceFileName := filepath.Base(FILESYSTEM_IMAGE_PATH)
	destinationPath := filepath.Join(tmpDir, id+"-"+sourceFileName)
	err := copyFileWithIO(FILESYSTEM_IMAGE_PATH, destinationPath)
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

func getVMConfig(vmID string, config *types.FirecrackerConfig, useDefaultDrive bool) firecracker.Config {
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
			VcpuCount:   firecracker.Int64(1),
			CPUTemplate: models.CPUTemplate("C3"),
			MemSizeMib:  firecracker.Int64(int64(config.MemSizeMib)),
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
		if CheckVMHealth(vm) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
