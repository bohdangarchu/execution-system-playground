package types

type FirecrackerConfig struct {
	MemSizeMib int `json:"memSizeMib"`
	CPUQuota   int `json:"cpuQuota"`
	CPUPeriod  int `json:"cpuPeriod"`
}

type DockerConfig struct {
	MaxMemSize int `json:"maxMemSize"`
	CPUQuota   int `json:"cpuQuota"`
	CPUPeriod  int `json:"cpuPeriod"`
}

type ProcessIsolationConfig struct {
	MaxMemSize int `json:"maxMemSize"`
	CPUQuota   int `json:"cpuQuota"`
	CPUPeriod  int `json:"cpuPeriod"`
}

type Config struct {
	Isolation        string                  `json:"isolation"`
	Workers          int                     `json:"workers"`
	Firecracker      *FirecrackerConfig      `json:"firecracker"`
	Docker           *DockerConfig           `json:"docker"`
	ProcessIsolation *ProcessIsolationConfig `json:"processIsolation"`
}
