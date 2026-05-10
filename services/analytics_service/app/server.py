import sys
import logging
import os
import signal
import time
from concurrent import futures

import grpc

# Добавляем путь к сгенерированным файлам
sys.path.insert(0, '/app/app/gen')

from analytics.v1 import analytics_pb2, analytics_pb2_grpc

from app.clickhouse_client import create_client
from app.config import load_config
from app.service import AnalyticsEngine



class AnalyticsGrpcService(analytics_pb2_grpc.AnalyticsServiceServicer):
    def __init__(self, engine: AnalyticsEngine):
        self._engine = engine

    def Health(self, request, context):  # noqa: N802
        status, ch = self._engine.health()
        return analytics_pb2.HealthResponse(
            status=status,
            clickhouse=ch,
            server_time_unix=int(time.time()),
        )

    def GetMatchWinProbabilities(self, request, context):  # noqa: N802
        try:
            res = self._engine.get_match_win_probabilities(
                home_team_id=request.home_team_id,
                away_team_id=request.away_team_id,
                season=request.season,
                matches_back=request.matches_back,
                max_goals=request.max_goals,
            )
        except Exception as e:  # noqa: BLE001
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return analytics_pb2.GetMatchWinProbabilitiesResponse()

        return analytics_pb2.GetMatchWinProbabilitiesResponse(
            p_home_win=res.p_home_win,
            p_draw=res.p_draw,
            p_away_win=res.p_away_win,
            lambda_home=res.lambda_home,
            lambda_away=res.lambda_away,
            model="poisson_recent_matches_v1",
            details=res.details,
        )

    def GetTeamFormIndex(self, request, context):  # noqa: N802
        try:
            res = self._engine.get_team_form_index(
                team_id=request.team_id,
                season=request.season,
                matches_back=request.matches_back,
                half_life_matches=request.half_life_matches,
            )
        except Exception as e:  # noqa: BLE001
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return analytics_pb2.GetTeamFormIndexResponse()

        return analytics_pb2.GetTeamFormIndexResponse(
            form_index=res.form_index,
            attack_index=res.attack_index,
            defense_index=res.defense_index,
            trend=res.trend,
            confidence=res.confidence,
            matches_used=res.matches_used,
            details=res.details,
        )

    def GetPlayerSimilarityTopK(self, request, context):  # noqa: N802
        try:
            res = self._engine.get_player_similarity_top_k(
                player_id=request.player_id,
                season=request.season,
                top_k=request.top_k,
                min_minutes=request.min_minutes,
            )
        except Exception as e:  # noqa: BLE001
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return analytics_pb2.GetPlayerSimilarityTopKResponse()

        return analytics_pb2.GetPlayerSimilarityTopKResponse(
            players=[
                analytics_pb2.SimilarPlayer(player_id=p.player_id, similarity=p.similarity)
                for p in res.top
            ],
            position=res.position,
            details=res.details,
        )


def serve() -> None:
    logging.basicConfig(level=os.getenv("LOG_LEVEL", "INFO"))
    cfg = load_config()

    ch = create_client(
        database_addr=cfg.clickhouse_addr,
        database=cfg.clickhouse_db,
        username=cfg.clickhouse_user,
        password=cfg.clickhouse_pass,
    )
    engine = AnalyticsEngine(ch)

    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=int(os.getenv("GRPC_WORKERS", "10"))),
        options=[
            ("grpc.max_receive_message_length", 8 * 1024 * 1024),
            ("grpc.max_send_message_length", 8 * 1024 * 1024),
        ],
    )
    analytics_pb2_grpc.add_AnalyticsServiceServicer_to_server(AnalyticsGrpcService(engine), server)
    server.add_insecure_port(f"0.0.0.0:{cfg.grpc_port}")

    stop = {"now": False}

    def _handle_sig(*_args):
        stop["now"] = True

    signal.signal(signal.SIGTERM, _handle_sig)
    signal.signal(signal.SIGINT, _handle_sig)

    logging.info("analytics_service gRPC listening on :%s", cfg.grpc_port)
    server.start()

    try:
        while not stop["now"]:
            time.sleep(0.2)
    finally:
        logging.info("stopping gRPC server...")
        server.stop(grace=2)


if __name__ == "__main__":
    serve()

