package firerunner

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetSocketPath(vmID string) string {
	filename := fmt.Sprintf("firecracker-%v.sock", vmID)
	dir := os.TempDir()
	return filepath.Join(dir, filename)
}

func GetJailerSocketPath(vmID string) string {
	return "/srv/jailer/firecracker/" + vmID + "/root/run/firecracker.socket"
}

func RemoveSocket(vmID string) error {
	socketPath := GetSocketPath(vmID)
	return os.Remove(socketPath)
}
