echo "js n = 1"
wrk -s js_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute
echo "lua n = 1"
wrk -s lua_submission.lua -t1 -c1 -d60s --timeout 60s http://localhost:8080/execute

echo "js n = 10"
wrk -s js_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute
echo "lua n = 10"
wrk -s lua_submission.lua -t1 -c10 -d60s --timeout 60s http://localhost:8080/execute

echo "js n = 50"
wrk -s js_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute
echo "lua n = 50"
wrk -s lua_submission.lua -t1 -c50 -d60s --timeout 60s http://localhost:8080/execute
