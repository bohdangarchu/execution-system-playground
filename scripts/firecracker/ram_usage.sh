#!/bin/sh

# start execution server in the background
# main: firerunner.RunStandaloneVM()

ps_output="ps_output.txt"
duration=30
interval=1

rm -f "$ps_output"  # Clear previous output file (if any)

# Run the command and save the output to the file
timeout "$duration"s bash -c "while true; do ps aux | grep firecracker | grep -v 'grep' | awk '{print \$6}'; sleep $interval; done" >> "$ps_output"

gnuplot << EOF
set term dumb
set title "RAM Usage of Firecracker Process"
set xlabel "Time (seconds)"
set ylabel "RAM Usage Kb"
plot "$ps_output" with lines
pause -1 "Press enter to exit"
EOF
