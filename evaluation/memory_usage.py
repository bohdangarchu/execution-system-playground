import time
import psutil
import sys
import threading
import matplotlib.pyplot as plt
from datetime import datetime

interval = 0.25

filter_keywords = ['python', 'stats', 'grep']

def is_valid(process_cmdline):
    for keyword in filter_keywords:
        if keyword in process_cmdline:
            return False
    return True

def get_process_info(keyword):
    process_dicts = []
    for process in psutil.process_iter(attrs=['pid', 'name', 'cmdline', 'memory_info', 'cpu_percent']):
        try:
            process_info = process.info
            pid = process_info['pid']
            process_name = process_info['name']
            process_cmdline = " ".join(process_info['cmdline'])
            process_memory = process_info['memory_info']
            if keyword in process_cmdline and is_valid(process_cmdline):
                # print(process_memory)
                process_dicts.append({
                    'pid': pid,
                    'name': process_name,
                    'cmdline': process_cmdline,
                    'memory': process_memory.rss,
                })
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            pass
    return process_dicts

def track_resource_usage(keyword):
    info_list = get_process_info(keyword)

    if len(info_list) == 0:
        return {
            'timestamp': datetime.fromtimestamp(time.time()),
            'memory': 0,
        }

    print(info_list)
    mem_sum = 0
    timestamp = time.time()
    for process_info in info_list:
        mem_sum += process_info.get('memory', 0)

    return {
        'timestamp': datetime.fromtimestamp(timestamp),
        'memory': mem_sum / (1024 * 1024),
    }

if __name__ == "__main__":
    duration = 20
    if len(sys.argv) > 1:
        duration = int(sys.argv[1])

    info_list = []
    for i in range(int(duration / interval)):
        info_dict = track_resource_usage('execution-system-bin')
        info_list.append(info_dict)
        time.sleep(interval)

    # print(info_list)
    timestamps = [info_dict['timestamp'] for info_dict in info_list]
    ram_usages = [info_dict['memory'] for info_dict in info_list]

    plt.plot(timestamps, ram_usages)
    plt.xlabel('Timestamp')
    plt.ylabel('RAM Usage (MB)')
    plt.title('RAM Usage over Time')
    # start at y=0
    axes = plt.gca()
    axes.set_ylim([0, None])

    plt.show()