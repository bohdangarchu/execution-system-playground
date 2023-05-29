package firerunner

const (
	socketPath = "/tmp/firecracker.sock"
	kernelPath = "/path/to/kernel/image"
	rootfsPath = "/path/to/rootfs/image"
)

// func main() {
// 	if err := startFirecracker(); err != nil {
// 		log.Fatal("Failed to start Firecracker:", err)
// 	}
// }

// func startFirecracker() error {
// 	// Create a new Firecracker microVM instance
// 	cfg := fc.Config{
// 		SocketPath: socketPath,
// 		LogLevel:   "Info",
// 	}

// 	if err := os.RemoveAll(socketPath); err != nil {
// 		return fmt.Errorf("failed to remove existing socket: %w", err)
// 	}

// 	vm, err := fc.NewMicroVM(context.Background(), cfg)
// 	if err != nil {
// 		return fmt.Errorf("failed to create microVM: %w", err)
// 	}

// 	// Configure the microVM
// 	if err := vm.Configure(
// 		fc.WithVCPUs(1),
// 		fc.WithMemSizeInMiB(256),
// 	); err != nil {
// 		return fmt.Errorf("failed to configure microVM: %w", err)
// 	}

// 	// Set up the network interface
// 	iface := fc.DefaultTapDevice()
// 	iface.HostDevName = "tap0"
// 	iface.GuestMac = net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x01}

// 	if err := vm.AddDevice(iface); err != nil {
// 		return fmt.Errorf("failed to add network device: %w", err)
// 	}

// 	// Set up the block device
// 	blockDev := fc.BlockDevice{
// 		HostPath:     rootfsPath,
// 		ReadOnly:     false,
// 		RootDevice:   true,
// 		IsRootDevice: true,
// 		IsReadOnly:   false,
// 		MonitorStdio: true,
// 		ShouldExist:  true,
// 		AutoDelete:   true,
// 		IsDirectIO:   false,
// 		RateLimiter:  fc.TokenBucketRateLimiter{Size: 10, OneTimeBurst: 10},
// 		Monitor:      fc.NewJSONFileMonitor(os.Stdout),
// 	}

// 	if err := vm.AddDevice(blockDev); err != nil {
// 		return fmt.Errorf("failed to add block device: %w", err)
// 	}

// 	// Load the kernel
// 	kernelImage, err := os.Open(kernelPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open kernel image: %w", err)
// 	}
// 	defer kernelImage.Close()

// 	kernel := fc.Kernel{
// 		Path: kernelPath,
// 	}
// 	if err := vm.LoadKernel(kernelImage, kernel); err != nil {
// 		return fmt.Errorf("failed to load kernel: %w", err)
// 	}

// 	// Start the microVM
// 	if err := vm.Start(context.Background()); err != nil {
// 		return fmt.Errorf("failed to start microVM: %w", err)
// 	}

// 	// Connect to the microVM's serial console
// 	cmd := exec.Command("sh", "-c", fmt.Sprintf("socat - UNIX-CONNECT:%s", socketPath))
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("failed to connect to serial console: %w", err)
// 	}

// 	return nil
// }
