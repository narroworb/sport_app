from __future__ import annotations


def recent_team_matches_sql() -> str:
    return """
SELECT
  match_id,
  date,
  home_team_id,
  away_team_id,
  home_score,
  away_score
FROM Matches m
INNER JOIN Tournaments t ON t.tournament_id = m.tournament_id
WHERE (m.home_team_id = {team_id:UInt32} OR m.away_team_id = {team_id:UInt32})
  AND m.status = 'Ended'
  AND t.season = {season:String}
ORDER BY m.date DESC
LIMIT {limit:UInt32}
"""


def player_position_sql() -> str:
    return """
SELECT position
FROM Athletes
WHERE athlete_id = {player_id:UInt32}
LIMIT 1
"""


def aggregated_player_features_sql() -> str:
    # This query intentionally stays conservative and uses columns already visible in core_api queries.
    return """
SELECT
  s.athlete_id,
  any(a.position) AS position,
  sum(s.minutes_played) AS minutes,
  avgIf(toFloat64(s.rating), s.rating != 0) AS rating_avg,
  sum(s.goals) AS goals,
  sum(s.assists) AS assists,
  sum(s.shot_on_target) AS shots_on_target,
  sum(s.total_shots) AS total_shots,
  sum(s.key_passes) AS key_passes,
  sum(s.duels) AS duels,
  sum(s.duels_won) AS duels_won,
  sum(s.pass_attempts) AS pass_attempts,
  sum(s.complete_passes) AS complete_passes,
  sum(s.interceptions) AS interceptions,
  sum(s.total_tackles) AS total_tackles
FROM Football_Player_Match_Stats s
INNER JOIN Matches m ON m.match_id = s.match_id
INNER JOIN Tournaments t ON t.tournament_id = m.tournament_id
INNER JOIN Athletes a ON a.athlete_id = s.athlete_id
WHERE t.season = {season:String}
  AND m.status = 'Ended'
GROUP BY s.athlete_id
HAVING minutes >= {min_minutes:UInt32}
"""

