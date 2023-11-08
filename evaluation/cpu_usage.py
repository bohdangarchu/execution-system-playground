import time
import psutil
import sys
import threading
import matplotlib.pyplot as plt
from datetime import datetime

interval = 0.1

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
            if not (process_info['cmdline'] and isinstance(process_info['cmdline'], list)):
                # print(process_info)
                continue
            process_cmdline = " ".join(process_info['cmdline'])
            process_memory = process_info['memory_info']
            if keyword in process_cmdline and is_valid(process_cmdline):
                print(process_memory)
                process_dicts.append({
                    'pid': pid,
                    'name': process_name,
                    'cmdline': process_cmdline,
                    'memory': process_memory.rss,
                })
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            pass
    return process_dicts

def track_cpu_usage(process_dict):
    pid = process_dict['pid']
    try:
        process = psutil.Process(pid)
        process_dict['cpu_percent'] = process.cpu_percent(interval=interval) / psutil.cpu_count(logical=True)
    except psutil.NoSuchProcess:
        pass

def track_resource_usage(keyword):
    info_list = get_process_info(keyword)
    threads = []

    if len(info_list) == 0:
        time.sleep(interval)
        return {
            'timestamp': time.time(),
            'cpu_percent': 0,
            'memory': 0,
        }

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
    for process_info in info_list:
        cpu_sum += process_info.get('cpu_percent', 0)
        mem_sum += process_info.get('memory', 0)

    return {
        'timestamp': time.time(),
        'cpu_percent': cpu_sum,
        'memory': mem_sum / (1024 * 1024),
    }

if __name__ == "__main__":
    duration = 50
    if len(sys.argv) > 2:
        duration = int(sys.argv[1])
        keyword = sys.argv[2]

    timestamp = time.time()
    info_list = []
    for i in range(int(duration/interval)):
        info_dict = track_resource_usage(keyword)
        info_list.append(info_dict)
        if time.time() - timestamp > duration:
            break

    timestamps = [info_dict['timestamp'] - timestamp for info_dict in info_list]
    cpu_usages = [info_dict['cpu_percent'] for info_dict in info_list]
    print(timestamps)
    print(cpu_usages)

    plt.plot(timestamps, cpu_usages)
    plt.xlabel('Time (s)')
    plt.ylabel('CPU Usage (%)')
    plt.title('CPU Usage over Time')
    # start at y=0
    axes = plt.gca()
    axes.set_ylim([0, None])
    plt.show()