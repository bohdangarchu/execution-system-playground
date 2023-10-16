import asyncio
import json
import random
import sys
import time
import aiohttp

url = "http://localhost:8080/execute"

async def make_request(session: aiohttp.ClientSession, response_status_codes, lock):
    print(f"Making request to {url}")
    submission = get_random_submission()
    # send a post request with submission to the URL
    async with session.post(url, data=submission) as response:
        response_text = await response.text()
        print(f"Response: {response_text} for submission {submission}")
        async with lock:
            response_status_codes.append(response.status)

async def main(concurrent_requests=5):
    async with aiohttp.ClientSession() as session:
        response_status_codes = []
        lock = asyncio.Lock()
        tasks = [make_request(session, response_status_codes, lock) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)
        
        # Count how many responses had a status code different from 200
        non_200_count = len([status for status in response_status_codes if status != 200])
        print(f"Number of responses with status code other than 200: {non_200_count}")


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
