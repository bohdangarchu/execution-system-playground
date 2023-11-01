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
		PathOnHost:   firecracker.String(ROOTFS_PATH),
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(true),
	}
}

func CopyBaseRootfs(id string) (string, error) {
	// copy rootfs.ext4 to /tmp/<id>-rootfs.ext4
	root_drive_path := ROOTFS_PATH
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
	root_drive_path := ROOTFS_PATH
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

	parent := filepath.Dir(dst)
	err = os.MkdirAll(parent, 0666)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func createHardlink(src, dst string) error {
	parent := filepath.Dir(dst)
	err := os.MkdirAll(parent, 0666)
	if err != nil {
		return err
	}
	return os.Link(src, dst)
}

func getVMConfig(vmID string, cpuCount int64, memSizeMib int64, useDefaultDrive bool) firecracker.Config {
	// options for logging:
	// LogPath:           "/tmp/fc.log",
	// LogLevel:          "Debug",
	var drive models.Drive
	if useDefaultDrive {
		drive = getDefaultDrive()
	} else {
		drive = GetUniqueDrive(vmID)
	}
	// uid := 1000
	// gid := 1000
	// numaNode := 0
	return firecracker.Config{
		SocketPath: GetJailerSocketPath(vmID),
		// KernelImagePath:   KERNEL_IMAGE_PATH,
		// KernelImagePath:   "/srv/jailer/firecracker/" + vmID + "/root/hello-vmlinux.bin",
		KernelImagePath:   "./hello-vmlinux.bin",
		KernelArgs:        "console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw",
		Drives:            []models.Drive{drive},
		NetworkInterfaces: getCNINetworkInterfaces(),
		MachineCfg: models.MachineConfiguration{
			VcpuCount:   &cpuCount,
			CPUTemplate: models.CPUTemplate("C3"),
			MemSizeMib:  &memSizeMib,
		},
		// JailerCfg: &firecracker.JailerConfig{
		// 	UID:            &uid,
		// 	GID:            &gid,
		// 	ID:             vmID,
		// 	ExecFile:       FIRECRACKER_BIN_PATH,
		// 	NumaNode:       &numaNode,
		// 	Daemonize:      false,
		// 	CgroupVersion:  "2",
		// 	ChrootStrategy: firecracker.NewNaiveChrootStrategy(KERNEL_IMAGE_PATH),
		// },
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
