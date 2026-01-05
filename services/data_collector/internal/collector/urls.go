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
	"2025/2026": {
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
	"2024/2025": {
		"English Premier League": "61627",
		"Ligue 1":                "61736",
		"Bundesliga":             "63516",
		"LA Liga":                "61643",
		"Serie A":                "63515",
		"Russian Premier League": "61712",
		"Liga Portugal":          "63670",
		"Eredevisie":             "61666",
		"Super Lig":              "63814",
	},
	"2023/2024": {
		"English Premier League": "52186",
		"Ligue 1":                "52571",
		"Bundesliga":             "52608",
		"LA Liga":                "52376",
		"Serie A":                "52760",
		"Russian Premier League": "52470",
		"Liga Portugal":          "52769",
		"Eredevisie":             "52554",
		"Super Lig":              "53190",
	},
	"2022/2023": {
		"English Premier League": "41886",
		"Ligue 1":                "42273",
		"Bundesliga":             "42268",
		"LA Liga":                "42409",
		"Serie A":                "42415",
		"Russian Premier League": "42388",
		"Liga Portugal":          "42655",
		"Eredevisie":             "42256",
		"Super Lig":              "42632",
	},
	"2021/2022": {
		"English Premier League": "37036",
		"Ligue 1":                "37167",
		"Bundesliga":             "37166",
		"LA Liga":                "37223",
		"Serie A":                "37475",
		"Russian Premier League": "37038",
		"Liga Portugal":          "37358",
		"Eredevisie":             "36890",
		"Super Lig":              "37466",
	},
	"2020/2021": {
		"English Premier League": "29415",
		"Ligue 1":                "28222",
		"Bundesliga":             "28210",
		"LA Liga":                "32501",
		"Serie A":                "32523",
		"Russian Premier League": "29200",
		"Liga Portugal":          "32456",
		"Eredevisie":             "29186",
		"Super Lig":              "29506",
	},
	"2019/2020": {
		"English Premier League": "23776",
		"Ligue 1":                "23872",
		"Bundesliga":             "23538",
		"LA Liga":                "24127",
		"Serie A":                "24644",
		"Russian Premier League": "23682",
		"Liga Portugal":          "24150",
		"Eredevisie":             "23873",
		"Super Lig":              "24407",
	},
	"2018/2019": {
		"English Premier League": "17359",
		"Ligue 1":                "17279",
		"Bundesliga":             "17597",
		"LA Liga":                "18020",
		"Serie A":                "17932",
		"Russian Premier League": "17753",
		"Liga Portugal":          "17714",
		"Eredevisie":             "17353",
		"Super Lig":              "17762",
	},
	"2017/2018": {
		"English Premier League": "13380",
		"Ligue 1":                "13384",
		"Bundesliga":             "13477",
		"LA Liga":                "13662",
		"Serie A":                "13768",
		"Russian Premier League": "13387",
		"Liga Portugal":          "13539",
		"Eredevisie":             "13399",
		"Super Lig":              "13575",
	},
	"2016/2017": {
		"English Premier League": "11733",
		"Ligue 1":                "11648",
		"Bundesliga":             "11818",
		"LA Liga":                "11906",
		"Serie A":                "11966",
		"Russian Premier League": "11868",
		"Liga Portugal":          "11924",
		"Eredevisie":             "11777",
		"Super Lig":              "11927",
	},
	"2015/2016": {
		"English Premier League": "10356",
		"Ligue 1":                "10373",
		"Bundesliga":             "10419",
		"LA Liga":                "10495",
		"Serie A":                "10596",
		"Russian Premier League": "10407",
		"Liga Portugal":          "10453",
		"Eredevisie":             "10370",
		"Super Lig":              "10470",
	},
}
