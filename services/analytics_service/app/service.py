from __future__ import annotations

import math
import time
from dataclasses import dataclass

import numpy as np

from app.math_utils import (
    cosine_similarity,
    ewma_weights,
    poisson_score_matrix,
    weighted_avg,
    zscore_normalize,
)
from app.queries import aggregated_player_features_sql, player_position_sql, recent_team_matches_sql


@dataclass(frozen=True)
class MatchWinProbs:
    p_home_win: float
    p_draw: float
    p_away_win: float
    lambda_home: float
    lambda_away: float
    details: str


@dataclass(frozen=True)
class TeamForm:
    form_index: float
    attack_index: float
    defense_index: float
    trend: float
    confidence: float
    matches_used: int
    details: str


@dataclass(frozen=True)
class SimilarPlayerRow:
    player_id: int
    similarity: float


@dataclass(frozen=True)
class PlayerSimilarity:
    position: str
    top: list[SimilarPlayerRow]
    details: str


def _clamp(x: float, lo: float, hi: float) -> float:
    return max(lo, min(hi, x))


def _safe_div(num: float, den: float) -> float:
    return float(num / den) if den else 0.0


class AnalyticsEngine:
    def __init__(self, ch_client):
        self._ch = ch_client

    def health(self) -> tuple[str, str]:
        try:
            self._ch.command("SELECT 1")
            return "ok", "ok"
        except Exception as e:  # noqa: BLE001
            return "degraded", f"error:{e}"

    def get_match_win_probabilities(
        self,
        home_team_id: int,
        away_team_id: int,
        season: str,
        matches_back: int,
        max_goals: int,
    ) -> MatchWinProbs:
        matches_back = int(matches_back) if matches_back else 10
        matches_back = max(3, min(40, matches_back))
        max_goals = int(max_goals) if max_goals else 10
        max_goals = max(5, min(12, max_goals))

        home_rows = self._ch.query(
            recent_team_matches_sql(),
            parameters={"team_id": int(home_team_id), "season": season, "limit": matches_back},
        ).result_rows
        away_rows = self._ch.query(
            recent_team_matches_sql(),
            parameters={"team_id": int(away_team_id), "season": season, "limit": matches_back},
        ).result_rows

        def extract_gf_ga(rows: list[tuple], team_id: int) -> tuple[list[int], list[int]]:
            gf: list[int] = []
            ga: list[int] = []
            for (_mid, _date, h_id, a_id, h_sc, a_sc) in rows:
                if int(h_id) == team_id:
                    gf.append(int(h_sc or 0))
                    ga.append(int(a_sc or 0))
                else:
                    gf.append(int(a_sc or 0))
                    ga.append(int(h_sc or 0))
            return gf, ga

        home_gf, home_ga = extract_gf_ga(home_rows, int(home_team_id))
        away_gf, away_ga = extract_gf_ga(away_rows, int(away_team_id))

        # Simple lambda estimation: blend scoring ability with opponent concession.
        home_attack = float(np.mean(home_gf)) if home_gf else 1.2
        home_defense_conc = float(np.mean(home_ga)) if home_ga else 1.2
        away_attack = float(np.mean(away_gf)) if away_gf else 1.2
        away_defense_conc = float(np.mean(away_ga)) if away_ga else 1.2

        # Add a mild home advantage and stabilize with priors.
        base_prior = 1.25
        lambda_home = 0.55 * home_attack + 0.35 * away_defense_conc + 0.10 * base_prior
        lambda_away = 0.55 * away_attack + 0.35 * home_defense_conc + 0.10 * base_prior
        lambda_home *= 1.08

        lambda_home = _clamp(lambda_home, 0.2, 3.5)
        lambda_away = _clamp(lambda_away, 0.2, 3.5)

        mat = poisson_score_matrix(lambda_home, lambda_away, max_goals)
        p_draw = float(np.trace(mat))
        p_home_win = float(np.tril(mat, k=-1).sum())
        p_away_win = float(np.triu(mat, k=1).sum())

        # Renormalize in case of tail truncation.
        s = p_home_win + p_draw + p_away_win
        if s > 0:
            p_home_win, p_draw, p_away_win = p_home_win / s, p_draw / s, p_away_win / s

        details = (
            f"recent_matches={matches_back}, max_goals={max_goals}, "
            f"home_attack={home_attack:.2f}, away_attack={away_attack:.2f}"
        )
        return MatchWinProbs(
            p_home_win=p_home_win,
            p_draw=p_draw,
            p_away_win=p_away_win,
            lambda_home=lambda_home,
            lambda_away=lambda_away,
            details=details,
        )

    def get_team_form_index(
        self,
        team_id: int,
        season: str,
        matches_back: int,
        half_life_matches: float,
    ) -> TeamForm:
        matches_back = int(matches_back) if matches_back else 10
        matches_back = max(3, min(30, matches_back))
        half_life_matches = float(half_life_matches) if half_life_matches else 5.0
        half_life_matches = _clamp(half_life_matches, 1.0, 20.0)

        rows = self._ch.query(
            recent_team_matches_sql(),
            parameters={"team_id": int(team_id), "season": season, "limit": matches_back},
        ).result_rows

        points: list[float] = []
        gf: list[float] = []
        ga: list[float] = []

        for (_mid, _date, h_id, a_id, h_sc, a_sc) in rows:
            h_id = int(h_id)
            a_id = int(a_id)
            h_sc = int(h_sc or 0)
            a_sc = int(a_sc or 0)
            is_home = h_id == int(team_id)
            team_goals = h_sc if is_home else a_sc
            opp_goals = a_sc if is_home else h_sc
            gf.append(float(team_goals))
            ga.append(float(opp_goals))
            if team_goals > opp_goals:
                points.append(3.0)
            elif team_goals == opp_goals:
                points.append(1.0)
            else:
                points.append(0.0)

        n = len(points)
        if n == 0:
            return TeamForm(
                form_index=0.0,
                attack_index=0.0,
                defense_index=0.0,
                trend=0.0,
                confidence=0.0,
                matches_used=0,
                details="no matches found",
            )

        w = ewma_weights(n, half_life_matches)
        # NOTE: rows are ordered newest->oldest, ewma_weights assumes ages 0..n-1 (newest age=0).
        pts_norm = [p / 3.0 for p in points]
        form_index = weighted_avg(pts_norm, w)

        # Attack/defense: normalize vs soft priors to fit 0..1.
        avg_gf = weighted_avg(gf, w)
        avg_ga = weighted_avg(ga, w)
        attack_index = _clamp(avg_gf / 3.0, 0.0, 1.0)
        defense_index = _clamp(1.0 - (avg_ga / 3.0), 0.0, 1.0)

        # Trend: compare first half vs second half (newest vs oldest) in weighted points.
        mid = max(1, n // 2)
        newest_pts = float(np.mean(pts_norm[:mid]))
        oldest_pts = float(np.mean(pts_norm[mid:])) if n > mid else newest_pts
        trend = _clamp(newest_pts - oldest_pts, -1.0, 1.0)

        confidence = _clamp(math.sqrt(n / matches_back), 0.0, 1.0)

        details = f"matches_used={n}, half_life={half_life_matches:.1f}, avg_gf={avg_gf:.2f}, avg_ga={avg_ga:.2f}"
        return TeamForm(
            form_index=float(form_index),
            attack_index=float(attack_index),
            defense_index=float(defense_index),
            trend=float(trend),
            confidence=float(confidence),
            matches_used=n,
            details=details,
        )

    def get_player_similarity_top_k(
        self,
        player_id: int,
        season: str,
        top_k: int,
        min_minutes: int,
    ) -> PlayerSimilarity:
        top_k = int(top_k) if top_k else 10
        top_k = max(1, min(30, top_k))
        min_minutes = int(min_minutes) if min_minutes else 600
        min_minutes = max(0, min(6000, min_minutes))

        pos_rows = self._ch.query(
            player_position_sql(), parameters={"player_id": int(player_id)}
        ).result_rows
        position = str(pos_rows[0][0]) if pos_rows else ""

        rows = self._ch.query(
            aggregated_player_features_sql(),
            parameters={"season": season, "min_minutes": int(min_minutes)},
        ).result_rows

        if not rows:
            return PlayerSimilarity(position=position, top=[], details="no player aggregates found")

        # Build feature matrix.
        ids: list[int] = []
        pos: list[str] = []
        feats: list[list[float]] = []

        for (
            athlete_id,
            p,
            minutes,
            rating_avg,
            goals,
            assists,
            shots_on_target,
            total_shots,
            key_passes,
            duels,
            duels_won,
            pass_attempts,
            complete_passes,
            interceptions,
            total_tackles,
        ) in rows:
            minutes = float(minutes or 0.0)
            per90 = 90.0 / minutes if minutes > 0 else 0.0
            pass_acc = _safe_div(float(complete_passes or 0.0), float(pass_attempts or 0.0))
            duel_win = _safe_div(float(duels_won or 0.0), float(duels or 0.0))

            # Vector is intentionally small/stable; per90 makes players comparable.
            v = [
                float(rating_avg or 0.0),
                float(goals or 0.0) * per90,
                float(assists or 0.0) * per90,
                float(shots_on_target or 0.0) * per90,
                float(total_shots or 0.0) * per90,
                float(key_passes or 0.0) * per90,
                pass_acc,
                duel_win,
                float(interceptions or 0.0) * per90,
                float(total_tackles or 0.0) * per90,
            ]
            ids.append(int(athlete_id))
            pos.append(str(p or ""))
            feats.append(v)

        mat = np.asarray(feats, dtype=float)
        mat_z, _, _ = zscore_normalize(mat)

        try:
            idx_query = ids.index(int(player_id))
        except ValueError:
            return PlayerSimilarity(position=position, top=[], details="query player not found in aggregates")

        query_pos = pos[idx_query] or position
        qv = mat_z[idx_query]

        scored: list[SimilarPlayerRow] = []
        for i, pid in enumerate(ids):
            if pid == int(player_id):
                continue
            if query_pos and (pos[i] != query_pos):
                continue
            sim = cosine_similarity(qv, mat_z[i])
            scored.append(SimilarPlayerRow(player_id=pid, similarity=sim))

        scored.sort(key=lambda x: x.similarity, reverse=True)
        top = scored[:top_k]
        details = f"candidates={len(scored)}, top_k={top_k}, min_minutes={min_minutes}, position={query_pos}"
        return PlayerSimilarity(position=query_pos, top=top, details=details)

