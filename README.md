Config example 

```
{
    "isolation": "firecracker",
    "workers": 1,
    "firecracker": {
        "memSizeMib": 256,
        "cpuQuota": 2000000,
        "cpuPeriod": 1000000
    },
    "docker": {
        "maxMemSize": 268435000,
        "cpuQuota": 1000000,
        "cpuPeriod": 1000000
    },
    "processIsolation": {
        "maxMemSize": 268435000,
        "cpuQuota": 1000000,
        "cpuPeriod": 1000000
    }
}
```

Submission example

```
{
  "functionName": "square",
  "code": "function square(a) {\n console.log(\"squaring\", a)\n return Math.pow(a, 2) \n}",
  "language": "js",
  "testCases": [
    {
      "id": "1",
      "input": ["5"]
    }
  ]
}
```