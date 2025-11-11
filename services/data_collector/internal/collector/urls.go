package collector

const (
	PREMIER_LEAGUE         = "https://api.sofascore.com/api/v1/unique-tournament/17/season/"
	LIGUE_1                = "https://api.sofascore.com/api/v1/unique-tournament/34/season/"
	BUNDESLIGA             = "https://api.sofascore.com/api/v1/unique-tournament/35/season/"
	LA_LIGA                = "https://api.sofascore.com/api/v1/unique-tournament/8/season/"
	SERIE_A                = "https://api.sofascore.com/api/v1/unique-tournament/23/season/"
	RUSSIAN_PREMIER_LEAGUE = "https://api.sofascore.com/api/v1/unique-tournament/203/season/"
	LIGA_PORTUGAL          = "https://api.sofascore.com/api/v1/unique-tournament/238/season/"
	EREDIVISIE             = "https://api.sofascore.com/api/v1/unique-tournament/37/season/"
	SUPER_LIG              = "https://api.sofascore.com/api/v1/unique-tournament/52/season/"
)

var tournaments = map[string]string{
	"English Premier League": PREMIER_LEAGUE,
	"Ligue 1":                LIGUE_1,
	"Bundesliga":             BUNDESLIGA,
	"LA Liga":                LA_LIGA,
	"Serie A":                SERIE_A,
	"Russian Premier League": RUSSIAN_PREMIER_LEAGUE,
	"Liga Portugal":          LIGA_PORTUGAL,
	"Eredevisie":             EREDIVISIE,
	"Super Lig":              SUPER_LIG,
}
var seasonsIDs = map[string]map[string]string{
	"2025/2026": map[string]string{
		"English Premier League": "76986",
		"Ligue 1":                "77356",
		"Bundesliga":             "77333",
		"LA Liga":                "77559",
		"Serie A":                "76457",
		"Russian Premier League": "77142",
		"Liga Portugal":          "77806",
		"Eredevisie":             "77012",
		"Super Lig":              "77805",
	},
}
