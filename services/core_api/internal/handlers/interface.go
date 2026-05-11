package handlers

import (
	"context"
	"time"

	"github.com/narroworb/core_api/internal/models"
)

type AnalyticDatabaseInterface interface {
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

type TransactionDatabaseInterface interface {
	GetFavoritePlayersIDs(ctx context.Context, userID int64) ([]uint32, error)
	GetFavoriteManagersIDs(ctx context.Context, userID int64) ([]uint32, error)
	GetFavoriteTeamsIDs(ctx context.Context, userID int64) ([]uint32, error)
	GetFavoriteTournamentIDs(ctx context.Context, userID int64) ([]uint32, error)

	SetFavoritePlayerByID(ctx context.Context, userID int64, playerID int64) error
	SetFavoriteManagerByID(ctx context.Context, userID int64, managerID int64) error
	SetFavoriteTeamByID(ctx context.Context, userID int64, teamID int64) error
	SetFavoriteTournamentByID(ctx context.Context, userID int64, tournamentID int64) error

	DeleteFavoritePlayerByID(ctx context.Context, userID int64, playerID int64) error
	DeleteFavoriteManagerByID(ctx context.Context, userID int64, managerID int64) error
	DeleteFavoriteTeamByID(ctx context.Context, userID int64, teamID int64) error
	DeleteFavoriteTournamentByID(ctx context.Context, userID int64, tournamentID int64) error
}

type SearchDatabaseInterface interface {
	Search(ctx context.Context, query string, entityType string, filters models.SearchFilters) ([]models.SearchResult, error)
	IndexPlayer(ctx context.Context, player models.Player) error
	IndexTeam(ctx context.Context, team models.Team) error
	IndexManager(ctx context.Context, manager models.Manager) error
	IndexTournament(ctx context.Context, tournament models.Tournament) error
	DeleteIndex(ctx context.Context, index string) error
}

