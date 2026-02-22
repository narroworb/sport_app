package handlers

import (
	"context"
	"time"

	"github.com/narroworb/core_api/internal/models"
)

type DatabaseInterface interface {
	GetPlayerByID(ctx context.Context, id uint32) (models.Player, error)
	GetManagerByID(ctx context.Context, id uint32) (models.Manager, error)
	GetTeamByID(ctx context.Context, id uint32) (models.Team, error)
	GetMatchByID(ctx context.Context, id uint32) (models.Match, error)
	GetTournamentByID(ctx context.Context, id uint32) (models.Tournament, error)
	GetPlayerPositionByID(ctx context.Context, id uint32) (string, error)
	GetPlayerFullStats(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	GetPlayerStatsBySeason(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	GetPlayerStatsByDates(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	GetGoalieFullStats(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	GetGoalieStatsBySeason(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	GetGoalieStatsByDates(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	GetPlayerFixtures(ctx context.Context, id, limit, offset uint32) ([]models.PlayerMatch, error)
	GetPlayerTeams(ctx context.Context, id uint32) ([]models.Team, error)
	GetTeamStatsBySeason(ctx context.Context, filter models.TeamStatsFilter) (models.TeamStatsInPeriod, error)
	GetTeamStatsByDates(ctx context.Context, filter models.TeamStatsFilter) (models.TeamStatsInPeriod, error)
	GetTeamFullStats(ctx context.Context, filter models.TeamStatsFilter) (models.TeamStatsInPeriod, error)
	GetStandingsByTeamAndSeason(ctx context.Context, teamID uint32, season string) ([]models.TableRow, error)
	GetTeamNextGame(ctx context.Context, id uint32) (models.ShortMatch, error)
	GetTeamPlayersBySeason(ctx context.Context, teamID uint32, season string) (map[string][]models.Player, error)
	GetTeamLastGames(ctx context.Context, teamID uint32, limit, offset uint32) ([]models.ShortMatch, error)
	GetCurrentManagerByTeam(ctx context.Context, id uint32) (models.Manager, error)
	GetTeamPlayersWithStatsBySeason(ctx context.Context, teamID uint32, season string) (map[string][]models.PlayerWithStats, error)
	GetTournamentTableByID(ctx context.Context, tournamentID uint32) ([]models.TableRow, error)
	GetTournamentFixturesByID(ctx context.Context, id uint32) (map[uint16][]models.ShortMatch, error)
	GetTournamentTeamsStatsByID(ctx context.Context, id uint32) ([]models.TeamStatsInPeriod, error)
	GetTournamentPlayersStatsByID(ctx context.Context, id, limit, offset uint32) ([]models.PlayerStatsInPeriod, error)
	GetTournamentTableGraphByID(ctx context.Context, tournamentID uint32) ([][]models.ShortTableRow, error)
	GetManagerFullStats(ctx context.Context, id uint32) (map[string]models.ManagerStatsInPeriod, error)
	GetManagerStatsBySeason(ctx context.Context, id uint32) (map[string]models.ManagerStatsInPeriod, error)
	GetManagerStatsByTeam(ctx context.Context, id uint32) (map[string]models.ManagerStatsInPeriod, error)
	GetManagerCardsByID(ctx context.Context, id uint32) (uint16, uint16, error)
	GetManagerTeams(ctx context.Context, id uint32) (map[string][]string, error)
	GetManagerFixtures(ctx context.Context, id, limit, offset uint32) ([]models.ManagerMatch, error)
	GetMatchesByDate(ctx context.Context, date time.Time) ([]models.Match, error)
	GetMatchPlayersStats(ctx context.Context, id uint32) (map[string][]models.PlayerStatsInMatch, error)
	GetMatchGoaliesStats(ctx context.Context, id uint32) (map[string][]models.GoalieStatsInMatch, error)
	GetMatchTeamsStats(ctx context.Context, id uint32) (models.TeamMatchStats, error)
}

type CacheInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}
