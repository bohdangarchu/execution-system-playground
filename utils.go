package main

import (
	"app/types"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var defaultConfig = types.Config{
	Isolation: "docker",
	Workers:   1,
	Firecracker: &types.FirecrackerConfig{
		CPUCount:        1,
		MemSizeMib:      128,
		BinPath:         "/home/bohdan/software/firecracker/build/cargo_target/x86_64-unknown-linux-musl/debug/firecracker",
		KernelImagePath: "/home/bohdan/workspace/assets/hello-vmlinux.bin",
		RootDrivePath:   "/home/bohdan/workspace/uni/thesis/worker/firecracker/rootfs.ext4",
	},
	Docker: &types.DockerConfig{
		ImageName:     "execution-server",
		MaxMemSize:    10000000,
		NanoCPUs:      1000000000,
		ContainerPort: 8080,
	},
	ProcessIsolation: &types.ProcessIsolationConfig{
		WorkerPath:   "../worker/main",
		CgroupName:   "worker.slice",
		CgroupMaxMem: 100000000,
		CgroupMaxCPU: 100,
	},
}

func LoadConfig(path string) types.Config {
	// Open and read the JSON file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode JSON data into the Config struct
	decoder := json.NewDecoder(file)
	var config types.Config
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode config file: %v", err))
	}
	if config.ProcessIsolation == nil {
		config.ProcessIsolation = defaultConfig.ProcessIsolation
	}
	if config.Docker == nil {
		config.Docker = defaultConfig.Docker
	}
	if config.Firecracker == nil {
		config.Firecracker = defaultConfig.Firecracker
	}
	if config.Workers < 0 {
		panic("Number of workers cannot be negative")
	}
	if !contains(allowedImplValues, config.Isolation) {
		panic(fmt.Sprintf("Isolation must be one of: %s", strings.Join(allowedImplValues, ", ")))
	}
	if config.Docker.MaxMemSize < 0 {
		panic("Docker max memory size cannot be negative")
	}
	if config.Docker.NanoCPUs < 0 {
		panic("Docker nano CPUs cannot be negative")
	}
	if config.Docker.ContainerPort < 0 {
		panic("Docker container port cannot be negative")
	}
	if config.Firecracker.CPUCount < 0 {
		panic("Firecracker CPU count cannot be negative")
	}
	if config.Firecracker.MemSizeMib < 0 {
		panic("Firecracker memory size cannot be negative")
	}
	return config
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

var jsonSubmission = `
{
	"functionName": "addTwoNumbers",
	"code": "function addTwoNumbers(a, b) {\n  return a + b;\n}",
	"testCases": [
	  {
		"input": [
		  {
			"value": 3,
			"type": "number"
		  },
		  {
			"value": -10,
			"type": "number"
		  }
		]
	  }
	]
  }
`
