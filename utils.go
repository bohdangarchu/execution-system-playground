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
		MemSizeMib: 128,
		CPUQuota:   125000,
		CPUPeriod:  1000000,
	},
	Docker: &types.DockerConfig{
		MaxMemSize: 10000000,
		CPUQuota:   125000,
		CPUPeriod:  1000000,
	},
	ProcessIsolation: &types.ProcessIsolationConfig{
		MaxMemSize: 100000000,
		CPUQuota:   125000,
		CPUPeriod:  1000000,
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
	if config.Docker.CPUQuota < 0 {
		panic("Docker CPU quota cannot be negative")
	}
	if config.Docker.CPUPeriod < 0 {
		panic("Docker CPU period cannot be negative")
	}
	if config.Firecracker.CPUQuota < 0 {
		panic("Firecracker CPU quota cannot be negative")
	}
	if config.Firecracker.CPUPeriod < 0 {
		panic("Firecracker CPU period cannot be negative")
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
