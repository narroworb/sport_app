import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Config:
    clickhouse_addr: str
    clickhouse_db: str
    clickhouse_user: str
    clickhouse_pass: str
    grpc_port: int


def load_config() -> Config:
    addr = os.getenv("DATABASE_ADDR", "clickhouse:8123").strip()
    db = os.getenv("DB_NAME", "default").strip()
    user = os.getenv("DB_USER", "default").strip()
    password = os.getenv("DB_PASS", "")
    grpc_port = int(os.getenv("ANALYTICS_GRPC_PORT", "50051"))

    if not addr:
        addr = "clickhouse:8123"

    return Config(
        clickhouse_addr=addr,
        clickhouse_db=db or "default",
        clickhouse_user=user or "default",
        clickhouse_pass=password,
        grpc_port=grpc_port,
    )

