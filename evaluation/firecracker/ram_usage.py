import time
import psutil

def get_memory_and_cpu_usage():
    total_memory = 0
    total_cpu = 0

    for process in psutil.process_iter(attrs=['pid', 'name', 'cmdline', 'memory_info', 'cpu_percent']):
        try:
            process_info = process.info
            process_name = process_info['name']
            process_cmdline = " ".join(process_info['cmdline'])
            process_memory = process_info['memory_info']
            process_cpu = process_info['cpu_percent']

            if "firecracker" in process_cmdline and "python" not in process_name:
                # print(f"Process: {process_cmdline}")
                total_memory += process_memory.rss
                total_cpu += process_cpu
                p = psutil.Process(process_info['pid'])
                cpu = p.cpu_percent(interval=1)
                print(f'pid {p.pid} cpu: {cpu:.2f} %')

        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            pass
    
    # Convert bytes to megabytes
    return total_memory / (1024 * 1024), total_cpu

if __name__ == "__main__":
    startTime = time.time()
    for i in range(100):
        total_mem, total_cpu = get_memory_and_cpu_usage()
        # print(f"Total memory usage of 'firecracker' processes: {total_mem:.2f} MB")
    print(f"--- {time.time() - startTime} seconds ---")
