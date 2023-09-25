package types

type FirecrackerConfig struct {
	CPUCount   int `json:"cpuCount"`
	MemSizeMib int `json:"memSizeMib"`
}

type DockerConfig struct {
	MaxMemSize int `json:"maxMemSize"`
	NanoCPUs   int `json:"nanoCPUs"`
}

type ProcessIsolationConfig struct {
	CgroupMaxMem int `json:"cgroupMaxMem"`
	CgroupMaxCPU int `json:"cgroupMaxCPU"`
}

type Config struct {
	Isolation        string                  `json:"isolation"`
	Workers          int                     `json:"workers"`
	Firecracker      *FirecrackerConfig      `json:"firecracker"`
	Docker           *DockerConfig           `json:"docker"`
	ProcessIsolation *ProcessIsolationConfig `json:"processIsolation"`
}
