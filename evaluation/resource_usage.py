import time
import psutil
import sys
import threading
import matplotlib.pyplot as plt
from datetime import datetime

def get_process_info(keyword="firecracker"):
    process_dicts = []
    for process in psutil.process_iter(attrs=['pid', 'name', 'cmdline', 'memory_info', 'cpu_percent']):
        try:
            process_info = process.info
            pid = process_info['pid']
            process_name = process_info['name']
            process_cmdline = " ".join(process_info['cmdline'])
            process_memory = process_info['memory_info']
            if keyword in process_cmdline and "python" not in process_name:
                process_dicts.append({
                    'pid': pid,
                    # 'cmdline': process_cmdline,
                    'memory': process_memory.rss,
                })
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            pass
    return process_dicts

def track_cpu_usage(process_dict):
    pid = process_dict['pid']
    try:
        process = psutil.Process(pid)
        process_dict['cpu_percent'] = process.cpu_percent(interval=1)
    except psutil.NoSuchProcess:
        pass

def track_resource_usage(keyword="firecracker"):
    info_list = get_process_info(keyword)
    threads = []

    for process_info in info_list:
        cpu_thread = threading.Thread(target=track_cpu_usage, args=(process_info,))
        cpu_thread.daemon = True
        cpu_thread.start()
        threads.append(cpu_thread)

    for thread in threads:
        thread.join()
    
    print(info_list)
    cpu_sum = 0
    mem_sum = 0
    timestamp = time.time()
    for process_info in info_list:
        cpu_sum += process_info['cpu_percent']
        mem_sum += process_info['memory']

    return {
        'timestamp': datetime.fromtimestamp(timestamp),
        'cpu_percent': cpu_sum,
        'memory': mem_sum / (1024 * 1024),
    }

if __name__ == "__main__":
    n = 50
    if len(sys.argv) > 1:
        n = int(sys.argv[1])

    info_list = []
    for i in range(n):
        info_dict = track_resource_usage()
        info_list.append(info_dict)

    print(info_list)
    timestamps = [info_dict['timestamp'] for info_dict in info_list]
    cpu_usages = [info_dict['cpu_percent'] for info_dict in info_list]
    ram_usages = [info_dict['memory'] for info_dict in info_list]

    # Create a figure with two subplots
    fig, (ax1, ax2) = plt.subplots(2, 1, sharex=True)

    # First subplot for CPU usage
    ax1.set_ylabel('CPU Usage (%)')
    ax1.plot(timestamps, cpu_usages, marker='o', linestyle='-')
    ax1.set_ylim(0, max(cpu_usages))  # Set the y-axis limits to start at 0

    # Second subplot for RAM usage
    ax2.set_xlabel('Timestamp')
    ax2.set_ylabel('RAM Usage (MB)')
    ax2.plot(timestamps, ram_usages, marker='o', linestyle='-')
    plt.xticks(rotation=45)
    plt.tight_layout()

    plt.show()