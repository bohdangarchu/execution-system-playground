import asyncio
import random
import time
import aiohttp

# Number of concurrent requests
concurrent_requests = 3

# URL to request
url = "http://localhost:8081/"

async def make_request(session: aiohttp.ClientSession):
    print(f"Making request to {url}")
    # send a post reqeust with submission to the url
    async with session.post(url, data=get_random_submission()) as response:
        response_text = await response.text()
        print(f"Response: {response_text}")

async def main():
    async with aiohttp.ClientSession() as session:
        tasks = [make_request(session) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)

def get_random_submission():
    val1 = random.randint(-100, 100)
    val2 = random.randint(-100, 100)
    return r"""
{
	"functionName": "addTwoNumbers",
	"code": "function addTwoNumbers(a, b) {\n  return a + b;\n}",
	"testCases": [
	  {
		"input": [
		  {
			"value": """ + str(val1) + """,
			"type": "number"
		  },
		  {
			"value": """ + str(val2) + """,
			"type": "number"
		  }
		]
	  }
	]
}
"""

if __name__ == "__main__":
    # measure time
    time_start = time.time()
    asyncio.run(main())
    time_end = time.time()
    print(f"Time taken: {time_end - time_start} seconds for {concurrent_requests} requests")
    print(f"Average time per request: {(time_end - time_start) / concurrent_requests} seconds")
