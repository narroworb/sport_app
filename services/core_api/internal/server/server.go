package server

import (
	"net/http"

	"github.com/narroworb/core_api/internal/handlers"
	"github.com/narroworb/core_api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type ServerRepo struct {
	handlers handlers.HandlerRepo
	port     string
}

func NewServerRepo(db handlers.DatabaseInterface, cacheDB handlers.CacheInterface, port string) *ServerRepo {
	return &ServerRepo{
		handlers: *handlers.NewHandlerRepo(db, cacheDB),
		port:     port,
	}
}

func (s *ServerRepo) Run() error {
	r := chi.NewRouter()

	r.Get("/api/search", s.handlers.Search)

	r.Mount("/api/player", s.playerRoutes())
	r.Mount("/api/team", s.teamRoutes())
	r.Mount("/api/tournament", s.tournamentRoutes())
	r.Mount("/api/manager", s.managerRoutes())
	r.Mount("/api/fixture", s.fixtureRoutes())

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTAuth)

		r.Get("/player/favorites", s.handlers.GetFavouritePlayers)
		r.Get("/team/favorites", s.handlers.GetFavouriteTeams)
		r.Get("/manager/favorites", s.handlers.GetFavouriteManagers)
	})

	return http.ListenAndServe(":"+s.port, r)
}

func (s *ServerRepo) playerRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/details", s.handlers.GetPlayerDetails)

		r.Get("/stats", s.handlers.GetPlayerStats)
		r.Get("/fixtures", s.handlers.GetPlayerFixtures)
		r.Get("/teams", s.handlers.GetPlayerTeams)
	})

	return r
}

func (s *ServerRepo) teamRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/details", s.handlers.GetTeamDetails)

		r.Get("/stats", s.handlers.GetTeamStats)
		r.Get("/next_game", s.handlers.GetTeamNextGame)
		r.Get("/standings", s.handlers.GetTeamStandings)
		r.Get("/players", s.handlers.GetTeamPlayers)
		r.Get("/fixtures", s.handlers.GetTeamLastGames)
		r.Get("/manager", s.handlers.GetTeamManager)
		r.Get("/players_stats", s.handlers.GetTeamPlayersWithStats)
	})

	return r
}

func (s *ServerRepo) tournamentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/details", s.handlers.GetTournamentDetails)
		r.Get("/table", s.handlers.GetTournamentTable)
		r.Get("/stats/teams", s.handlers.GetTournamentTeamsStats)
		r.Get("/stats/players", s.handlers.GetTournamentPlayersStats)
		r.Get("/table/graph", s.handlers.GetTournamentTableGraph)
		r.Get("/fixtures", s.handlers.GetTournamentFixtures)
	})

	return r
}

func (s *ServerRepo) managerRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/details", s.handlers.GetManagerDetails)
		r.Get("/stats", s.handlers.GetManagerStats)
		r.Get("/teams", s.handlers.GetManagerTeams)
		r.Get("/fixtures", s.handlers.GetManagerFixtures)
	})

	return r
}

func (s *ServerRepo) fixtureRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/details", s.handlers.GetFixtureDetails)
		r.Get("/stats/players", s.handlers.GetFixturePlayersStats)
		r.Get("/stats/goalies", s.handlers.GetFixtureGoaliesStats)
		r.Get("/stats/teams", s.handlers.GetFixtureTeamsStats)
	})

	r.Get("/", s.handlers.GetFixturesByDate)

	return r
}
