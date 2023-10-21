import asyncio
import json
import random
import sys
import time
import aiohttp

url = "http://localhost:8080/execute"

async def make_request(session: aiohttp.ClientSession, non_200_responses, lock):
    start_time = time.time()
    print(f"Making request to {url}")
    submission = get_submission()
    # send a post request with submission to the URL
    async with session.post(url, data=submission) as response:
        response_text = await response.text()
        end_time = time.time()
        print(f"Time taken: {end_time - start_time} seconds")
        print(f"Response: {response_text}")
        async with lock:
            if response.status != 200:
                non_200_responses.append(response_text)

async def main(concurrent_requests=5):
    async with aiohttp.ClientSession() as session:
        non_200_responses = []
        lock = asyncio.Lock()
        tasks = [make_request(session, non_200_responses, lock) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)
        
        print(f"Number of non-200 responses: {len(non_200_responses)}")
        print("Non-200 responses:")
        for response in non_200_responses:
            print(response)


def get_submission():
    return r"""
{
	"functionName": "addTwoNumbers",
	"code": "function addTwoNumbers(a, b) {\n  return a + b;\n}",
    "language": "js",
	"testCases": [
	  {
        "id": "1",
		"input": ["1", "2"]
	  }
	]
}
"""

if __name__ == "__main__":
    concurrent_requests = 5
    if len(sys.argv) > 1:
        concurrent_requests = int(sys.argv[1])

    time_start = time.time()
    asyncio.run(main(concurrent_requests))
    time_end = time.time()
    print(f"Time taken: {time_end - time_start} seconds for {concurrent_requests} requests")
    print(f"Average time per request: {(time_end - time_start) / concurrent_requests} seconds")
