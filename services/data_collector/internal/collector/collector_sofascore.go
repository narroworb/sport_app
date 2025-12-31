package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/narroworb/data_collector/internal/models"
)

type DatabaseInterface interface {
	GetUnactualTournamentsAndTours(context.Context) ([]UnactualTournamentsAndTours, error)
	GetFootballTournamentID(ctx context.Context, name, season string) (uint32, error)
	GetFootballTeamID(ctx context.Context, name string) (uint32, error)
	InsertFootballManager(ctx context.Context, manager *models.Manager) (uint32, error)
	GetFootballManagerID(ctx context.Context, manager *models.Manager) (uint32, error)
	GetFootballMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error)
	GetFootballMatchStatus(ctx context.Context, matchID uint32) (string, error)
	UpdateFootballMatch(ctx context.Context, match *models.Match) error
	InsertFootballMatch(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error)
	GetFootballPlayerID(ctx context.Context, name string, dateOfBirth time.Time) (uint32, error)
	InsertFootballPlayer(ctx context.Context, player *models.Player) (uint32, error)
	IncrementYellowCardsManager(ctx context.Context, managerID uint32) error
	IncrementRedCardsManager(ctx context.Context, managerID uint32) error
	InsertFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error)
	GetFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error)
	InsertFootballPlayerMatchStats(ctx context.Context, stats models.PlayerStatsInMatch, matchID uint32) (uint32, error)
	GetFootballPlayerMatchStats(ctx context.Context, stats models.PlayerStatsInMatch, matchID uint32) (uint32, error)
	InsertFootballGoalieMatchStats(ctx context.Context, stats models.GoalieStatsInMatch, matchID uint32) (uint32, error)
	GetFootballGoalieMatchStats(ctx context.Context, stats models.GoalieStatsInMatch, matchID uint32) (uint32, error)
	InsertFootballGoalieMatchStatsBatch(ctx context.Context, statsBatch map[uint32]*models.GoalieStatsInMatch, matchID uint32) error
	InsertFootballPlayerMatchStatsBatch(ctx context.Context, statsBatch map[uint32]*models.PlayerStatsInMatch, matchID uint32) error
	GetFootballNotPlayedMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error)
	GetFootballPlayedMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error)
	GetCountPlayersStatsByMatchID(ctx context.Context, matchID uint32) (uint64, error)
	GetFootballTeamTournamentPerformanceID(ctx context.Context, tournamentID uint32, teamID uint32) (uint32, error)
	InsertFootballTeamTournamentPerformanceID(ctx context.Context, rowTable *models.TableRow, tournamentID uint32) error
	UpdateFootballTeamTournamentPerformanceID(ctx context.Context, rowTable *models.TableRow, tournamentID uint32, statID uint32) error
	GetUpcomingTours(ctx context.Context) ([]UnactualTournamentsAndTours, error)
}

type SofaApi interface {
	FetchBodyConc(ctx context.Context, url string) (string, error)
	FindManagersOfMatch(ctx context.Context, url string) (homeID, awayID string)
}

type Producer interface {
	Send(ctx context.Context, key string, msg any) error
}

type Updater struct {
	db       DatabaseInterface
	api      SofaApi
	producer Producer
}

type ApiClient struct {
	browserCtx context.Context
}

func newApiClient() *ApiClient {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ExecPath("/usr/bin/chromium"),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/117.0.0.0 Safari/537.36"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, _ := chromedp.NewContext(allocCtx)

	return &ApiClient{browserCtx: browserCtx}
}

func NewUpdater(db DatabaseInterface, producer Producer) *Updater {
	return &Updater{
		db:       db,
		api:      newApiClient(),
		producer: producer,
	}
}

type tournamentForApi struct {
	baseURL    string
	leagueName string
	season     string
	tours      []uint16
}

func (u *Updater) StartUpdate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataToCollect, err := u.db.GetUnactualTournamentsAndTours(ctx)
	if err != nil {
		log.Printf("error in get data from clickhouse: %v\n", err)
		return
	}

	if len(dataToCollect) == 0 {
		log.Println("there is nothing to update wits stats, end of operation!")
		return
	}

	dataToCollectForApi := make([]tournamentForApi, 0, len(dataToCollect))

	for _, d := range dataToCollect {
		dataToCollectForApi = append(dataToCollectForApi,
			tournamentForApi{
				baseURL:    tournaments[d.LeagueName] + seasonsIDs[d.Season][d.LeagueName],
				leagueName: d.LeagueName,
				season:     d.Season,
				tours:      d.Tours,
			},
		)
	}

	const maxWorkers = 5
	tasks := make(chan func(context.Context), 50)
	wg := &sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		go func() {
			workerCtx, _ := chromedp.NewContext(u.api.(*ApiClient).browserCtx)

			// fmt.Println(u.api.FindManagersOfMatch(workerCtx, "https://www.sofascore.com/football/match/sc-heerenveen-az-alkmaar/ajbsojb#id:14053614"))

			for task := range tasks {
				// ctx, cancel := context.WithTimeout(workerCtx, 40*time.Second)
				task(workerCtx)
				// cancel()
				wg.Done()
			}
		}()
	}

	for _, tournamentToUpdateInfo := range dataToCollectForApi {
		wg.Add(1)

		tournamentToUpdateInfo := tournamentToUpdateInfo
		tasks <- func(ctx context.Context) {
			url := tournamentToUpdateInfo.baseURL + "/standings/total"
			body, err := u.api.FetchBodyConc(ctx, url)
			if err != nil {
				log.Printf("Ошибка при запросе %s: %v\n", url, err)
				return
			}

			var r StandingsResponse
			if len(body) < 5 {
				return
			}
			body = strings.TrimSuffix(body[5:], `</pre><div class="json-formatter-container"></div>`)
			if err := json.Unmarshal([]byte(body), &r); err != nil {
				log.Printf("Ошибка JSON: %v\n", err)
				return
			}

			season := models.Season{
				Year:  tournamentToUpdateInfo.season,
				Table: make(map[uint8]models.TableRow, len(r.Standings[0].Rows)),
				Teams: make(map[string]*models.Team, len(r.Standings[0].Rows)),
				Tours: make(map[uint16][]models.Match, (len(r.Standings[0].Rows)-1)*2),
			}

			log.Printf("Start to fetch season %s %s\n", tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season)

			var seasonID uint32
			if seasonID, err = u.db.GetFootballTournamentID(ctx, tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season); err != nil || seasonID == 0 {
				log.Printf("Season %s %s not found in db. err: %v\n", tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season, err)
				return
			} else {
				log.Printf("Season %s %s is currently in DB with ID %d\n", tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season, seasonID)
			}

			for _, row := range r.Standings[0].Rows {
				// log.Printf("Start to fetch team %s\n", row.Team.Name)

				teamID, err := u.db.GetFootballTeamID(ctx, row.Team.Name)
				if err != nil {
					log.Printf("error in get team %s from db: %v\n", row.Team.Name, err)
					return
				}

				team := &models.Team{
					Name: row.Team.Name,
					ID:   teamID,
				}

				row := models.TableRow{
					Team:          team,
					Points:        row.Points,
					Pos:           row.Pos,
					Matches:       row.Matches,
					Wins:          row.Wins,
					Losses:        row.Losses,
					Draws:         row.Draws,
					ScoresFor:     row.ScoresFor,
					ScoresAgainst: row.ScoresAgainst,
				}
				season.Table[row.Pos] = row
				season.Teams[row.Team.Name] = team

				if id, err := u.db.GetFootballTeamTournamentPerformanceID(ctx, seasonID, row.Team.ID); err != nil {
					if err := u.db.InsertFootballTeamTournamentPerformanceID(ctx, &row, seasonID); err != nil {
						fmt.Printf("error in insert new team(id=%d) tournament(id=%d) performance: %v", row.Team.ID, seasonID, err)
					}
				} else {
					if err := u.db.UpdateFootballTeamTournamentPerformanceID(ctx, &row, seasonID, id); err != nil {
						fmt.Printf("error in update team(id=%d) tournament(id=%d) performance: %v", row.Team.ID, seasonID, err)
					}
				}
			}

			fmt.Println("Обработка туров", tournamentToUpdateInfo.leagueName, season.Year)
			for _, tour := range tournamentToUpdateInfo.tours {
				log.Printf("парсинг %d тура лиги %s %s\n", tour, tournamentToUpdateInfo.leagueName, season.Year)
				matches, err := u.fetchMatches(ctx, tournamentToUpdateInfo.baseURL, tour, season.Teams)
				if err != nil {
					log.Printf("Ошибка при обработке матчей %d тура лиги %s %s: %v\n", tour, tournamentToUpdateInfo.leagueName, season.Year, err)
					continue
				}
				log.Printf("данные %d тура лиги %s %s получены.\n", tour, tournamentToUpdateInfo.leagueName, season.Year)

				for indMatch, match := range matches {
					log.Printf("Матч %s - %s от %s в обработке. \n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
					if id, err := u.db.GetFootballPlayedMatchID(ctx, &(matches[indMatch]), seasonID); err == nil && id != 0 {
						if cnt, err := u.db.GetCountPlayersStatsByMatchID(ctx, id); err == nil && cnt >= 20 {
							log.Printf("Матч %s - %s от %s уже в базе и полностью обработан.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
							continue
						}
						log.Printf("Матч %s - %s от %s уже в базе и требует доп обработки.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
						matches[indMatch].IDAppDB = id
					} else if id, err := u.db.GetFootballNotPlayedMatchID(ctx, &(matches[indMatch]), seasonID); err == nil && id != 0 {
						log.Printf("Матч %s - %s от %s уже в базе со статусом 'Не сыгран'.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
						matches[indMatch].IDAppDB = id
						err := u.db.UpdateFootballMatch(ctx, &(matches[indMatch]))
						if err != nil {
							log.Printf("error in updating match %+v: %v\n", matches[indMatch], err)
							continue
						}
					} else if id, err := u.db.GetFootballMatchID(ctx, &(matches[indMatch]), seasonID); err == nil && id != 0 {
						log.Printf("Матч %s - %s от %s уже в базе.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
						matches[indMatch].IDAppDB = id
						status, err := u.db.GetFootballMatchStatus(ctx, id)
						if err != nil {
							log.Printf("error in find match status: %v\n", err)
							continue
						}
						if status != matches[indMatch].Status {
							err := u.db.UpdateFootballMatch(ctx, &(matches[indMatch]))
							if err != nil {
								log.Printf("error in updating match %+v: %v\n", matches[indMatch], err)
								continue
							}
						} else {
							log.Printf("Матч %s - %s от %s был обработан ранее.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
							continue
						}
					} else {
						id, err := u.db.InsertFootballMatch(ctx, &(matches[indMatch]), seasonID)
						matches[indMatch].IDAppDB = id
						if err != nil {
							log.Printf("error in insert match: %v\n", err)
							continue
						}
					}

					s, err := u.fetchAllStatisticsFromMatches(ctx, fmt.Sprint(match.ID))
					if err != nil {
						log.Println("Ошибка в обработке статистики матча: ", err)
					}
					u.addStatsToDB(ctx, s, matches[indMatch].IDAppDB)
					log.Printf("Матч %s - %s от %s обработан.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
				}

				log.Printf("Тур %d лиги %s %s обработан полностью\n", tour, tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season)
			}
		}
	}
	close(tasks)
	wg.Wait()
	chromedp.Cancel(u.api.(*ApiClient).browserCtx)
	log.Println("update done succesfully")
}

func (u *Updater) StartUpdateWithoutStatistics() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataToCollect, err := u.db.GetUpcomingTours(ctx)
	if err != nil {
		log.Printf("error in get data from clickhouse: %v\n", err)
		return
	}

	dataToCollectForApi := make([]tournamentForApi, 0, len(dataToCollect))

	for _, d := range dataToCollect {
		dataToCollectForApi = append(dataToCollectForApi,
			tournamentForApi{
				baseURL:    tournaments[d.LeagueName] + seasonsIDs[d.Season][d.LeagueName],
				leagueName: d.LeagueName,
				season:     d.Season,
				tours:      d.Tours,
			},
		)
	}

	const maxWorkers = 5
	tasks := make(chan func(context.Context), 50)
	wg := &sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		go func() {
			workerCtx, _ := chromedp.NewContext(u.api.(*ApiClient).browserCtx)

			for task := range tasks {
				task(workerCtx)
				wg.Done()
			}
		}()
	}

	for _, tournamentToUpdateInfo := range dataToCollectForApi {
		wg.Add(1)

		tournamentToUpdateInfo := tournamentToUpdateInfo
		tasks <- func(ctx context.Context) {
			url := tournamentToUpdateInfo.baseURL + "/standings/total"
			body, err := u.api.FetchBodyConc(ctx, url)
			if err != nil {
				log.Printf("Ошибка при запросе %s: %v\n", url, err)
				return
			}

			var r StandingsResponse
			if len(body) < 5 {
				return
			}
			body = strings.TrimSuffix(body[5:], `</pre><div class="json-formatter-container"></div>`)
			if err := json.Unmarshal([]byte(body), &r); err != nil {
				log.Printf("Ошибка JSON: %v\n", err)
				return
			}

			season := models.Season{
				Year:  tournamentToUpdateInfo.season,
				Table: make(map[uint8]models.TableRow, len(r.Standings[0].Rows)),
				Teams: make(map[string]*models.Team, len(r.Standings[0].Rows)),
				Tours: make(map[uint16][]models.Match, (len(r.Standings[0].Rows)-1)*2),
			}

			log.Printf("Start to fetch season %s %s\n", tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season)

			var seasonID uint32
			if seasonID, err = u.db.GetFootballTournamentID(ctx, tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season); err != nil || seasonID == 0 {
				log.Printf("Season %s %s not found in db. err: %v\n", tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season, err)
				return
			} else {
				log.Printf("Season %s %s is currently in DB with ID %d\n", tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season, seasonID)
			}

			for _, row := range r.Standings[0].Rows {
				// log.Printf("Start to fetch team %s\n", row.Team.Name)

				teamID, err := u.db.GetFootballTeamID(ctx, row.Team.Name)
				if err != nil {
					log.Printf("error in get team %s from db: %v\n", row.Team.Name, err)
					return
				}

				team := &models.Team{
					Name: row.Team.Name,
					ID:   teamID,
				}

				row := models.TableRow{
					Team:          team,
					Points:        row.Points,
					Pos:           row.Pos,
					Matches:       row.Matches,
					Wins:          row.Wins,
					Losses:        row.Losses,
					Draws:         row.Draws,
					ScoresFor:     row.ScoresFor,
					ScoresAgainst: row.ScoresAgainst,
				}
				season.Table[row.Pos] = row
				season.Teams[row.Team.Name] = team

				if id, err := u.db.GetFootballTeamTournamentPerformanceID(ctx, seasonID, row.Team.ID); err != nil {
					if err := u.db.InsertFootballTeamTournamentPerformanceID(ctx, &row, seasonID); err != nil {
						fmt.Printf("error in insert new team(id=%d) tournament(id=%d) performance: %v", row.Team.ID, seasonID, err)
					}
				} else {
					if err := u.db.UpdateFootballTeamTournamentPerformanceID(ctx, &row, seasonID, id); err != nil {
						fmt.Printf("error in update team(id=%d) tournament(id=%d) performance: %v", row.Team.ID, seasonID, err)
					}
				}
			}

			fmt.Println("Обработка туров", tournamentToUpdateInfo.leagueName, season.Year)
			for _, tour := range tournamentToUpdateInfo.tours {
				log.Printf("парсинг %d тура лиги %s %s\n", tour, tournamentToUpdateInfo.leagueName, season.Year)
				matches, err := u.fetchMatches(ctx, tournamentToUpdateInfo.baseURL, tour, season.Teams)
				if err != nil {
					log.Printf("Ошибка при обработке матчей %d тура лиги %s %s: %v\n", tour, tournamentToUpdateInfo.leagueName, season.Year, err)
					continue
				}
				log.Printf("данные %d тура лиги %s %s получены.\n", tour, tournamentToUpdateInfo.leagueName, season.Year)

				for indMatch, match := range matches {
					log.Printf("Матч %s - %s от %s в обработке. \n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
					if id, err := u.db.GetFootballPlayedMatchID(ctx, &(matches[indMatch]), seasonID); err == nil && id != 0 {
						log.Printf("Матч %s - %s от %s уже в базе и полностью обработан.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
						continue
					} else if id, err := u.db.GetFootballNotPlayedMatchID(ctx, &(matches[indMatch]), seasonID); err == nil && id != 0 {
						log.Printf("Матч %s - %s от %s уже в базе со статусом 'Не сыгран'.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
						matches[indMatch].IDAppDB = id
						err := u.db.UpdateFootballMatch(ctx, &(matches[indMatch]))
						if err != nil {
							log.Printf("error in updating match %+v: %v\n", matches[indMatch], err)
							continue
						}
					} else if id, err := u.db.GetFootballMatchID(ctx, &(matches[indMatch]), seasonID); err == nil && id != 0 {
						log.Printf("Матч %s - %s от %s уже в базе.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
						matches[indMatch].IDAppDB = id
						status, err := u.db.GetFootballMatchStatus(ctx, id)
						if err != nil {
							log.Printf("error in find match status: %v\n", err)
							continue
						}
						if status != matches[indMatch].Status {
							err := u.db.UpdateFootballMatch(ctx, &(matches[indMatch]))
							if err != nil {
								log.Printf("error in updating match %+v: %v\n", matches[indMatch], err)
								continue
							}
						} else {
							log.Printf("Матч %s - %s от %s был обработан ранее.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
							continue
						}
					} else {
						id, err := u.db.InsertFootballMatch(ctx, &(matches[indMatch]), seasonID)
						matches[indMatch].IDAppDB = id
						if err != nil {
							log.Printf("error in insert match: %v\n", err)
							continue
						}
					}

					log.Printf("Матч %s - %s от %s обновлен.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Date)
				}

				log.Printf("Тур %d лиги %s %s обработан полностью\n", tour, tournamentToUpdateInfo.leagueName, tournamentToUpdateInfo.season)
			}
		}
	}
	close(tasks)
	wg.Wait()
	chromedp.Cancel(u.api.(*ApiClient).browserCtx)
	log.Println("update without statistics done succesfully")
}

func (a *ApiClient) FetchBodyConc(ctx context.Context, url string) (string, error) {
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(15*time.Second),
		chromedp.InnerHTML("body", &body),
	)
	if len(body) < 5 || strings.TrimSuffix(body[5:], `</pre><div class="json-formatter-container"></div>`) == `{"error": {"code": 403, "reason": "challenge" }}` {
		log.Printf("Error in fetching body of url: %s\n", url)
		fmt.Println(body)
		return "", fmt.Errorf(`get error with code 403 by reason "challenge"`)
	}
	return body, err
}

type StandingsResponse struct {
	Standings []struct {
		Rows []struct {
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
			Points        uint16 `json:"points"`
			Pos           uint8  `json:"position"`
			Matches       uint16 `json:"matches"`
			Wins          uint16 `json:"wins"`
			Losses        uint16 `json:"losses"`
			Draws         uint16 `json:"draws"`
			ScoresFor     uint16 `json:"scoresFor"`
			ScoresAgainst uint16 `json:"scoresAgainst"`
		} `json:"rows"`
	} `json:"standings"`
}

type MatchesResponse struct {
	Events []struct {
		ID        int    `json:"id"`
		CustomID  string `json:"customID"`
		Slug      string `json:"slug"`
		StartTime int64  `json:"startTimestamp"`
		HomeTeam  struct {
			Name string `json:"name"`
		} `json:"homeTeam"`
		AwayTeam struct {
			Name string `json:"name"`
		} `json:"awayTeam"`
		HomeScore struct {
			Display int `json:"display"`
		} `json:"homeScore"`
		AwayScore struct {
			Display int `json:"display"`
		} `json:"awayScore"`
		Status struct {
			Description string `json:"description"`
		} `json:"status"`
	} `json:"events"`
}

func (u *Updater) fetchMatches(ctx context.Context, url_base string, roundID uint16, allTeams map[string]*models.Team) ([]models.Match, error) {
	url := fmt.Sprintf("%s/events/round/%d", url_base, roundID)
	body, err := u.api.FetchBodyConc(ctx, url)
	if err != nil {
		return nil, err
	}

	var result MatchesResponse
	if len(body) < 5 {
		return nil, fmt.Errorf("too short body of response")
	}
	body = strings.TrimSuffix(body[5:], `</pre><div class="json-formatter-container"></div>`)
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}

	matches := make([]models.Match, 0, len(result.Events))
	for _, e := range result.Events {
		url := fmt.Sprintf("https://www.sofascore.com/football/match/%s/%s#id:%d", e.Slug, e.CustomID, e.ID)
		homeManagerID, awayManagerID := u.api.FindManagersOfMatch(ctx, url)
		if homeManagerID == "" && awayManagerID == "" {
			log.Printf("first attempt to find managers of match %s - %s by url=%s is failed\n", e.HomeTeam.Name, e.AwayTeam.Name, url)
			for attempt := 2; attempt <= 3; attempt++ {
				time.Sleep(time.Second * 5)
				homeManagerID, awayManagerID = u.api.FindManagersOfMatch(ctx, url)
				if homeManagerID != "" || awayManagerID != "" {
					log.Printf("attempt %d to find managers of match %s - %s by url=%s is success\n", attempt, e.HomeTeam.Name, e.AwayTeam.Name, url)
					break
				}
				log.Printf("attempt %d to find managers of match %s - %s by url=%s is failed\n", attempt, e.HomeTeam.Name, e.AwayTeam.Name, url)
			}
		}
		homeManager, err := u.fetchManager(ctx, homeManagerID)
		if err != nil || (homeManager.FirstName == "" && homeManager.LastName == "") {
			log.Printf("Home manager from match %s - %s by url %s not found.\n", e.HomeTeam.Name, e.AwayTeam.Name, url)
			homeManager = models.Manager{FirstName: "Not", LastName: "Find"}
		}

		awayManager, err := u.fetchManager(ctx, awayManagerID)
		if err != nil || (awayManager.FirstName == "" && awayManager.LastName == "") {
			log.Printf("Away manager from match %s - %s by url %s not found.\n", e.HomeTeam.Name, e.AwayTeam.Name, url)
			awayManager = models.Manager{FirstName: "Not", LastName: "Find"}
		}

		if id, err := u.db.GetFootballManagerID(ctx, &homeManager); err != nil || id == 0 {
			_, err := u.db.InsertFootballManager(ctx, &homeManager)
			if err != nil {
				log.Printf("error in insert manager: %v\n", err)
			}
		}

		if id, err := u.db.GetFootballManagerID(ctx, &awayManager); err != nil || id == 0 {
			_, err := u.db.InsertFootballManager(ctx, &awayManager)
			if err != nil {
				log.Printf("error in insert manager: %v\n", err)
			}
		}

		matches = append(matches, models.Match{
			ID:          e.ID,
			Date:        time.Unix(e.StartTime, 0).Add(time.Hour * 3),
			HomeTeam:    allTeams[e.HomeTeam.Name],
			AwayTeam:    allTeams[e.AwayTeam.Name],
			HomeGoals:   uint16(e.HomeScore.Display),
			AwayGoals:   uint16(e.AwayScore.Display),
			Round:       uint16(roundID),
			HomeManager: &homeManager,
			AwayManager: &awayManager,
			Status:      e.Status.Description,
		})
	}

	return matches, nil
}

func (a *ApiClient) FindManagersOfMatch(ctx context.Context, url string) (homeID, awayID string) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// ctx, cancel = chromedp.NewContext(ctxTimeout)
	// defer cancel()

	var htmls []string

	js := `Array.from(document.querySelectorAll('img[src*="manager"]')).map(el => el.outerHTML)`

	err := chromedp.Run(ctxTimeout,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Evaluate(js, &htmls),
	)
	if err != nil {
		log.Printf("error in fetch managers from match %s: %v\n", url, err)
		return
	}

	if len(htmls) != 2 {
		return
	}
	start := strings.Index(htmls[0], "/manager/") + 9
	end := strings.Index(htmls[0], "/image")
	homeID = htmls[0][start:end]
	start = strings.Index(htmls[1], "/manager/") + 9
	end = strings.Index(htmls[1], "/image")
	awayID = htmls[1][start:end]

	return
}

type ManagerResponse struct {
	Manager struct {
		Name   string `json:"name"`
		Nation struct {
			Name string `json:"name"`
		} `json:"country"`
	} `json:"manager"`
}

func (u *Updater) fetchManager(ctx context.Context, managerID string) (manager models.Manager, err error) {
	url := fmt.Sprintf("https://api.sofascore.com/api/v1/manager/%s", managerID)
	body, err := u.api.FetchBodyConc(ctx, url)
	if err != nil {
		return
	}
	var result ManagerResponse
	if len(body) < 5 {
		return models.Manager{}, fmt.Errorf("too short body of response")
	}
	body = strings.TrimSuffix(body[5:], `</pre><div class="json-formatter-container"></div>`)

	if err = json.Unmarshal([]byte(body), &result); err != nil {
		return
	}

	// fmt.Printf("%+v", result)

	manager = models.Manager{
		FirstName: strings.Split(result.Manager.Name, " ")[0],
		LastName:  strings.Join(strings.Split(result.Manager.Name, " ")[1:], " "),
		Nation:    result.Manager.Nation.Name,
	}

	return
}

type StatsFromMatch struct {
	teamStats        models.TeamMatchStats
	playersStatsHome map[uint32]*models.PlayerStatsInMatch
	goalieStatsHome  map[uint32]*models.GoalieStatsInMatch
	playersStatsAway map[uint32]*models.PlayerStatsInMatch
	goalieStatsAway  map[uint32]*models.GoalieStatsInMatch
}

func (u *Updater) fetchAllStatisticsFromMatches(ctx context.Context, matchID string) (StatsFromMatch, error) {
	var res StatsFromMatch

	urlPlayers := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%s/lineups", matchID)
	urlTeams := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%s/statistics", matchID)
	urlIncidents := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%s/incidents", matchID)

	teamStatsBody, err := u.api.FetchBodyConc(ctx, urlTeams)
	if err != nil {
		return StatsFromMatch{}, err
	}

	if len(teamStatsBody) < 5 {
		return StatsFromMatch{}, fmt.Errorf("too short body of response")
	}
	teamStatsBody = strings.TrimSuffix(teamStatsBody[5:], `</pre><div class="json-formatter-container"></div>`)

	var teamStatsMatchResp TeamStatsResponse

	if err := json.Unmarshal([]byte(teamStatsBody), &teamStatsMatchResp); err != nil {
		log.Println("Ошибка в обработке json по url: ", urlTeams)
		return StatsFromMatch{}, err
	}

	res.teamStats = TSResponseToTSStruct(teamStatsMatchResp)

	playersStatsBody, err := u.api.FetchBodyConc(ctx, urlPlayers)
	if err != nil {
		return StatsFromMatch{}, err
	}

	if len(playersStatsBody) < 5 {
		return StatsFromMatch{}, fmt.Errorf("too short body of response")
	}
	playersStatsBody = strings.TrimSuffix(playersStatsBody[5:], `</pre><div class="json-formatter-container"></div>`)

	var playerStatsMatchResp PlayerStatsResponse

	if err := json.Unmarshal([]byte(playersStatsBody), &playerStatsMatchResp); err != nil {
		return StatsFromMatch{}, err
	}

	_, _ = u.pSResponseToPSStructs(ctx, playerStatsMatchResp, &res)

	incidentsBody, err := u.api.FetchBodyConc(ctx, urlIncidents)
	if err != nil {
		return StatsFromMatch{}, err
	}

	if len(incidentsBody) < 5 {
		return StatsFromMatch{}, fmt.Errorf("too short body of response")
	}
	incidentsBody = strings.TrimSuffix(incidentsBody[5:], `</pre><div class="json-formatter-container"></div>`)

	var incidentsResp IncidentsResponse

	if err := json.Unmarshal([]byte(incidentsBody), &incidentsResp); err != nil {
		return StatsFromMatch{}, err
	}

	u.iResponseToPS(ctx, incidentsResp, &res, urlIncidents)

	return res, err
}

type TeamStatsResponse struct {
	Statistics []struct {
		Period string `json:"period"`
		Groups []struct {
			Stats []struct {
				Name string  `json:"name"`
				Home float32 `json:"homeValue"`
				Away float32 `json:"awayValue"`
			} `json:"statisticsItems"`
		} `json:"groups"`
	} `json:"statistics"`
}

type PlayerStatsResponse struct {
	HomeTeam struct {
		Players []struct {
			Player struct {
				Name     string `json:"name"`
				ID       int    `json:"id"`
				Position string `json:"position"`
				Height   uint16 `json:"height"`
				Country  struct {
					Name string `json:"name"`
				} `json:"Country"`
				DateOfBirthTimeStamp int64 `json:"dateOfBirthTimestamp"`
			} `json:"player"`
			BenchPlayer bool `json:"substitute"`
			Statistics  struct {
				Rating        float32 `json:"rating"`
				MinutesPlayed uint16  `json:"minutesPlayed"`
				Goals         uint8   `json:"goals"`
				Assists       uint8   `json:"goalAssist"`
				BlockedShots  uint8   `json:"outfielderBlock"`
				Interceptions uint8   `json:"interceptionWon"`
				TotalTackles  uint8   `json:"totalTackle"`

				DribbledPast uint8 `json:"challengeLost"`

				DuelsWon         uint8  `json:"duelWon"`
				DuelsLost        uint8  `json:"duelLost"`
				Fouls            uint8  `json:"fouls"`
				WasFouled        uint8  `json:"wasFouled"`
				PassAttempts     uint16 `json:"totalPass"`
				CompletePasses   uint16 `json:"accuratePass"`
				KeyPasses        uint8  `json:"keyPass"`
				ShotsOffTarget   uint8  `json:"shotOffTarget"`
				ShotsOnTarget    uint8  `json:"onTargetScoringAttempt"`
				ShotsBlocked     uint8  `json:"blockedScoringAttempt"`
				DribbleAttempts  uint8  `json:"totalContest"`
				CompleteDribbles uint8  `json:"wonContest"`
				PenaltyScored    uint8  `json:"penaltyWon"`
				PenaltyMissed    uint8  `json:"penaltyMiss"`
				YellowCards      uint8  `json:""`
				RedCards         uint8  `json:""`

				GoalsConceded   uint8 `json:""`
				Saves           uint8 `json:"saves"`
				PenaltySaved    uint8 `json:"penaltySave"`
				PenaltyConceded uint8 `json:""`

				Captain bool `json:"captain"`
			} `json:"statistics"`
		} `json:"players"`
	} `json:"home"`

	AwayTeam struct {
		Players []struct {
			Player struct {
				Name     string `json:"name"`
				ID       int    `json:"id"`
				Position string `json:"position"`
				Height   uint16 `json:"height"`
				Country  struct {
					Name string `json:"name"`
				} `json:"Country"`
				DateOfBirthTimeStamp int64 `json:"dateOfBirthTimestamp"`
			} `json:"player"`
			BenchPlayer bool `json:"substitute"`
			Statistics  struct {
				Rating        float32 `json:"rating"`
				MinutesPlayed uint16  `json:"minutesPlayed"`
				Goals         uint8   `json:"goals"`
				Assists       uint8   `json:"goalAssist"`
				BlockedShots  uint8   `json:"outfielderBlock"`
				Interceptions uint8   `json:"interceptionWon"`
				TotalTackles  uint8   `json:"totalTackle"`

				DribbledPast uint8 `json:"challengeLost"`

				DuelsWon         uint8  `json:"duelWon"`
				DuelsLost        uint8  `json:"duelLost"`
				Fouls            uint8  `json:"fouls"`
				WasFouled        uint8  `json:"wasFouled"`
				PassAttempts     uint16 `json:"totalPass"`
				CompletePasses   uint16 `json:"accuratePass"`
				KeyPasses        uint8  `json:"keyPass"`
				ShotsOffTarget   uint8  `json:"shotOffTarget"`
				ShotsOnTarget    uint8  `json:"onTargetScoringAttempt"`
				ShotsBlocked     uint8  `json:"blockedScoringAttempt"`
				DribbleAttempts  uint8  `json:"totalContest"`
				CompleteDribbles uint8  `json:"wonContest"`
				PenaltyScored    uint8  `json:"penaltyWon"`
				PenaltyMissed    uint8  `json:"penaltyMiss"`
				YellowCards      uint8  `json:""`
				RedCards         uint8  `json:""`

				GoalsConceded   uint8 `json:""`
				Saves           uint8 `json:"saves"`
				PenaltySaved    uint8 `json:"penaltySave"`
				PenaltyConceded uint8 `json:""`

				Captain bool `json:"captain"`
			} `json:"statistics"`
		} `json:"players"`
	} `json:"away"`
}

type IncidentsResponse struct {
	Incidents []struct {
		Player struct {
			Name                 string `json:"name"`
			Position             string `json:"position"`
			Height               uint16 `json:"height"`
			DateOfBirthTimeStamp int64  `json:"dateOfBirthTimestamp"`
		} `json:"player"`
		Manager struct {
			Name string `json:"name"`
		} `json:"manager"`
		FootballPassingNetworkAction []struct {
			EventType  string `json:"eventType"`
			Goalkeeper struct {
				Name                 string `json:"name"`
				Position             string `json:"position"`
				Height               uint16 `json:"height"`
				DateOfBirthTimeStamp int64  `json:"dateOfBirthTimestamp"`
			} `json:"goalkeeper"`
		} `json:"footballPassingNetworkAction"`
		IncidentType  string `json:"incidentType"`
		IncidentClass string `json:"incidentClass"`
		IsHome        bool   `json:"isHome"`
		Time          int16  `json:"time"`
	} `json:"incidents"`
}

func TSResponseToTSStruct(resp TeamStatsResponse) models.TeamMatchStats {
	var res models.TeamMatchStats
	for _, stats := range resp.Statistics {
		if stats.Period != "ALL" {
			continue
		}

		for _, group := range stats.Groups {
			for _, stat := range group.Stats {
				switch stat.Name {
				case "Shots on target":
					res.ShotsOnGoalHome = uint16(stat.Home)
					res.ShotsOnGoalAway = uint16(stat.Away)
				case "Total shots":
					res.TotalShotsHome = uint16(stat.Home)
					res.TotalShotsAway = uint16(stat.Away)
				case "Blocked shots":
					res.BlockedShotsHome = uint16(stat.Home)
					res.BlockedShotsAway = uint16(stat.Away)
				case "Fouls":
					res.FoulsHome = uint16(stat.Home)
					res.FoulsAway = uint16(stat.Away)
				case "Corner kicks":
					res.CornerKicksHome = uint16(stat.Home)
					res.CornerKicksAway = uint16(stat.Away)
				case "Ball possession":
					res.BallPossessionHome = uint8(stat.Home)
					res.BallPossessionAway = uint8(stat.Away)
				case "Yellow cards":
					res.YellowCardsHome = uint8(stat.Home)
					res.YellowCardsAway = uint8(stat.Away)
				case "Red cards":
					res.RedCardsHome = uint8(stat.Home)
					res.RedCardsAway = uint8(stat.Away)
				case "Passes":
					res.TotalPassesHome = uint16(stat.Home)
					res.TotalPassesAway = uint16(stat.Away)
				case "Accurate passes":
					res.CompletePassesHome = uint16(stat.Home)
					res.CompletePassesAway = uint16(stat.Away)
				case "Offsides":
					res.OffsidesHome = uint8(stat.Home)
					res.OffsidesAway = uint8(stat.Away)
				case "Shots inside box":
					res.ShotsInsideBoxHome = uint8(stat.Home)
					res.ShotsInsideBoxAway = uint8(stat.Away)
				}
			}
		}
		break
	}
	return res
}

func (u *Updater) pSResponseToPSStructs(ctx context.Context, resp PlayerStatsResponse, fullStats *StatsFromMatch) (RedCardsHome uint8, RedCardsAway uint8) {
	fullStats.playersStatsHome = make(map[uint32]*models.PlayerStatsInMatch, len(resp.HomeTeam.Players))
	fullStats.playersStatsAway = make(map[uint32]*models.PlayerStatsInMatch, len(resp.AwayTeam.Players))
	fullStats.goalieStatsHome = make(map[uint32]*models.GoalieStatsInMatch, 3)
	fullStats.goalieStatsAway = make(map[uint32]*models.GoalieStatsInMatch, 3)

	for _, p := range resp.HomeTeam.Players {
		id, err := u.db.GetFootballPlayerID(ctx, p.Player.Name, time.Unix(p.Player.DateOfBirthTimeStamp, 0))
		if err != nil || id == 0 {
			player, err := u.fetchPlayer(ctx, p.Player.ID)
			if err != nil {
				log.Printf("Ошибка в получении данных игрока %s: %v\n", p.Player.Name, err)
				continue
			}
			_, err = u.db.InsertFootballPlayer(ctx, player)
			if err != nil {
				log.Printf("error in inserting player %s: %v\n", p.Player.Name, err)
				continue
			}
			id = player.ID
		}
		if p.Player.Position == "G" {
			stats := &models.GoalieStatsInMatch{
				IDPlayer:        id,
				StartPlayer:     !p.BenchPlayer,
				Rating:          p.Statistics.Rating,
				MinutesPlayed:   p.Statistics.MinutesPlayed,
				Goals:           p.Statistics.Goals,
				Assists:         p.Statistics.Assists,
				GoalsConceded:   p.Statistics.GoalsConceded,
				Saves:           p.Statistics.Saves,
				PassAttempts:    p.Statistics.PassAttempts,
				CompletePasses:  p.Statistics.CompletePasses,
				KeyPasses:       p.Statistics.KeyPasses,
				PenaltySaved:    p.Statistics.PenaltySaved,
				PenaltyConceded: p.Statistics.PenaltyConceded,
				Fouls:           p.Statistics.Fouls,
				WasFouled:       p.Statistics.WasFouled,
				YellowCards:     p.Statistics.YellowCards,
				RedCards:        p.Statistics.RedCards,
				Captain:         p.Statistics.Captain,
				HomeTeamPlayer:  true,
			}
			RedCardsHome += stats.RedCards
			fullStats.goalieStatsHome[stats.IDPlayer] = stats
		} else {
			stats := &models.PlayerStatsInMatch{
				IDPlayer:      id,
				StartPlayer:   !p.BenchPlayer,
				Rating:        p.Statistics.Rating,
				MinutesPlayed: p.Statistics.MinutesPlayed,
				Goals:         p.Statistics.Goals,
				Assists:       p.Statistics.Assists,
				BlockedShots:  p.Statistics.BlockedShots,
				Interceptions: p.Statistics.Interceptions,
				TotalTackles:  p.Statistics.TotalTackles,
				DribbledPast:  p.Statistics.DribbledPast,
				Duels:         p.Statistics.DuelsLost + p.Statistics.DuelsWon,
				DuelsWon:      p.Statistics.DuelsWon,
				Fouls:         p.Statistics.Fouls,
				WasFouled:     p.Statistics.WasFouled,

				PassAttempts:     p.Statistics.PassAttempts,
				CompletePasses:   p.Statistics.CompletePasses,
				KeyPasses:        p.Statistics.KeyPasses,
				ShotsOnTarget:    p.Statistics.ShotsOnTarget,
				TotalShots:       p.Statistics.ShotsOffTarget + p.Statistics.ShotsOnTarget,
				DribbleAttempts:  p.Statistics.DribbleAttempts,
				CompleteDribbles: p.Statistics.CompleteDribbles,
				PenaltyScored:    p.Statistics.PenaltyScored,
				PenaltyMissed:    p.Statistics.PenaltyMissed,

				YellowCards: p.Statistics.YellowCards,
				RedCards:    p.Statistics.RedCards,

				Captain:        p.Statistics.Captain,
				HomeTeamPlayer: true,
			}
			RedCardsHome += stats.RedCards
			fullStats.playersStatsHome[stats.IDPlayer] = stats
		}
	}

	for _, p := range resp.AwayTeam.Players {
		id, err := u.db.GetFootballPlayerID(ctx, p.Player.Name, time.Unix(p.Player.DateOfBirthTimeStamp, 0))
		if err != nil || id == 0 {
			player, err := u.fetchPlayer(ctx, p.Player.ID)
			if err != nil {
				log.Printf("Ошибка в получении данных игрока %s: %v\n", p.Player.Name, err)
				continue
			}
			_, err = u.db.InsertFootballPlayer(ctx, player)
			if err != nil {
				log.Printf("error in inserting player %s: %v\n", p.Player.Name, err)
				continue
			}
			id = player.ID
		}
		if p.Player.Position == "G" {
			stats := &models.GoalieStatsInMatch{
				IDPlayer:        id,
				StartPlayer:     !p.BenchPlayer,
				Rating:          p.Statistics.Rating,
				MinutesPlayed:   p.Statistics.MinutesPlayed,
				Goals:           p.Statistics.Goals,
				Assists:         p.Statistics.Assists,
				GoalsConceded:   p.Statistics.GoalsConceded,
				Saves:           p.Statistics.Saves,
				PassAttempts:    p.Statistics.PassAttempts,
				CompletePasses:  p.Statistics.CompletePasses,
				KeyPasses:       p.Statistics.KeyPasses,
				PenaltySaved:    p.Statistics.PenaltySaved,
				PenaltyConceded: p.Statistics.PenaltyConceded,
				Fouls:           p.Statistics.Fouls,
				WasFouled:       p.Statistics.WasFouled,
				YellowCards:     p.Statistics.YellowCards,
				RedCards:        p.Statistics.RedCards,
				Captain:         p.Statistics.Captain,
				HomeTeamPlayer:  false,
			}
			RedCardsAway += stats.RedCards
			fullStats.goalieStatsAway[stats.IDPlayer] = stats
		} else {
			stats := &models.PlayerStatsInMatch{
				IDPlayer:      id,
				StartPlayer:   !p.BenchPlayer,
				Rating:        p.Statistics.Rating,
				MinutesPlayed: p.Statistics.MinutesPlayed,
				Goals:         p.Statistics.Goals,
				Assists:       p.Statistics.Assists,
				BlockedShots:  p.Statistics.BlockedShots,
				Interceptions: p.Statistics.Interceptions,
				TotalTackles:  p.Statistics.TotalTackles,
				DribbledPast:  p.Statistics.DribbledPast,
				Duels:         p.Statistics.DuelsLost + p.Statistics.DuelsWon,
				DuelsWon:      p.Statistics.DuelsWon,
				Fouls:         p.Statistics.Fouls,
				WasFouled:     p.Statistics.WasFouled,

				PassAttempts:     p.Statistics.PassAttempts,
				CompletePasses:   p.Statistics.CompletePasses,
				KeyPasses:        p.Statistics.KeyPasses,
				ShotsOnTarget:    p.Statistics.ShotsOnTarget,
				TotalShots:       p.Statistics.ShotsOffTarget + p.Statistics.ShotsOnTarget,
				DribbleAttempts:  p.Statistics.DribbleAttempts,
				CompleteDribbles: p.Statistics.CompleteDribbles,
				PenaltyScored:    p.Statistics.PenaltyScored,
				PenaltyMissed:    p.Statistics.PenaltyMissed,

				YellowCards: p.Statistics.YellowCards,
				RedCards:    p.Statistics.RedCards,

				Captain:        p.Statistics.Captain,
				HomeTeamPlayer: false,
			}
			RedCardsAway += stats.RedCards
			fullStats.playersStatsAway[stats.IDPlayer] = stats
		}
	}

	return
}

func (u *Updater) iResponseToPS(ctx context.Context, incidents IncidentsResponse, stats *StatsFromMatch, url string) {
	for _, incident := range incidents.Incidents {
		if incident.IncidentType == "goal" {
			if incident.IncidentClass == "regular" {
				if incident.IsHome {
					if len(incident.FootballPassingNetworkAction) == 0 {
						for _, g := range stats.goalieStatsAway {
							if incident.Time < 0 {
								continue
							}
							if (g.StartPlayer && g.MinutesPlayed <= uint16(incident.Time)) || (!g.StartPlayer && g.MinutesPlayed > 0) {
								g.GoalsConceded++
							}
							break
						}
						continue
					}
					keeper := incident.FootballPassingNetworkAction[len(incident.FootballPassingNetworkAction)-1].Goalkeeper
					// id := db.GetPlayerByName(strings.TrimSpace(incident.FootballPassingNetworkAction[len(incident.FootballPassingNetworkAction)-1].Goalkeeper.Name)).ID
					id, err := u.db.GetFootballPlayerID(ctx, strings.TrimSpace(keeper.Name), time.Unix(keeper.DateOfBirthTimeStamp, 0))
					if err != nil {
						log.Printf("error in search in incidents home goal: %v, url: %s, incident: %+v\n", err, url, incident)
						continue
					}
					if stats.goalieStatsAway[id] == nil {
						log.Printf("error in increment goals conceded in incidents home goal: %v, url: %s, incident: %+v\n", err, url, incident)
						continue
					}
					stats.goalieStatsAway[id].GoalsConceded++
				} else {
					if len(incident.FootballPassingNetworkAction) == 0 {
						for _, g := range stats.goalieStatsHome {
							if incident.Time < 0 {
								continue
							}
							if (g.StartPlayer && g.MinutesPlayed <= uint16(incident.Time)) || (!g.StartPlayer && g.MinutesPlayed > 0) {
								g.GoalsConceded++
							}
							break
						}
						continue
					}
					keeper := incident.FootballPassingNetworkAction[len(incident.FootballPassingNetworkAction)-1].Goalkeeper
					// id := db.GetPlayerByName(strings.TrimSpace(incident.FootballPassingNetworkAction[len(incident.FootballPassingNetworkAction)-1].Goalkeeper.Name)).ID
					id, err := u.db.GetFootballPlayerID(ctx, strings.TrimSpace(keeper.Name), time.Unix(keeper.DateOfBirthTimeStamp, 0))
					if err != nil {
						log.Printf("error in search in incidents away goal: %v, url: %s, incident: %+v\n", err, url, incident)
						continue
					}
					if stats.goalieStatsHome[id] == nil {
						log.Printf("error in increment goals conceded in incidents away goal: %v, url: %s, incident: %+v\n", err, url, incident)
						continue
					}
					stats.goalieStatsHome[id].GoalsConceded++
				}
			} else if incident.IncidentClass == "penalty" {
				if incident.IsHome {
					for _, g := range stats.goalieStatsAway {
						if incident.Time < 0 {
							continue
						}
						if (g.StartPlayer && g.MinutesPlayed <= uint16(incident.Time)) || (!g.StartPlayer && g.MinutesPlayed > 0) {
							g.PenaltyConceded++
							break
						}
					}
				} else {
					for _, g := range stats.goalieStatsHome {
						if incident.Time < 0 {
							continue
						}
						if (g.StartPlayer && g.MinutesPlayed <= uint16(incident.Time)) || (!g.StartPlayer && g.MinutesPlayed > 0) {
							g.PenaltyConceded++
							break
						}
					}
				}
			}
		} else if incident.IncidentType == "card" {
			switch incident.IncidentClass {
			case "yellow":
				if incident.Manager.Name != "" {
					id, err := u.db.GetFootballManagerID(ctx, &models.Manager{FirstName: strings.Split(incident.Manager.Name, " ")[0], LastName: strings.Join(strings.Split(incident.Manager.Name, " ")[1:], " ")})
					if err != nil {
						log.Printf("error in search in yellow manager: %v\n", err)
						continue
					}
					if err := u.producer.Send(ctx, "IncrementYellowCardsManager|"+fmt.Sprint(id), struct{}{}); err != nil {
						log.Printf("error in sending in yellow manager: %v\n", err)
						continue
					}
					// u.db.IncrementYellowCardsManager(ctx, id)
					break
				}
				id, err := u.db.GetFootballPlayerID(ctx, strings.TrimSpace(incident.Player.Name), time.Unix(incident.Player.DateOfBirthTimeStamp, 0))
				if err != nil {
					log.Printf("error in search in incidents yellow card: %v, url: %s, incident: %+v\n", err, url, incident)
					continue
				}
				if incident.IsHome {
					if incident.Player.Position == "G" {
						if stats.goalieStatsHome[id] == nil {
							continue
						}
						stats.goalieStatsHome[id].YellowCards++
					} else {
						if stats.playersStatsHome[id] == nil {
							continue
						}
						stats.playersStatsHome[id].YellowCards++
					}
				} else {
					if incident.Player.Position == "G" {
						if stats.goalieStatsAway[id] == nil {
							continue
						}
						stats.goalieStatsAway[id].YellowCards++
					} else {
						if stats.playersStatsAway[id] == nil {
							continue
						}
						stats.playersStatsAway[id].YellowCards++
					}
				}
			case "yellowRed":
				if incident.Manager.Name != "" {
					id, err := u.db.GetFootballManagerID(ctx, &models.Manager{FirstName: strings.Split(incident.Manager.Name, " ")[0], LastName: strings.Join(strings.Split(incident.Manager.Name, " ")[1:], " ")})
					if err != nil {
						log.Printf("error in search in yellowRed manager: %v\n", err)
						continue
					}
					if err := u.producer.Send(ctx, "IncrementYellowCardsManager|"+fmt.Sprint(id), struct{}{}); err != nil {
						log.Printf("error in sending yellow in yellowred manager: %v\n", err)
						continue
					}
					// u.db.IncrementYellowCardsManager(ctx, id)
					if err := u.producer.Send(ctx, "IncrementRedCardsManager|"+fmt.Sprint(id), struct{}{}); err != nil {
						log.Printf("error in sending red in yellowred manager: %v\n", err)
						continue
					}
					// u.db.IncrementRedCardsManager(ctx, id)
					break
				}
				id, err := u.db.GetFootballPlayerID(ctx, strings.TrimSpace(incident.Player.Name), time.Unix(incident.Player.DateOfBirthTimeStamp, 0))
				if err != nil {
					log.Printf("error in search in incidents yellow red card: %v, url: %s, incident: %+v\n", err, url, incident)
					continue
				}
				if incident.IsHome {
					if incident.Player.Position == "G" {
						if stats.goalieStatsHome[id] == nil {
							continue
						}
						stats.goalieStatsHome[id].YellowCards++
						stats.goalieStatsHome[id].RedCards++
					} else {
						if stats.playersStatsHome[id] == nil {
							continue
						}
						stats.playersStatsHome[id].YellowCards++
						stats.playersStatsHome[id].RedCards++
					}
				} else {
					if incident.Player.Position == "G" {
						if stats.goalieStatsAway[id] == nil {
							continue
						}
						stats.goalieStatsAway[id].YellowCards++
						stats.goalieStatsAway[id].RedCards++
					} else {
						if stats.playersStatsAway[id] == nil {
							continue
						}
						stats.playersStatsAway[id].YellowCards++
						stats.playersStatsAway[id].RedCards++
					}
				}
			case "red":
				if incident.Manager.Name != "" {
					id, err := u.db.GetFootballManagerID(ctx, &models.Manager{FirstName: strings.Split(incident.Manager.Name, " ")[0], LastName: strings.Join(strings.Split(incident.Manager.Name, " ")[1:], " ")})
					if err != nil {
						log.Printf("error in search in red manager: %v\n", err)
						continue
					}
					if err := u.producer.Send(ctx, "IncrementRedCardsManager|"+fmt.Sprint(id), struct{}{}); err != nil {
						log.Printf("error in sending in red manager: %v\n", err)
						continue
					}
					// u.db.IncrementRedCardsManager(ctx, id)
					break
				}
				id, err := u.db.GetFootballPlayerID(ctx, strings.TrimSpace(incident.Player.Name), time.Unix(incident.Player.DateOfBirthTimeStamp, 0))
				if err != nil {
					log.Printf("error in search in incidents red card: %v, url: %s, incident: %+v\n", err, url, incident)
					continue
				}
				if incident.IsHome {
					if incident.Player.Position == "G" {
						if stats.goalieStatsHome[id] == nil {
							continue
						}
						stats.goalieStatsHome[id].RedCards++
					} else {
						if stats.playersStatsHome[id] == nil {
							continue
						}
						stats.playersStatsHome[id].RedCards++
					}
				} else {
					if incident.Player.Position == "G" {
						if stats.goalieStatsAway[id] == nil {
							continue
						}
						stats.goalieStatsAway[id].RedCards++
					} else {
						if stats.playersStatsAway[id] == nil {
							continue
						}
						stats.playersStatsAway[id].RedCards++
					}
				}
			}
		}
	}
}

func (u *Updater) fetchPlayer(ctx context.Context, idPlayer int) (*models.Player, error) {
	url := fmt.Sprintf("https://api.sofascore.com/api/v1/player/%d", idPlayer)
	body, err := u.api.FetchBodyConc(ctx, url)
	if err != nil {
		return nil, err
	}

	var playerResp PlayerResponse
	if len(body) < 5 {
		return nil, fmt.Errorf("too short body of response")
	}
	body = strings.TrimSuffix(body[5:], `</pre><div class="json-formatter-container"></div>`)
	if err := json.Unmarshal([]byte(body), &playerResp); err != nil {
		return nil, err
	}

	var res models.Player

	if playerResp.Player.Retired {
		res.CurrentStatus = "Retired"
	} else {
		res.CurrentStatus = "Active"
	}
	res.DateOfBirth = time.Unix(playerResp.Player.DateOfBirthTimeStamp, 0)
	names := strings.Split(playerResp.Player.Name, " ")
	res.FirstName = strings.TrimSpace(names[0])
	res.LastName = strings.TrimSpace(strings.Join(names[1:], " "))
	if res.LastName == "" {
		res.LastName = res.FirstName
		res.FirstName = ""
	}
	res.Height = uint16(playerResp.Player.Height)
	res.Nation = playerResp.Player.Country.Name
	res.Position = playerResp.Player.Position
	res.PreferredFoot = playerResp.Player.PreferredFoot

	return &res, nil
}

type PlayerResponse struct {
	Player struct {
		Name                 string `json:"name"`
		Position             string `json:"position"`
		Height               int    `json:"height"`
		PreferredFoot        string `json:"preferredFoot"`
		Retired              bool   `json:"retired"`
		DateOfBirthTimeStamp int64  `json:"dateOfBirthTimestamp"`
		Country              struct {
			Name string `json:"name"`
		} `json:"country"`
	} `json:"player"`
}

func (u *Updater) addStatsToDB(ctx context.Context, statsFromMatch StatsFromMatch, matchID uint32) {
	if id, err := u.db.GetFootballMatchStats(ctx, statsFromMatch.teamStats, matchID); err != nil || id == 0 {
		// _, err := u.db.InsertFootballMatchStats(ctx, statsFromMatch.teamStats, matchID)
		// if err != nil {
		// 	log.Printf("error in inserting team stats: %v\n", err)
		// }

		if err := u.producer.Send(ctx, "InsertFootballMatchStats|"+fmt.Sprint(matchID), statsFromMatch.teamStats); err != nil {
			log.Printf("error in sending team stats: %v\n", err)
		}
	}

	// if err := u.db.InsertFootballGoalieMatchStatsBatch(ctx, statsFromMatch.goalieStatsAway, matchID); err != nil {
	// 	log.Printf("error in inserting away goalie stats: %v\n", err)
	// }

	if err := u.producer.Send(ctx, "InsertFootballGoalieMatchStatsBatch|"+fmt.Sprint(matchID), statsFromMatch.goalieStatsAway); err != nil {
		log.Printf("error in sending away goalie stats: %v\n", err)
	}

	// if err := u.db.InsertFootballGoalieMatchStatsBatch(ctx, statsFromMatch.goalieStatsHome, matchID); err != nil {
	// 	log.Printf("error in inserting home goalie stats: %v\n", err)
	// }

	if err := u.producer.Send(ctx, "InsertFootballGoalieMatchStatsBatch|"+fmt.Sprint(matchID), statsFromMatch.goalieStatsHome); err != nil {
		log.Printf("error in sending home goalie stats: %v\n", err)
	}

	// if err := u.db.InsertFootballPlayerMatchStatsBatch(ctx, statsFromMatch.playersStatsAway, matchID); err != nil {
	// 	log.Printf("error in inserting away player stats: %v\n", err)
	// }

	if err := u.producer.Send(ctx, "InsertFootballPlayerMatchStatsBatch|"+fmt.Sprint(matchID), statsFromMatch.playersStatsAway); err != nil {
		log.Printf("error in sending away player stats: %v\n", err)
	}

	// if err := u.db.InsertFootballPlayerMatchStatsBatch(ctx, statsFromMatch.playersStatsHome, matchID); err != nil {
	// 	log.Printf("error in inserting home player stats: %v\n", err)
	// }

	if err := u.producer.Send(ctx, "InsertFootballPlayerMatchStatsBatch|"+fmt.Sprint(matchID), statsFromMatch.playersStatsHome); err != nil {
		log.Printf("error in sending home player stats: %v\n", err)
	}
}
