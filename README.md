Config example 

```
{
    "isolation": "docker",
    "workers": 0,
    "firecracker": {
        "cpuCount": 1,
        "memSizeMib": 128
    },
    "docker": {
        "maxMemSize": 10000000,
        "nanoCPUs": 1000000000
    },
    "processIsolation": {
        "cgroupMaxMem": 100000000,
        "cgroupMaxCPU": 100
    }
}
```