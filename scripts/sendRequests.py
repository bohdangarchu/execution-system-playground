import asyncio
import aiohttp

# Number of concurrent requests
concurrent_requests = 10

# URL to request
url = "http://localhost:8080/"

async def make_request(session):
    print(f"Making request to {url}")
    async with session.get(url) as response:
        response_text = await response.text()
        print(f"Response: {response_text}")

async def main():
    async with aiohttp.ClientSession() as session:
        tasks = [make_request(session) for _ in range(concurrent_requests)]
        await asyncio.gather(*tasks)

if __name__ == "__main__":
    asyncio.run(main())
