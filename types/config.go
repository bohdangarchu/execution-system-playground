package types

type WorkerConfig struct {
	Workers int `json:"workers"`
}

type FirecrackerConfig struct {
	CPUCount        int    `json:"cpuCount"`
	MemSizeMib      int    `json:"memSizeMib"`
	BinPath         string `json:"binPath"`
	KernelImagePath string `json:"kernelImagePath"`
	RootDrivePath   string `json:"rootDrivePath"`
}

type DockerConfig struct {
	ImageName     string `json:"imageName"`
	MaxMemSize    int    `json:"maxMemSize"`
	NanoCPUs      int    `json:"nanoCPUs"`
	ContainerPort int    `json:"containerPort"`
}

type ProcessIsolationConfig struct {
	WorkerPath   string `json:"workerPath"`
	CgroupName   string `json:"cgroupName"`
	CgroupMaxMem int    `json:"cgroupMaxMem"`
	CgroupMaxCPU int    `json:"cgroupMaxCPU"`
}

type Config struct {
	Isolation        string                  `json:"isolation"`
	Workers          int                     `json:"workers"`
	Firecracker      *FirecrackerConfig      `json:"firecracker"`
	Docker           *DockerConfig           `json:"docker"`
	ProcessIsolation *ProcessIsolationConfig `json:"processIsolation"`
}
