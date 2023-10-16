import asyncio
import json
import random
import sys
import time
import aiohttp

url = "http://localhost:8080/execute"

async def make_request(session: aiohttp.ClientSession, non_200_responses, lock):
    print(f"Making request to {url}")
    submission = get_random_submission()
    # send a post request with submission to the URL
    async with session.post(url, data=submission) as response:
        response_text = await response.text()
        async with lock:
            if response.status != 200:
                non_200_responses.append(response_text)

async def main(concurrent_requests=5):
    async with aiohttp.ClientSession() as session:
        non_200_responses = []
        lock = asyncio.Lock()
        tasks = [make_request(session, non_200_responses, lock) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)
        
        for response in non_200_responses:
            print(response)


def get_random_submission():
    val = random.randint(-1000, 1000)
    return r"""
{
	"functionName": "addTwoNumbers",
	"code": "function addTwoNumbers(a, b) {\n  return a + b + """ + str(val) + r""";\n}",
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
