#!/bin/sh

# start execution server in the background
docker run -p 8080:8080 execution-server &

ps_output="ps_output.txt"
duration=20
interval=1

# Clear previous output file (if any)
rm -f "$ps_output"  

# Run the command and save the output to the file
timeout "$duration"s bash -c "while true; do ps aux | grep 'execution-server' | grep -v 'grep' | awk '{print \$3}'; sleep $interval; done" >> "$ps_output"

gnuplot << EOF
set term dumb
set title "CPU Usage of Docker Process"
set xlabel "Time (seconds)"
set ylabel "CPU Usage (%)"
plot "$ps_output" with lines
pause -1 "Press enter to exit"
EOF
