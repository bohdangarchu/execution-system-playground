while true; do
  curl --unix-socket /tmp/worker-cl3bvbfvrg9rh076smng.sock 'http://localhost/execute' \
  --header 'Content-Type: application/json' \
  --data '{
    "functionName": "factorial",
    "code": "function factorial(n) {\n  if (n <= 1) {\n    return 1;\n  } else {\n    return n * factorial(n - 1);\n  }\n}",
    "language": "js",
    "testCases": [
      {
        "id": "1",
        "input": ["5"]
      }
    ]
  }'
done
