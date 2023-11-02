#!/bin/bash

set -e

n=5
echo "benchmark baseline, worker per request"

echo "one client, js"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

sleep "$n"

echo "one client, lua"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

sleep "$n"

echo "ten clients, js"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

sleep "$n"

echo "ten clients, lua"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

sleep "$n"

echo "fifty clients, js"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

sleep "$n"

echo "fifty clients, lua"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute

curl -X GET http://localhost:8080/kill
