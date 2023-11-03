
duration=$1
keyword=$2
python3 cpu_usage.py $duration $keyword &
python3 memory_usage.py $duration $keyword