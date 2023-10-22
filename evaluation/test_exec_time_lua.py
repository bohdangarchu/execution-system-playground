import testing
import sys


def get_submission():
    return r"""
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
"""

if __name__ == "__main__":
    concurrent_requests = 5
    if len(sys.argv) > 1:
        concurrent_requests = int(sys.argv[1])

    testing.run(concurrent_requests, get_submission(), debug=False)