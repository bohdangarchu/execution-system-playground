curl --unix-socket /tmp/worker.sock 'http://localhost/execute' \
--header 'Content-Type: application/json' \
--data '{
  "functionName": "addTwoNumbers",
  "code": "function addTwoNumbers(a, b) {\n  console.log(\"args: \" + a + \", \" + b)\n  return a + b;\n}",
  "testCases": [ 
    {
      "id": "1",
      "input": ["1", "2"]
    }
  ]
}'