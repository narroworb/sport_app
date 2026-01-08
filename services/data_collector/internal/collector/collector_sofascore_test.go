package collector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/narroworb/data_collector/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchManager(t *testing.T) {
	api := new(MockApi)
	db := new(MockDB)
	updater := &Updater{db: db, api: api}

	jsonBody := `xxxxx{"manager": {"name": "Pep Guardiola","country": {"name": "Spain"}}}`

	api.On("FetchBodyConc", mock.Anything, "https://api.sofascore.com/api/v1/manager/123").
		Return(jsonBody, nil)

	mgr, err := updater.fetchManager(context.Background(), "123")

	require.NoError(t, err)
	assert.Equal(t, "Pep", mgr.FirstName)
	assert.Equal(t, "Guardiola", mgr.LastName)
	assert.Equal(t, "Spain", mgr.Nation)
}

func TestFetchMatches(t *testing.T) {
	api := new(MockApi)
	db := new(MockDB)
	u := &Updater{db: db, api: api}

	body := `xxxxx{
		"events":[{
			"id":1,
			"customID":"abc",
			"slug":"team1-team2",
			"startTimestamp":1700000000,
			"homeTeam":{"name":"A"},
			"awayTeam":{"name":"B"},
			"homeScore":{"display":2},
			"awayScore":{"display":1},
			"status":{"description":"finished"}
		}]}`
	api.On("FetchBodyConc", mock.Anything, "url/events/round/10").
		Return(body, nil)

	db.On("GetMatchesByTournamentAndRound", mock.Anything, mock.Anything, mock.Anything).
		Return(map[ShortTypeMatch]ManagersOfMatch{}, nil)

	api.On("FindManagersOfMatch", mock.Anything, mock.AnythingOfType("string")).
		Return("10", "20")

	api.On("FetchBodyConc", mock.Anything, "https://api.sofascore.com/api/v1/manager/10").
		Return(`xxxxx{"manager":{"name":"John Doe","country":{"name":"England"}}}`, nil)

	api.On("FetchBodyConc", mock.Anything, "https://api.sofascore.com/api/v1/manager/20").
		Return(`xxxxx{"manager":{"name":"Carlos Ruiz","country":{"name":"Mexico"}}}`, nil)

	db.On("GetFootballManagerID", mock.Anything, mock.Anything).
		Return(uint32(1), nil)

	teams := map[string]*models.Team{
		"A": {Name: "A", ID: 1},
		"B": {Name: "B", ID: 2},
	}

	matches, err := u.fetchMatches(context.Background(), "url", 10, teams, 1)

	require.NoError(t, err)
	require.Len(t, matches, 1)

	m := matches[0]
	assert.Equal(t, uint16(2), m.HomeGoals)
	assert.Equal(t, uint16(1), m.AwayGoals)
	assert.Equal(t, "finished", m.Status)
	assert.Equal(t, "A", m.HomeTeam.Name)
	assert.Equal(t, "B", m.AwayTeam.Name)
}

func TestFetchAllStatisticsFromMatches(t *testing.T) {
	api := new(MockApi)
	db := new(MockDB)

	u := &Updater{
		api: api,
		db:  db,
	}

	ctx := context.Background()
	matchID := "123"

	urlPlayers := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%s/lineups", matchID)
	urlTeams := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%s/statistics", matchID)
	urlIncidents := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%s/incidents", matchID)

	api.On("FetchBodyConc", mock.Anything, urlTeams).
		Return(teamStatsJSON, nil)

	api.On("FetchBodyConc", mock.Anything, urlPlayers).
		Return(playersStatsJSON, nil)

	api.On("FetchBodyConc", mock.Anything, urlIncidents).
		Return(incidentsJSON, nil)

	db.On("GetFootballPlayerID", mock.Anything, "Pepe Reina", time.Unix(399600000, 0)).
		Return(uint32(26937), nil)

	db.On("GetFootballPlayerID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(uint32(110), nil)

	stats, err := u.fetchAllStatisticsFromMatches(ctx, matchID)
	require.NoError(t, err)

	assert.Equal(t, uint16(10), stats.teamStats.ShotsOnGoalHome)
	assert.Equal(t, uint16(2), stats.teamStats.ShotsOnGoalAway)

	assert.Equal(t, uint16(6), stats.teamStats.FoulsHome)
	assert.Equal(t, uint16(9), stats.teamStats.FoulsAway)

	assert.Equal(t, uint8(0), stats.teamStats.RedCardsHome)
	assert.Equal(t, uint8(1), stats.teamStats.RedCardsAway)

	assert.Equal(t, uint8(0), stats.teamStats.YellowCardsHome)
	assert.Equal(t, uint8(4), stats.teamStats.YellowCardsAway)

	assert.Equal(t, uint16(18), stats.goalieStatsHome[26937].PassAttempts)
	assert.Equal(t, uint16(18), stats.goalieStatsHome[26937].CompletePasses)
	assert.Equal(t, uint8(1), stats.goalieStatsHome[26937].Saves)
	assert.Equal(t, float32(6.5), stats.goalieStatsHome[26937].Rating)
}
