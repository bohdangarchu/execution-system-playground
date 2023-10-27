wrk.method = "POST"
wrk.body   = [[
    {
        "functionName": "factorial",
        "code": "function factorial(n)\n  if n <= 1 then\n    return 1\n  else\n    return n * factorial(n - 1)\n  end\nend",
        "language": "lua",
        "testCases": [
          {
            "id": "1",
            "input": ["5"]
          }
        ]
    }
]]
wrk.headers["Content-Type"] = "application/json"