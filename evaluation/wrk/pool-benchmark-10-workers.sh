#!/bin/bash

set -e
echo "benchmark 10 workers"
echo "firecracker benchmark"

echo "one client, js"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "one client, lua"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, js"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "ten clients, lua"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, js"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

echo "fifty clients, lua"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill

