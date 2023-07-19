#!/bin/sh

ps_output="ps_output.txt"
duration=20
interval=1

rm -f "$ps_output"  # Clear previous output file (if any)

# Run the command and save the output to the file
timeout "$duration"s bash -c "while true; do ps aux | grep firecracker | grep -v 'grep' | awk '{print \$3}'; sleep $interval; done" >> "$ps_output"

gnuplot << EOF
set term dumb
set title "CPU Usage of Firecracker Process"
set xlabel "Time (seconds)"
set ylabel "CPU Usage (%)"
plot "$ps_output" with lines
pause -1 "Press enter to exit"
EOF
