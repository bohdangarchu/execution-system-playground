import asyncio
import json
import random
import sys
import time
import aiohttp

# TODO change json format
url = "http://localhost:8080/execute"

async def make_request(session: aiohttp.ClientSession):
    print(f"Making request to {url}")
    # send a post reqeust with submission to the url
    async with session.post(url, data=get_heavy_submission()) as response:
        response_text = await response.text()
        print(f"Response: {response_text}")

async def main(concurrent_requests=5):
    async with aiohttp.ClientSession() as session:
        tasks = [make_request(session) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)

def get_heavy_submission():
    fun = """function estimatePI(iterations) {
  let insideCircle = 0;
  for (let i = 0; i < iterations; i++) {
    const x = Math.random();
    const y = Math.random();
    const distance = Math.sqrt(x * x + y * y);
    if (distance <= 1) {
      insideCircle++;
    }
  }
  return 4 * (insideCircle / iterations);
}
"""
    sub = {
        "functionName": "estimatePI",
        "code": fun,
        "testCases": [
            {
                "input": ["1000000"],
            }
        ]
    }
    return json.dumps(sub, indent = 4)


if __name__ == "__main__":
    concurrent_requests = 5
    if len(sys.argv) > 1:
        concurrent_requests = int(sys.argv[1])

    time_start = time.time()
    asyncio.run(main(concurrent_requests))
    time_end = time.time()
    print(f"Time taken: {time_end - time_start} seconds for {concurrent_requests} requests")
    print(f"Average time per request: {(time_end - time_start) / concurrent_requests} seconds")

