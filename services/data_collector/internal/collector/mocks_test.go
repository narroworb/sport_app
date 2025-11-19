package collector

import (
	"context"

	"github.com/narroworb/data_collector/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetUnactualTournamentsAndTours(ctx context.Context) ([]UnactualTournamentsAndTours, error) {
	args := m.Called(ctx)
	return args.Get(0).([]UnactualTournamentsAndTours), args.Error(1)
}

func (m *MockDB) GetFootballTournamentID(ctx context.Context, name, season string) (uint32, error) {
	args := m.Called(ctx, name, season)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballTeamID(ctx context.Context, name string) (uint32, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) InsertFootballManager(ctx context.Context, manager *models.Manager) (uint32, error) {
	args := m.Called(ctx, manager)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballManagerID(ctx context.Context, manager *models.Manager) (uint32, error) {
	args := m.Called(ctx, manager)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	args := m.Called(ctx, match, tournamentID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballMatchStatus(ctx context.Context, matchID uint32) (string, error) {
	args := m.Called(ctx, matchID)
	return args.String(0), args.Error(1)
}

func (m *MockDB) UpdateFootballMatch(ctx context.Context, match *models.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockDB) InsertFootballMatch(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	args := m.Called(ctx, match)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballPlayerID(ctx context.Context, name string, position string, height uint16) (uint32, error) {
	args := m.Called(ctx, name, position, height)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) InsertFootballPlayer(ctx context.Context, player *models.Player) (uint32, error) {
	args := m.Called(ctx, player)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) IncrementYellowCardsManager(ctx context.Context, managerID uint32) error {
	args := m.Called(ctx, managerID)
	return args.Error(0)
}

func (m *MockDB) IncrementRedCardsManager(ctx context.Context, managerID uint32) error {
	args := m.Called(ctx, managerID)
	return args.Error(0)
}

func (m *MockDB) InsertFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error) {
	args := m.Called(ctx, stats, matchID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error) {
	args := m.Called(ctx, stats, matchID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) InsertFootballPlayerMatchStats(ctx context.Context, stats models.PlayerStatsInMatch, matchID uint32) (uint32, error) {
	args := m.Called(ctx, stats, matchID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballPlayerMatchStats(ctx context.Context, stats models.PlayerStatsInMatch, matchID uint32) (uint32, error) {
	args := m.Called(ctx, stats, matchID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) InsertFootballGoalieMatchStats(ctx context.Context, stats models.GoalieStatsInMatch, matchID uint32) (uint32, error) {
	args := m.Called(ctx, stats, matchID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballGoalieMatchStats(ctx context.Context, stats models.GoalieStatsInMatch, matchID uint32) (uint32, error) {
	args := m.Called(ctx, stats, matchID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) InsertFootballGoalieMatchStatsBatch(ctx context.Context, statsBatch map[uint32]*models.GoalieStatsInMatch, matchID uint32) error {
	args := m.Called(ctx, statsBatch, matchID)
	return args.Error(0)
}

func (m *MockDB) InsertFootballPlayerMatchStatsBatch(ctx context.Context, statsBatch map[uint32]*models.PlayerStatsInMatch, matchID uint32) error {
	args := m.Called(ctx, statsBatch, matchID)
	return args.Error(0)
}

func (m *MockDB) GetFootballNotPlayedMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	args := m.Called(ctx, match, tournamentID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetFootballPlayedMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	args := m.Called(ctx, match, tournamentID)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockDB) GetCountPlayersStatsByMatchID(ctx context.Context, matchID uint32) (uint64, error) {
	args := m.Called(ctx, matchID)
	return args.Get(0).(uint64), args.Error(1)
}

type MockApi struct {
	mock.Mock
}

func (m *MockApi) FetchBodyConc(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (m *MockApi) FindManagersOfMatch(url string) (homeID, awayID string) {
	args := m.Called(url)
	return args.String(0), args.String(1)
}
