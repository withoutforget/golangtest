import asyncio
from dataclasses import asdict, dataclass
import datetime
import json
import random
import time
from typing import Any, Callable
import aiohttp
import adaptix
import structlog
import uvloop

logger = structlog.get_logger("test_logger")

class RequestException(Exception):
    def __init__(self, method: str, path: str, status: int, response: aiohttp.ClientResponse) -> None:
        super().__init__("[{method}] {status} | {path} | {response=}".format(method=method, status=status, path = path, response=response))

@dataclass(slots=True)
class TestSettings:
    host: str
    port: int
    secure: bool

    @property
    def url(self) -> str:
        s = 'https' if self.secure else 'http'
        return '{secure}://{host}:{port}'.format(secure = s, host = self.host, port = self.port)
    
class TestClient:
    def __init__(self, session: aiohttp.ClientSession, settings: TestSettings) -> None:
        self.session = session
        self.settings = settings
        self.url = settings.url
    
    async def post(self, path: str, data: Any) -> dict: 
        async with self.session.post(self.url + path, json = data) as response:
            if not (200 <= response.status < 300):
                raise RequestException('post', path, response.status, response)
            return json.loads(await response.text())
    async def get(self, path: str, queries: dict) -> dict:
        async with self.session.get(self.url + path, params = queries) as response:
            if not (200 <= response.status < 300):
                raise RequestException('get', path, response.status, response)
            return json.loads(await response.text())
        
async def generate_multiple_async(
        func: Callable[..., dict],
        count: int,
        generator: Callable[[], list],
) -> list[Callable]:
    return [
        func(*generator())
        for _ in range(count)
    ]

async def test_get(client: aiohttp.ClientSession, count: int):
        data = await generate_multiple_async(
            TestClient.get,
            count = count,
            generator = lambda: [client, '/get', {}]
        )
        start = time.time()
        logger.info("Starting...")

        results = await asyncio.gather(*data, return_exceptions=True)
        finish = time.time()

        executed = sum(1 for i in results if not isinstance(i, BaseException))
        logger.info(f"Успешно отправлено {executed}/{count} запросов ({(finish - start)})")  


@dataclass(slots=True)
class AppendLogRequest:
    raw: str
    level: str 
    created_at: datetime.datetime
    source: str | None
    request_id: str| None
    logger_name: str | None

import string

def random_string(sz: int) -> str:
    return ''.join(random.choice(string.ascii_letters) for _ in range(sz))

async def test_append(client: aiohttp.ClientSession, count: int):
        data = await generate_multiple_async(
            TestClient.post,
            count = count,
            generator = lambda: [client, '/append', asdict(AppendLogRequest(
                raw = random_string(300),
                level = "INFO",
                created_at = datetime.datetime.now(tz = datetime.UTC).isoformat(),
                source = None,
                request_id = None,
                logger_name = None,
            ))]
        )
        logger.info("Starting...")
        start = time.time()
        results = await asyncio.gather(*data, return_exceptions=True)
        finish = time.time()
        executed = sum(1 for i in results if not isinstance(i, BaseException))
        logger.info(f"Успешно отправлено {executed}/{count} запросов ({(finish - start)})")  

async def main():
    connector = aiohttp.TCPConnector(limit = 1000, limit_per_host=1000)
    async with aiohttp.ClientSession(connector=connector) as session:
        settings = TestSettings('localhost', 80, False)
        client = TestClient(session, settings)

        await test_append(client, 10_000)
        await test_get(client, 10_000)


asyncio.run(main())