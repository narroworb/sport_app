from __future__ import annotations

from dataclasses import dataclass

import clickhouse_connect


@dataclass(frozen=True)
class ClickHouseConnInfo:
    host: str
    port: int


def _parse_addr(addr: str) -> ClickHouseConnInfo:
    addr = addr.strip()
    if "://" in addr:
        # Accept http://host:8123 style, but clickhouse-connect wants host/port.
        addr = addr.split("://", 1)[1]
    if "/" in addr:
        addr = addr.split("/", 1)[0]
    if ":" in addr:
        host, port_s = addr.rsplit(":", 1)
        return ClickHouseConnInfo(host=host, port=int(port_s))
    return ClickHouseConnInfo(host=addr, port=8123)


def create_client(database_addr: str, database: str, username: str, password: str):
    info = _parse_addr(database_addr)
    # Use HTTP protocol (8123) to keep runtime simple.
    return clickhouse_connect.get_client(
        host=info.host,
        port=info.port,
        username=username,
        password=password,
        database=database,
        connect_timeout=5,
        send_receive_timeout=20,
    )

