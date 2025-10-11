from dataclasses import asdict, dataclass, field
import datetime
import json
from typing import Any
import requests
import structlog

logger = structlog.get_logger('test_logger')

@dataclass(slots=True, frozen=True, kw_only=True)
class AppendLogRequest:
    raw:str
    level: str
    created_at: datetime.datetime
    source: str | None
    request_id: str | None
    logger_name: str | None

def test_append():

    req = AppendLogRequest(
        raw = "[info] all is okay",
        level = "info",
        created_at = datetime.datetime.now(tz=datetime.UTC).isoformat(),
        source = None,
        request_id=None,
        logger_name=None,
    )
    data = json.dumps(asdict(req), default=str)
    result = requests.post(
        url = 'http://localhost/append',
        data = data,
    )

    resp_data = json.loads(result.text)
    resp_code = result.status_code

    logger.info("Everything is okay!", code = resp_code, data = resp_data)

@dataclass(slots=True, frozen=True, kw_only=True)
class GetLogRequest:
    since: datetime.datetime | None
    before: datetime.datetime | None
    level: list[str]
    source: str | None
    request_id: str | None
    logger_name: str | None
    order: str | None

def test_get():
    req = GetLogRequest(
        since = None,
        before = None,
        level = ["info"],
        source = None,
        request_id=None,
        logger_name=None,
        order = None
    )

    res = requests.get('http://localhost/get', params=asdict(req))

    resp_data = json.loads(res.text)
    resp_code = res.status_code

    logger.info("Everything is okay!", code = resp_code, data = resp_data)

test_append()
test_get()