import asyncio
import json
import random
import sys
import time
from typing import List
import aiohttp

url = "http://localhost:8080/execute"

async def make_request(session: aiohttp.ClientSession, submission: str, non_200_responses: List[str], lock: asyncio.Lock, debug: bool = False):
    start_time = time.time()
    log(f"Making request to {url}", debug=debug)
    # send a post request with submission to the URL
    async with session.post(url, data=submission) as response:
        response_text = await response.text()
        end_time = time.time()
        log(f"Time taken: {end_time - start_time} seconds", debug=debug)
        log(f"Response: {response_text}", debug=debug)
        async with lock:
            if response.status != 200:
                non_200_responses.append(response_text)

async def main(concurrent_requests=5, submission=None, debug: bool = False):
    async with aiohttp.ClientSession() as session:
        non_200_responses = []
        lock = asyncio.Lock()
        tasks = [make_request(session, submission, non_200_responses, lock) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)
        
        print(f"Number of non-200 responses: {len(non_200_responses)}")
        print("Non-200 responses:")
        for response in non_200_responses:
            print(response)

def log(message: str, debug: bool):
    if debug:
        print(message)
            
def run(concurrent_requests=5, submission=None, debug=False):
    time_start = time.time()
    asyncio.run(main(concurrent_requests, submission, debug=debug))
    time_end = time.time()
    print(f"Time taken: {time_end - time_start} seconds for {concurrent_requests} requests")
    print(f"Average time per request: {(time_end - time_start) / concurrent_requests} seconds")
    print(f"Throughput: {concurrent_requests / (time_end - time_start)} requests per second")