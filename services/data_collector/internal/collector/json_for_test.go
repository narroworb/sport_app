package collector

const (
	teamStatsJSON = `xxxxx{
  "statistics": [
    {
      "period": "ALL",
      "groups": [
        {
          "groupName": "Match overview",
          "statisticsItems": [
            {
              "name": "Ball possession",
              "home": "75%",
              "away": "25%",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 75,
              "awayValue": 25,
              "renderType": 2,
              "key": "ballPossession"
            },
            {
              "name": "Total shots",
              "home": "19",
              "away": "3",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 19,
              "awayValue": 3,
              "renderType": 1,
              "key": "totalShotsOnGoal"
            },
            {
              "name": "Goalkeeper saves",
              "home": "1",
              "away": "9",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 9,
              "renderType": 1,
              "key": "goalkeeperSaves"
            },
            {
              "name": "Corner kicks",
              "home": "14",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 14,
              "awayValue": 2,
              "renderType": 1,
              "key": "cornerKicks"
            },
            {
              "name": "Fouls",
              "home": "6",
              "away": "9",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 6,
              "awayValue": 9,
              "renderType": 1,
              "key": "fouls"
            },
            {
              "name": "Passes",
              "home": "783",
              "away": "269",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 783,
              "awayValue": 269,
              "renderType": 1,
              "key": "passes"
            },
            {
              "name": "Free kicks",
              "home": "9",
              "away": "5",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 9,
              "awayValue": 5,
              "renderType": 1,
              "key": "freeKicks"
            },
            {
              "name": "Yellow cards",
              "home": "0",
              "away": "4",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 4,
              "renderType": 1,
              "key": "yellowCards"
            },
            {
              "name": "Red cards",
              "home": "0",
              "away": "1",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 1,
              "renderType": 1,
              "key": "redCards"
            }
          ]
        },
        {
          "groupName": "Shots",
          "statisticsItems": [
            {
              "name": "Total shots",
              "home": "19",
              "away": "3",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 19,
              "awayValue": 3,
              "renderType": 1,
              "key": "totalShotsOnGoal"
            },
            {
              "name": "Shots on target",
              "home": "10",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 10,
              "awayValue": 2,
              "renderType": 1,
              "key": "shotsOnGoal"
            },
            {
              "name": "Hit woodwork",
              "home": "1",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 0,
              "renderType": 1,
              "key": "hitWoodwork"
            },
            {
              "name": "Shots off target",
              "home": "5",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 5,
              "awayValue": 0,
              "renderType": 1,
              "key": "shotsOffGoal"
            },
            {
              "name": "Blocked shots",
              "home": "4",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 4,
              "awayValue": 1,
              "renderType": 1,
              "key": "blockedScoringAttempt"
            },
            {
              "name": "Shots inside box",
              "home": "13",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 13,
              "awayValue": 1,
              "renderType": 1,
              "key": "totalShotsInsideBox"
            },
            {
              "name": "Shots outside box",
              "home": "6",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 6,
              "awayValue": 2,
              "renderType": 1,
              "key": "totalShotsOutsideBox"
            }
          ]
        },
        {
          "groupName": "Attack",
          "statisticsItems": [
            {
              "name": "Through balls",
              "home": "1",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 0,
              "renderType": 1,
              "key": "accurateThroughBall"
            },
            {
              "name": "Offsides",
              "home": "0",
              "away": "4",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 4,
              "renderType": 1,
              "key": "offsides"
            }
          ]
        },
        {
          "groupName": "Passes",
          "statisticsItems": [
            {
              "name": "Accurate passes",
              "home": "684",
              "away": "179",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 684,
              "awayValue": 179,
              "renderType": 1,
              "key": "accuratePasses"
            },
            {
              "name": "Throw-ins",
              "home": "33",
              "away": "19",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 33,
              "awayValue": 19,
              "renderType": 1,
              "key": "throwIns"
            }
          ]
        },
        {
          "groupName": "Duels",
          "statisticsItems": [
            {
              "name": "Duels",
              "home": "60%",
              "away": "40%",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 60,
              "awayValue": 40,
              "renderType": 2,
              "key": "duelWonPercent"
            },
            {
              "name": "Ground duels",
              "home": "36/62 (58%)",
              "away": "26/62 (42%)",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "team",
              "homeValue": 36,
              "awayValue": 26,
              "homeTotal": 62,
              "awayTotal": 62,
              "renderType": 3,
              "key": "groundDuelsPercentage"
            },
            {
              "name": "Aerial duels",
              "home": "16/25 (64%)",
              "away": "9/25 (36%)",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "team",
              "homeValue": 16,
              "awayValue": 9,
              "homeTotal": 25,
              "awayTotal": 25,
              "renderType": 3,
              "key": "aerialDuelsPercentage"
            }
          ]
        },
        {
          "groupName": "Goalkeeping",
          "statisticsItems": [
            {
              "name": "Total saves",
              "home": "1",
              "away": "9",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 9,
              "renderType": 1,
              "key": "goalkeeperSaves"
            },
            {
              "name": "Goal kicks",
              "home": "2",
              "away": "14",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 2,
              "awayValue": 14,
              "renderType": 1,
              "key": "goalKicks"
            }
          ]
        }
      ]
    },
    {
      "period": "1ST",
      "groups": [
        {
          "groupName": "Match overview",
          "statisticsItems": [
            {
              "name": "Ball possession",
              "home": "74%",
              "away": "26%",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 74,
              "awayValue": 26,
              "renderType": 2,
              "key": "ballPossession"
            },
            {
              "name": "Total shots",
              "home": "7",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 7,
              "awayValue": 1,
              "renderType": 1,
              "key": "totalShotsOnGoal"
            },
            {
              "name": "Goalkeeper saves",
              "home": "0",
              "away": "2",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 2,
              "renderType": 1,
              "key": "goalkeeperSaves"
            },
            {
              "name": "Corner kicks",
              "home": "4",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 4,
              "awayValue": 1,
              "renderType": 1,
              "key": "cornerKicks"
            },
            {
              "name": "Passes",
              "home": "431",
              "away": "149",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 431,
              "awayValue": 149,
              "renderType": 1,
              "key": "passes"
            },
            {
              "name": "Free kicks",
              "home": "6",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 6,
              "awayValue": 2,
              "renderType": 1,
              "key": "freeKicks"
            },
            {
              "name": "Yellow cards",
              "home": "0",
              "away": "2",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 2,
              "renderType": 1,
              "key": "yellowCards"
            },
            {
              "name": "Red cards",
              "home": "0",
              "away": "0",
              "compareCode": 3,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 0,
              "renderType": 1,
              "key": "redCards"
            }
          ]
        },
        {
          "groupName": "Shots",
          "statisticsItems": [
            {
              "name": "Total shots",
              "home": "7",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 7,
              "awayValue": 1,
              "renderType": 1,
              "key": "totalShotsOnGoal"
            },
            {
              "name": "Shots on target",
              "home": "2",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 2,
              "awayValue": 1,
              "renderType": 1,
              "key": "shotsOnGoal"
            },
            {
              "name": "Hit woodwork",
              "home": "1",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 0,
              "renderType": 1,
              "key": "hitWoodwork"
            },
            {
              "name": "Shots off target",
              "home": "2",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 2,
              "awayValue": 0,
              "renderType": 1,
              "key": "shotsOffGoal"
            },
            {
              "name": "Blocked shots",
              "home": "3",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 3,
              "awayValue": 0,
              "renderType": 1,
              "key": "blockedScoringAttempt"
            },
            {
              "name": "Shots inside box",
              "home": "5",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 5,
              "awayValue": 1,
              "renderType": 1,
              "key": "totalShotsInsideBox"
            },
            {
              "name": "Shots outside box",
              "home": "2",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 2,
              "awayValue": 0,
              "renderType": 1,
              "key": "totalShotsOutsideBox"
            }
          ]
        },
        {
          "groupName": "Attack",
          "statisticsItems": [
            {
              "name": "Through balls",
              "home": "1",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 0,
              "renderType": 1,
              "key": "accurateThroughBall"
            },
            {
              "name": "Offsides",
              "home": "0",
              "away": "1",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 1,
              "renderType": 1,
              "key": "offsides"
            }
          ]
        },
        {
          "groupName": "Passes",
          "statisticsItems": [
            {
              "name": "Accurate passes",
              "home": "383",
              "away": "102",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 383,
              "awayValue": 102,
              "renderType": 1,
              "key": "accuratePasses"
            },
            {
              "name": "Throw-ins",
              "home": "15",
              "away": "8",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 15,
              "awayValue": 8,
              "renderType": 1,
              "key": "throwIns"
            }
          ]
        },
        {
          "groupName": "Duels",
          "statisticsItems": [
            {
              "name": "Duels",
              "home": "58%",
              "away": "42%",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 58,
              "awayValue": 42,
              "renderType": 2,
              "key": "duelWonPercent"
            },
            {
              "name": "Ground duels",
              "home": "16/28 (57%)",
              "away": "12/28 (43%)",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "team",
              "homeValue": 16,
              "awayValue": 12,
              "homeTotal": 28,
              "awayTotal": 28,
              "renderType": 3,
              "key": "groundDuelsPercentage"
            },
            {
              "name": "Aerial duels",
              "home": "9/15 (60%)",
              "away": "6/15 (40%)",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "team",
              "homeValue": 9,
              "awayValue": 6,
              "homeTotal": 15,
              "awayTotal": 15,
              "renderType": 3,
              "key": "aerialDuelsPercentage"
            }
          ]
        },
        {
          "groupName": "Goalkeeping",
          "statisticsItems": [
            {
              "name": "Total saves",
              "home": "0",
              "away": "2",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 2,
              "renderType": 1,
              "key": "goalkeeperSaves"
            },
            {
              "name": "Goal kicks",
              "home": "1",
              "away": "7",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 7,
              "renderType": 1,
              "key": "goalKicks"
            }
          ]
        }
      ]
    },
    {
      "period": "2ND",
      "groups": [
        {
          "groupName": "Match overview",
          "statisticsItems": [
            {
              "name": "Ball possession",
              "home": "75%",
              "away": "25%",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 75,
              "awayValue": 25,
              "renderType": 2,
              "key": "ballPossession"
            },
            {
              "name": "Total shots",
              "home": "12",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 12,
              "awayValue": 2,
              "renderType": 1,
              "key": "totalShotsOnGoal"
            },
            {
              "name": "Goalkeeper saves",
              "home": "1",
              "away": "7",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 7,
              "renderType": 1,
              "key": "goalkeeperSaves"
            },
            {
              "name": "Corner kicks",
              "home": "10",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 10,
              "awayValue": 1,
              "renderType": 1,
              "key": "cornerKicks"
            },
            {
              "name": "Passes",
              "home": "352",
              "away": "120",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 352,
              "awayValue": 120,
              "renderType": 1,
              "key": "passes"
            },
            {
              "name": "Free kicks",
              "home": "3",
              "away": "3",
              "compareCode": 3,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 3,
              "awayValue": 3,
              "renderType": 1,
              "key": "freeKicks"
            },
            {
              "name": "Yellow cards",
              "home": "0",
              "away": "2",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 2,
              "renderType": 1,
              "key": "yellowCards"
            },
            {
              "name": "Red cards",
              "home": "0",
              "away": "1",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 1,
              "renderType": 1,
              "key": "redCards"
            }
          ]
        },
        {
          "groupName": "Shots",
          "statisticsItems": [
            {
              "name": "Total shots",
              "home": "12",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 12,
              "awayValue": 2,
              "renderType": 1,
              "key": "totalShotsOnGoal"
            },
            {
              "name": "Shots on target",
              "home": "8",
              "away": "1",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 8,
              "awayValue": 1,
              "renderType": 1,
              "key": "shotsOnGoal"
            },
            {
              "name": "Hit woodwork",
              "home": "0",
              "away": "0",
              "compareCode": 3,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 0,
              "renderType": 1,
              "key": "hitWoodwork"
            },
            {
              "name": "Shots off target",
              "home": "3",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 3,
              "awayValue": 0,
              "renderType": 1,
              "key": "shotsOffGoal"
            },
            {
              "name": "Blocked shots",
              "home": "1",
              "away": "1",
              "compareCode": 3,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 1,
              "renderType": 1,
              "key": "blockedScoringAttempt"
            },
            {
              "name": "Shots inside box",
              "home": "8",
              "away": "0",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 8,
              "awayValue": 0,
              "renderType": 1,
              "key": "totalShotsInsideBox"
            },
            {
              "name": "Shots outside box",
              "home": "4",
              "away": "2",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 4,
              "awayValue": 2,
              "renderType": 1,
              "key": "totalShotsOutsideBox"
            }
          ]
        },
        {
          "groupName": "Attack",
          "statisticsItems": [
            {
              "name": "Through balls",
              "home": "0",
              "away": "0",
              "compareCode": 3,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 0,
              "renderType": 1,
              "key": "accurateThroughBall"
            },
            {
              "name": "Offsides",
              "home": "0",
              "away": "3",
              "compareCode": 2,
              "statisticsType": "negative",
              "valueType": "event",
              "homeValue": 0,
              "awayValue": 3,
              "renderType": 1,
              "key": "offsides"
            }
          ]
        },
        {
          "groupName": "Passes",
          "statisticsItems": [
            {
              "name": "Accurate passes",
              "home": "301",
              "away": "77",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 301,
              "awayValue": 77,
              "renderType": 1,
              "key": "accuratePasses"
            },
            {
              "name": "Throw-ins",
              "home": "18",
              "away": "11",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 18,
              "awayValue": 11,
              "renderType": 1,
              "key": "throwIns"
            }
          ]
        },
        {
          "groupName": "Duels",
          "statisticsItems": [
            {
              "name": "Duels",
              "home": "61%",
              "away": "39%",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 61,
              "awayValue": 39,
              "renderType": 2,
              "key": "duelWonPercent"
            },
            {
              "name": "Ground duels",
              "home": "20/34 (59%)",
              "away": "14/34 (41%)",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "team",
              "homeValue": 20,
              "awayValue": 14,
              "homeTotal": 34,
              "awayTotal": 34,
              "renderType": 3,
              "key": "groundDuelsPercentage"
            },
            {
              "name": "Aerial duels",
              "home": "7/10 (70%)",
              "away": "3/10 (30%)",
              "compareCode": 1,
              "statisticsType": "positive",
              "valueType": "team",
              "homeValue": 7,
              "awayValue": 3,
              "homeTotal": 10,
              "awayTotal": 10,
              "renderType": 3,
              "key": "aerialDuelsPercentage"
            }
          ]
        },
        {
          "groupName": "Goalkeeping",
          "statisticsItems": [
            {
              "name": "Total saves",
              "home": "1",
              "away": "7",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 7,
              "renderType": 1,
              "key": "goalkeeperSaves"
            },
            {
              "name": "Goal kicks",
              "home": "1",
              "away": "7",
              "compareCode": 2,
              "statisticsType": "positive",
              "valueType": "event",
              "homeValue": 1,
              "awayValue": 7,
              "renderType": 1,
              "key": "goalKicks"
            }
          ]
        }
      ]
    }
  ]
}</pre><div class="json-formatter-container"></div>`

	playersStatsJSON = `xxxxx{
  "confirmed": true,
  "home": {
    "players": [
      {
        "player": {
          "name": "Pepe Reina",
          "slug": "pepe-reina",
          "shortName": "P. Reina",
          "position": "G",
          "height": 188,
          "userCount": 1338,
          "gender": "M",
          "id": 26937,
          "country": {
            "alpha2": "ES",
            "alpha3": "ESP",
            "name": "Spain",
            "slug": "spain"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 399600000,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "بيبي رينا",
              "hi": "पेपे रीना",
              "bn": "পেপে রেইনা"
            },
            "shortNameTranslation": {
              "ar": "ب. رينا",
              "hi": "पी. रीना",
              "bn": "পি. রেইনা"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 25,
        "jerseyNumber": "25",
        "position": "G",
        "substitute": false,
        "statistics": {
          "totalPass": 18,
          "accuratePass": 18,
          "totalLongBalls": 2,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 17,
          "totalOwnHalfPasses": 17,
          "accurateOppositionHalfPasses": 1,
          "totalOppositionHalfPasses": 1,
          "ballRecovery": 5,
          "saves": 1,
          "minutesPlayed": 90,
          "touches": 20,
          "rating": 6.5,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Elseid Hysaj",
          "slug": "elseid-hysaj",
          "shortName": "E. Hysaj",
          "position": "D",
          "jerseyNumber": "23",
          "height": 182,
          "userCount": 1321,
          "gender": "M",
          "id": 136322,
          "country": {
            "alpha2": "AL",
            "alpha3": "ALB",
            "name": "Albania",
            "slug": "albania"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 760147200,
          "proposedMarketValueRaw": {
            "value": 1500000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "إلسيد هيساج",
              "hi": "एल्सीड ह्यसाज",
              "bn": "এলসিদ হাইসাজ"
            },
            "shortNameTranslation": {
              "ar": "إ. هيساج",
              "hi": "ई. ह्यसाज",
              "bn": "ই. হাইসাজ"
            }
          }
        },
        "teamId": 2699,
        "shirtNumber": 2,
        "jerseyNumber": "2",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 57,
          "accuratePass": 53,
          "totalLongBalls": 3,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 17,
          "totalOwnHalfPasses": 17,
          "accurateOppositionHalfPasses": 36,
          "totalOppositionHalfPasses": 45,
          "totalCross": 5,
          "aerialLost": 2,
          "aerialWon": 1,
          "duelLost": 3,
          "duelWon": 6,
          "dispossessed": 1,
          "onTargetScoringAttempt": 1,
          "interceptionWon": 2,
          "ballRecovery": 9,
          "totalTackle": 4,
          "unsuccessfulTouch": 1,
          "wasFouled": 1,
          "minutesPlayed": 90,
          "touches": 84,
          "rating": 7.2,
          "totalShots": 1,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Nikola Maksimović",
          "slug": "nikola-maksimovic",
          "shortName": "N. Maksimović",
          "position": "D",
          "jerseyNumber": "15",
          "height": 193,
          "userCount": 274,
          "gender": "M",
          "id": 68076,
          "country": {
            "alpha2": "RS",
            "alpha3": "SRB",
            "name": "Serbia",
            "slug": "serbia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 691027200,
          "proposedMarketValueRaw": {
            "value": 275000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "نيكولا ماكسيموفيش",
              "hi": "निकोला मक्सिमोविच",
              "bn": "নিকোলা মাকসিমোভিচ"
            },
            "shortNameTranslation": {
              "ar": "ن. ماكسيموفيش",
              "hi": "एन. मक्सिमोविच",
              "bn": "এন. মাকসিমোভিচ"
            }
          }
        },
        "teamId": 118516,
        "shirtNumber": 19,
        "jerseyNumber": "19",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 81,
          "accuratePass": 72,
          "totalLongBalls": 3,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 36,
          "totalOwnHalfPasses": 39,
          "accurateOppositionHalfPasses": 36,
          "totalOppositionHalfPasses": 43,
          "totalCross": 1,
          "aerialLost": 3,
          "aerialWon": 4,
          "duelLost": 3,
          "duelWon": 6,
          "totalContest": 1,
          "wonContest": 1,
          "totalClearance": 2,
          "outfielderBlock": 1,
          "interceptionWon": 3,
          "ballRecovery": 7,
          "totalTackle": 1,
          "fouls": 1,
          "minutesPlayed": 90,
          "touches": 90,
          "rating": 7,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Raúl Albiol",
          "firstName": "",
          "lastName": "",
          "slug": "raul-albiol",
          "shortName": "R. Albiol",
          "position": "D",
          "jerseyNumber": "39",
          "height": 190,
          "userCount": 1073,
          "gender": "M",
          "id": 3041,
          "country": {
            "alpha2": "ES",
            "alpha3": "ESP",
            "name": "Spain",
            "slug": "spain"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 494640000,
          "proposedMarketValueRaw": {
            "value": 950000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "راؤول ألبيول",
              "hi": "राउल एल्बिओल",
              "bn": "রাউল আলবিওল"
            },
            "shortNameTranslation": {
              "ar": "ر. ألبيول",
              "hi": "आर. एल्बिओल",
              "bn": "আর. আলবিওল"
            }
          }
        },
        "teamId": 2737,
        "shirtNumber": 33,
        "jerseyNumber": "33",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 112,
          "accuratePass": 100,
          "totalLongBalls": 12,
          "accurateLongBalls": 6,
          "accurateOwnHalfPasses": 26,
          "totalOwnHalfPasses": 30,
          "accurateOppositionHalfPasses": 74,
          "totalOppositionHalfPasses": 82,
          "aerialWon": 7,
          "duelLost": 3,
          "duelWon": 9,
          "challengeLost": 2,
          "totalContest": 1,
          "wonContest": 1,
          "shotOffTarget": 2,
          "totalClearance": 6,
          "interceptionWon": 2,
          "ballRecovery": 11,
          "totalTackle": 1,
          "fouls": 1,
          "minutesPlayed": 90,
          "touches": 127,
          "rating": 7.1,
          "keyPass": 1,
          "totalShots": 2,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Faouzi Ghoulam",
          "slug": "faouzi-ghoulam",
          "shortName": "F. Ghoulam",
          "position": "D",
          "height": 188,
          "userCount": 2055,
          "gender": "M",
          "id": 133236,
          "country": {
            "alpha2": "DZ",
            "alpha3": "DZA",
            "name": "Algeria",
            "slug": "algeria"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 665366400,
          "proposedMarketValueRaw": {
            "value": 210000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "فوزي غلام",
              "hi": "फ़ौज़ी ग़ुलाम",
              "bn": "ফৌজি গোলাম"
            },
            "shortNameTranslation": {
              "ar": "ف. غلام",
              "hi": "एफ. ग़ुलाम",
              "bn": "এফ. গোলাম"
            }
          }
        },
        "teamId": 1109346,
        "shirtNumber": 31,
        "jerseyNumber": "31",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 65,
          "accuratePass": 54,
          "totalLongBalls": 2,
          "accurateLongBalls": 1,
          "accurateOwnHalfPasses": 11,
          "totalOwnHalfPasses": 12,
          "accurateOppositionHalfPasses": 46,
          "totalOppositionHalfPasses": 76,
          "totalCross": 23,
          "accurateCross": 3,
          "aerialWon": 3,
          "duelLost": 1,
          "duelWon": 6,
          "dispossessed": 1,
          "interceptionWon": 1,
          "ballRecovery": 6,
          "totalTackle": 2,
          "unsuccessfulTouch": 2,
          "wasFouled": 1,
          "minutesPlayed": 90,
          "touches": 112,
          "rating": 6.6,
          "keyPass": 1,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Allan",
          "slug": "allan",
          "shortName": "Allan",
          "position": "M",
          "jerseyNumber": "25",
          "height": 175,
          "userCount": 2949,
          "gender": "M",
          "id": 158277,
          "country": {
            "alpha2": "BR",
            "alpha3": "BRA",
            "name": "Brazil",
            "slug": "brazil"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 663292800,
          "proposedMarketValueRaw": {
            "value": 1900000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ألان",
              "hi": "एलन",
              "bn": "অ্যালান"
            },
            "shortNameTranslation": {
              "ar": "ألان",
              "hi": "एलन",
              "bn": "অ্যালান"
            }
          }
        },
        "teamId": 1958,
        "shirtNumber": 5,
        "jerseyNumber": "5",
        "position": "M",
        "substitute": false,
        "statistics": {
          "totalPass": 27,
          "accuratePass": 23,
          "accurateOwnHalfPasses": 6,
          "totalOwnHalfPasses": 6,
          "accurateOppositionHalfPasses": 17,
          "totalOppositionHalfPasses": 22,
          "totalCross": 1,
          "aerialLost": 2,
          "duelLost": 5,
          "duelWon": 2,
          "challengeLost": 1,
          "dispossessed": 2,
          "ballRecovery": 2,
          "totalTackle": 2,
          "minutesPlayed": 54,
          "touches": 34,
          "rating": 6.1,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Jorginho",
          "slug": "jorginho",
          "shortName": "Jorginho",
          "position": "M",
          "jerseyNumber": "21",
          "height": 180,
          "userCount": 30295,
          "gender": "M",
          "sofascoreId": "jorginhofrello",
          "id": 132874,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 693187200,
          "proposedMarketValueRaw": {
            "value": 6500000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "جورجينيو",
              "hi": "जोर्गिनहो",
              "bn": "জর্গিনহো"
            },
            "shortNameTranslation": {
              "ar": "جورجينيو",
              "hi": "जोर्गिनहो",
              "bn": "জর্গিনহো"
            }
          }
        },
        "teamId": 5981,
        "shirtNumber": 8,
        "jerseyNumber": "8",
        "position": "M",
        "substitute": false,
        "statistics": {
          "totalPass": 105,
          "accuratePass": 91,
          "totalLongBalls": 9,
          "accurateLongBalls": 5,
          "accurateOwnHalfPasses": 19,
          "totalOwnHalfPasses": 20,
          "accurateOppositionHalfPasses": 72,
          "totalOppositionHalfPasses": 87,
          "totalCross": 2,
          "duelLost": 2,
          "duelWon": 4,
          "challengeLost": 1,
          "ballRecovery": 3,
          "totalTackle": 3,
          "wasFouled": 1,
          "fouls": 1,
          "minutesPlayed": 63,
          "touches": 114,
          "rating": 7.7,
          "keyPass": 3,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Marek Hamšík",
          "slug": "marek-hamsik",
          "shortName": "M. Hamšík",
          "position": "M",
          "jerseyNumber": "17",
          "height": 183,
          "userCount": 1392,
          "id": 25985,
          "country": {
            "alpha2": "SK",
            "alpha3": "SVK",
            "name": "Slovakia",
            "slug": "slovakia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 554342400,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ماريك هامشيك",
              "hi": "मारेक हामसिक",
              "bn": "মারেক হ্যামসিক"
            },
            "shortNameTranslation": {
              "ar": "م. هامشيك",
              "hi": "एम. हामसिक",
              "bn": "এম. হামসিক"
            }
          }
        },
        "teamId": 828626,
        "shirtNumber": 17,
        "jerseyNumber": "17",
        "position": "M",
        "substitute": false,
        "captain": true,
        "statistics": {
          "totalPass": 127,
          "accuratePass": 118,
          "totalLongBalls": 16,
          "accurateLongBalls": 14,
          "accurateOwnHalfPasses": 18,
          "totalOwnHalfPasses": 19,
          "accurateOppositionHalfPasses": 102,
          "totalOppositionHalfPasses": 111,
          "totalCross": 3,
          "accurateCross": 2,
          "aerialWon": 1,
          "duelLost": 3,
          "duelWon": 9,
          "dispossessed": 2,
          "totalContest": 2,
          "wonContest": 2,
          "bigChanceCreated": 1,
          "ballRecovery": 13,
          "totalTackle": 3,
          "wasFouled": 3,
          "fouls": 1,
          "minutesPlayed": 90,
          "touches": 142,
          "rating": 8.4,
          "keyPass": 4,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "José Callejón",
          "slug": "jose-callejon",
          "shortName": "J. Callejón",
          "position": "F",
          "height": 178,
          "userCount": 770,
          "gender": "M",
          "id": 40965,
          "country": {
            "alpha2": "ES",
            "alpha3": "ESP",
            "name": "Spain",
            "slug": "spain"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 540000000,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "خوسيه كاييخون",
              "hi": "जोस कैलेजोन",
              "bn": "জোসে ক্যালেজন"
            },
            "shortNameTranslation": {
              "ar": "خ. كاييخون",
              "hi": "जे. कैलेजोन",
              "bn": "জে. ক্যালেজন"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 7,
        "jerseyNumber": "7",
        "position": "F",
        "substitute": false,
        "statistics": {
          "totalPass": 42,
          "accuratePass": 39,
          "totalLongBalls": 2,
          "accurateLongBalls": 1,
          "accurateOwnHalfPasses": 5,
          "totalOwnHalfPasses": 5,
          "accurateOppositionHalfPasses": 35,
          "totalOppositionHalfPasses": 47,
          "totalCross": 10,
          "accurateCross": 1,
          "aerialLost": 1,
          "duelLost": 4,
          "duelWon": 2,
          "challengeLost": 2,
          "bigChanceMissed": 1,
          "shotOffTarget": 1,
          "onTargetScoringAttempt": 1,
          "ballRecovery": 5,
          "totalTackle": 1,
          "wasFouled": 1,
          "fouls": 1,
          "minutesPlayed": 90,
          "touches": 66,
          "rating": 6.8,
          "keyPass": 4,
          "totalShots": 2,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Dries Mertens",
          "slug": "dries-mertens",
          "shortName": "D. Mertens",
          "position": "M",
          "height": 169,
          "userCount": 11859,
          "gender": "M",
          "sofascoreId": "driesmertens",
          "id": 32493,
          "country": {
            "alpha2": "BE",
            "alpha3": "BEL",
            "name": "Belgium",
            "slug": "belgium"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 547257600,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "دريس ميرتينز",
              "hi": "ड्रीस मेर्टेंस",
              "bn": "ড্রাইস মার্টেনস"
            },
            "shortNameTranslation": {
              "ar": "د. ميرتنز",
              "hi": "डी. मेर्टेंस",
              "bn": "ডি. মার্টেনস"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 14,
        "jerseyNumber": "14",
        "position": "F",
        "substitute": false,
        "statistics": {
          "totalPass": 29,
          "accuratePass": 21,
          "totalLongBalls": 1,
          "accurateLongBalls": 1,
          "accurateOwnHalfPasses": 1,
          "totalOwnHalfPasses": 2,
          "accurateOppositionHalfPasses": 21,
          "totalOppositionHalfPasses": 30,
          "totalCross": 3,
          "accurateCross": 1,
          "duelLost": 4,
          "duelWon": 5,
          "challengeLost": 1,
          "dispossessed": 2,
          "totalContest": 4,
          "wonContest": 3,
          "bigChanceCreated": 1,
          "shotOffTarget": 1,
          "onTargetScoringAttempt": 4,
          "blockedScoringAttempt": 3,
          "hitWoodwork": 1,
          "goals": 1,
          "ballRecovery": 3,
          "unsuccessfulTouch": 1,
          "wasFouled": 2,
          "minutesPlayed": 90,
          "touches": 54,
          "rating": 8.4,
          "keyPass": 2,
          "totalShots": 8,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Lorenzo Insigne",
          "firstName": "",
          "lastName": "",
          "slug": "lorenzo-insigne",
          "shortName": "L. Insigne",
          "position": "F",
          "height": 163,
          "userCount": 4625,
          "gender": "M",
          "sofascoreId": "IlMagnifico",
          "id": 106258,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 675993600,
          "proposedMarketValueRaw": {
            "value": 1600000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "لورينزو إنسيني",
              "hi": "लोरेंजो इनसिग्ने",
              "bn": "লরেঞ্জো ইনসাইন"
            },
            "shortNameTranslation": {
              "ar": "ل. إنسيني",
              "hi": "एल. इनसिग्ने",
              "bn": "এল. ইনসাইন"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 24,
        "jerseyNumber": "24",
        "position": "F",
        "substitute": false,
        "statistics": {
          "totalPass": 83,
          "accuratePass": 67,
          "totalLongBalls": 2,
          "accurateOwnHalfPasses": 10,
          "totalOwnHalfPasses": 11,
          "accurateOppositionHalfPasses": 57,
          "totalOppositionHalfPasses": 77,
          "totalCross": 5,
          "duelWon": 2,
          "totalContest": 2,
          "wonContest": 2,
          "bigChanceMissed": 1,
          "shotOffTarget": 1,
          "onTargetScoringAttempt": 2,
          "blockedScoringAttempt": 1,
          "ballRecovery": 5,
          "unsuccessfulTouch": 1,
          "minutesPlayed": 90,
          "touches": 96,
          "rating": 7.1,
          "totalShots": 4,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Piotr Zieliński",
          "slug": "piotr-zielinski",
          "shortName": "P. Zieliński",
          "position": "M",
          "jerseyNumber": "7",
          "height": 180,
          "userCount": 18011,
          "gender": "M",
          "sofascoreId": "zielu_94",
          "id": 138605,
          "country": {
            "alpha2": "PL",
            "alpha3": "POL",
            "name": "Poland",
            "slug": "poland"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 769392000,
          "proposedMarketValueRaw": {
            "value": 10300000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "بيوتر زيلينسكي",
              "hi": "पिओट्र ज़िलिंस्की",
              "bn": "পিওতর জিলিনিস্কি"
            },
            "shortNameTranslation": {
              "ar": "ب. زيلينسكي",
              "hi": "पी. ज़िलिंस्की",
              "bn": "পি. জিলিনিস্কি"
            }
          }
        },
        "teamId": 2697,
        "shirtNumber": 20,
        "jerseyNumber": "20",
        "position": "M",
        "substitute": true,
        "statistics": {
          "totalPass": 34,
          "accuratePass": 27,
          "totalLongBalls": 5,
          "accurateLongBalls": 4,
          "goalAssist": 1,
          "accurateOwnHalfPasses": 3,
          "totalOwnHalfPasses": 3,
          "accurateOppositionHalfPasses": 25,
          "totalOppositionHalfPasses": 32,
          "totalCross": 1,
          "accurateCross": 1,
          "aerialLost": 1,
          "duelLost": 6,
          "duelWon": 1,
          "challengeLost": 3,
          "dispossessed": 1,
          "onTargetScoringAttempt": 1,
          "totalTackle": 1,
          "fouls": 1,
          "minutesPlayed": 36,
          "touches": 40,
          "rating": 6.9,
          "keyPass": 2,
          "totalShots": 1,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Leonardo Pavoletti",
          "slug": "leonardo-pavoletti",
          "shortName": "L. Pavoletti",
          "position": "F",
          "jerseyNumber": "30",
          "height": 188,
          "userCount": 523,
          "gender": "M",
          "id": 132704,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 596505600,
          "proposedMarketValueRaw": {
            "value": 640000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ليوناردو بافوليتي",
              "hi": "लियोनार्डो पावोलेटी",
              "bn": "লিওনার্দো পাভোলেটি"
            },
            "shortNameTranslation": {
              "ar": "ل. بافوليتي",
              "hi": "एल. पावोलेटी",
              "bn": "এল. পাভোলেটি"
            }
          }
        },
        "teamId": 2719,
        "shirtNumber": 32,
        "jerseyNumber": "32",
        "position": "F",
        "substitute": true,
        "statistics": {
          "totalPass": 3,
          "accuratePass": 1,
          "accurateOppositionHalfPasses": 1,
          "totalOppositionHalfPasses": 3,
          "duelLost": 1,
          "dispossessed": 1,
          "onTargetScoringAttempt": 1,
          "unsuccessfulTouch": 1,
          "minutesPlayed": 27,
          "touches": 7,
          "rating": 6.6,
          "totalShots": 1,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Rafael Cabral",
          "slug": "rafael-cabral",
          "shortName": "Rafael Cabral",
          "position": "G",
          "jerseyNumber": "1",
          "height": 186,
          "userCount": 372,
          "gender": "M",
          "id": 124980,
          "country": {
            "alpha2": "BR",
            "alpha3": "BRA",
            "name": "Brazil",
            "slug": "brazil"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 643161600,
          "proposedMarketValueRaw": {
            "value": 185000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "رافائيل كابرال"
            },
            "shortNameTranslation": {
              "ar": "ر. كابرال"
            }
          }
        },
        "teamId": 5133,
        "shirtNumber": 1,
        "jerseyNumber": "1",
        "position": "G",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Luigi Sepe",
          "slug": "luigi-sepe",
          "shortName": "L. Sepe",
          "position": "G",
          "height": 185,
          "userCount": 180,
          "gender": "M",
          "id": 46714,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 673660800,
          "proposedMarketValueRaw": {
            "value": 240000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "لويجي سيبي",
              "hi": "लुइगी सेपे",
              "bn": "লুইগি সেপে"
            },
            "shortNameTranslation": {
              "ar": "ل. سيبي",
              "hi": "एल. सेपे",
              "bn": "এল. সেপে"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 22,
        "jerseyNumber": "22",
        "position": "G",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Christian Maggio",
          "slug": "christian-maggio",
          "shortName": "C. Maggio",
          "position": "D",
          "height": 184,
          "userCount": 49,
          "id": 14374,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 382233600,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "كريستيان ماجيو",
              "hi": "क्रिश्चियन मैगियो",
              "bn": "ক্রিশ্চিয়ান ম্যাগিও"
            },
            "shortNameTranslation": {
              "ar": "ك. ماجيو",
              "hi": "सी. मैगियो",
              "bn": "সি. ম্যাজিও"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 11,
        "jerseyNumber": "11",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Marko Rog",
          "slug": "marko-rog",
          "shortName": "M. Rog",
          "position": "M",
          "jerseyNumber": "20",
          "height": 180,
          "userCount": 1223,
          "gender": "M",
          "id": 573616,
          "country": {
            "alpha2": "HR",
            "alpha3": "HRV",
            "name": "Croatia",
            "slug": "croatia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 806112000,
          "proposedMarketValueRaw": {
            "value": 685000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ماركو روغ",
              "hi": "मार्को रोग",
              "bn": "মার্কো রোগ"
            },
            "shortNameTranslation": {
              "ar": "م. روغ",
              "hi": "एम. रोग",
              "bn": "এম. রোগ"
            }
          }
        },
        "teamId": 2719,
        "shirtNumber": 30,
        "jerseyNumber": "30",
        "position": "M",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Amadou Diawara",
          "firstName": "",
          "lastName": "",
          "slug": "amadou-diawara",
          "shortName": "A. Diawara",
          "position": "M",
          "jerseyNumber": "24",
          "height": 183,
          "userCount": 2252,
          "gender": "M",
          "id": 788895,
          "country": {
            "alpha2": "GN",
            "alpha3": "GIN",
            "name": "Guinea",
            "slug": "guinea"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 869097600,
          "proposedMarketValueRaw": {
            "value": 720000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "أمادو دياوارا",
              "hi": "अमादौ दियारा",
              "bn": "আমাদৌ দিওয়ারা"
            },
            "shortNameTranslation": {
              "ar": "أ. دياوارا",
              "hi": "ए. दियारा",
              "bn": "এ. দিওয়ারা"
            }
          }
        },
        "teamId": 2845,
        "shirtNumber": 42,
        "jerseyNumber": "42",
        "position": "M",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Vlad Chiricheş",
          "slug": "vlad-chiriches",
          "shortName": "V. Chiricheş",
          "position": "M",
          "jerseyNumber": "21",
          "height": 184,
          "userCount": 768,
          "gender": "M",
          "id": 69295,
          "country": {
            "alpha2": "RO",
            "alpha3": "ROU",
            "name": "Romania",
            "slug": "romania"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 627004800,
          "proposedMarketValueRaw": {
            "value": 270000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "فلاد كيريكيش",
              "hi": "व्लाद चिरिचेस",
              "bn": "ভ্লাদ চিরিচেস"
            },
            "shortNameTranslation": {
              "ar": "ف. كيريكيش",
              "hi": "वी. चिरिचेस",
              "bn": "ভি. চিরিচেস"
            }
          }
        },
        "teamId": 3301,
        "shirtNumber": 21,
        "jerseyNumber": "21",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Emanuele Giaccherini",
          "slug": "emanuele-giaccherini",
          "shortName": "E. Giaccherini",
          "position": "M",
          "height": 167,
          "userCount": 108,
          "id": 76133,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 484099200,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "إيمانويل جياكيريني",
              "hi": "इमानुएल गियाचेरिनी",
              "bn": "ইমানুয়েল গিয়াচেরিনি"
            },
            "shortNameTranslation": {
              "ar": "إ. جياكيريني",
              "hi": "ई. गियाचेरिनी",
              "bn": "ই. গিয়াচেরিনি"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 4,
        "jerseyNumber": "4",
        "position": "M",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Manolo Gabbiadini",
          "slug": "manolo-gabbiadini",
          "shortName": "M. Gabbiadini",
          "position": "F",
          "jerseyNumber": "11",
          "height": 186,
          "userCount": 446,
          "gender": "M",
          "id": 53816,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 691113600,
          "proposedMarketValueRaw": {
            "value": 650000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "مانولو غابياديني",
              "hi": "मनोलो गाबियादिनी",
              "bn": "মানোলো গ্যাবিয়াদিনি"
            },
            "shortNameTranslation": {
              "ar": "م. غابياديني",
              "hi": "एम. गाबियादिनी",
              "bn": "এম. গ্যাবিয়াদিনি"
            }
          }
        },
        "teamId": 46332,
        "shirtNumber": 23,
        "jerseyNumber": "23",
        "position": "F",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      }
    ],
    "supportStaff": [],
    "formation": "4-3-3",
    "playerColor": {
      "primary": "ffffff",
      "number": "008cea",
      "outline": "ffffff",
      "fancyNumber": "222226"
    },
    "goalkeeperColor": {
      "primary": "ff9900",
      "number": "000000",
      "outline": "ff9900",
      "fancyNumber": "222226"
    }
  },
  "away": {
    "players": [
      {
        "player": {
          "name": "Josip Posavec",
          "slug": "josip-posavec",
          "shortName": "J. Posavec",
          "position": "G",
          "jerseyNumber": "12",
          "height": 190,
          "userCount": 244,
          "gender": "M",
          "id": 798416,
          "country": {
            "alpha2": "HR",
            "alpha3": "HRV",
            "name": "Croatia",
            "slug": "croatia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 826416000,
          "proposedMarketValueRaw": {
            "value": 310000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "جوزيب بوسافيك",
              "hi": "जोसिप पोसावेक",
              "bn": "জোসিপ পোসাভেক"
            },
            "shortNameTranslation": {
              "ar": "ج. بوسافيك",
              "hi": "जे. पोसावेक",
              "bn": "জে. পোসাভেক"
            }
          }
        },
        "teamId": 36246,
        "shirtNumber": 1,
        "jerseyNumber": "1",
        "position": "G",
        "substitute": false,
        "statistics": {
          "totalPass": 31,
          "accuratePass": 14,
          "totalLongBalls": 28,
          "accurateLongBalls": 11,
          "accurateOwnHalfPasses": 8,
          "totalOwnHalfPasses": 10,
          "accurateOppositionHalfPasses": 6,
          "totalOppositionHalfPasses": 21,
          "totalClearance": 5,
          "ballRecovery": 11,
          "errorLeadToAGoal": 1,
          "goodHighClaim": 1,
          "savedShotsFromInsideTheBox": 6,
          "saves": 9,
          "punches": 2,
          "totalKeeperSweeper": 2,
          "accurateKeeperSweeper": 2,
          "minutesPlayed": 90,
          "touches": 48,
          "rating": 8.2,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Andrea Rispoli",
          "slug": "andrea-rispoli",
          "shortName": "A. Rispoli",
          "position": "D",
          "height": 188,
          "userCount": 30,
          "gender": "M",
          "id": 43910,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 591494400,
          "proposedMarketValueRaw": {
            "value": 105000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "أندريا ريسبولي",
              "hi": "एंड्रिया रिस्पोली",
              "bn": "আন্দ্রেয়া রিসপোলি"
            },
            "shortNameTranslation": {
              "ar": "أ. ريسبولي",
              "hi": "ए. रिस्पोली",
              "bn": "এ. রিসপোলি"
            }
          }
        },
        "teamId": 2718,
        "shirtNumber": 3,
        "jerseyNumber": "3",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 31,
          "accuratePass": 18,
          "totalLongBalls": 10,
          "accurateLongBalls": 3,
          "goalAssist": 1,
          "accurateOwnHalfPasses": 11,
          "totalOwnHalfPasses": 17,
          "accurateOppositionHalfPasses": 9,
          "totalOppositionHalfPasses": 16,
          "totalCross": 2,
          "accurateCross": 2,
          "duelLost": 4,
          "challengeLost": 2,
          "totalContest": 1,
          "totalClearance": 9,
          "interceptionWon": 5,
          "ballRecovery": 9,
          "fouls": 1,
          "minutesPlayed": 90,
          "touches": 61,
          "rating": 7.1,
          "keyPass": 1,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Edoardo Goldaniga",
          "firstName": "",
          "lastName": "",
          "slug": "edoardo-goldaniga",
          "shortName": "E. Goldaniga",
          "position": "D",
          "jerseyNumber": "5",
          "height": 188,
          "userCount": 303,
          "gender": "M",
          "id": 295133,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 752198400,
          "proposedMarketValueRaw": {
            "value": 2900000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "إدواردو جولدانيجا",
              "hi": "एडोआर्डो गोल्डानिगा",
              "bn": "এদোয়ার্দো গোল্ডানিগা"
            },
            "shortNameTranslation": {
              "ar": "إ. جولدانيجا",
              "hi": "ई. गोल्डानिगा",
              "bn": "ই. গোল্ডানিগা"
            }
          }
        },
        "teamId": 2704,
        "shirtNumber": 6,
        "jerseyNumber": "6",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 15,
          "accuratePass": 10,
          "totalLongBalls": 6,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 6,
          "totalOwnHalfPasses": 8,
          "accurateOppositionHalfPasses": 4,
          "totalOppositionHalfPasses": 7,
          "aerialWon": 1,
          "duelLost": 3,
          "duelWon": 3,
          "dispossessed": 1,
          "totalClearance": 8,
          "outfielderBlock": 1,
          "interceptionWon": 2,
          "ballRecovery": 3,
          "totalTackle": 2,
          "fouls": 2,
          "minutesPlayed": 89,
          "touches": 31,
          "rating": 6.6,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Giancarlo González",
          "slug": "giancarlo-gonzalez",
          "shortName": "G. González",
          "position": "D",
          "jerseyNumber": "26",
          "height": 186,
          "userCount": 100,
          "gender": "M",
          "id": 117793,
          "country": {
            "alpha2": "CR",
            "alpha3": "CRI",
            "name": "Costa Rica",
            "slug": "costa-rica"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 571276800,
          "proposedMarketValueRaw": {
            "value": 52000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "جيانكارلو غونزاليس",
              "hi": "जियानकार्लो गोंजालेज",
              "bn": "জিয়ানকার্লো গনজালেজ"
            },
            "shortNameTranslation": {
              "ar": "ج. غونزاليس",
              "hi": "जी. गोंजालेज",
              "bn": "জি. গনজালেজ"
            }
          }
        },
        "teamId": 262906,
        "shirtNumber": 12,
        "jerseyNumber": "12",
        "position": "D",
        "substitute": false,
        "captain": true,
        "statistics": {
          "totalPass": 21,
          "accuratePass": 14,
          "totalLongBalls": 8,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 10,
          "totalOwnHalfPasses": 13,
          "accurateOppositionHalfPasses": 4,
          "totalOppositionHalfPasses": 8,
          "aerialLost": 1,
          "duelLost": 1,
          "duelWon": 1,
          "totalClearance": 13,
          "outfielderBlock": 2,
          "interceptionWon": 4,
          "ballRecovery": 2,
          "totalTackle": 1,
          "minutesPlayed": 90,
          "touches": 44,
          "rating": 7.6,
          "keyPass": 1,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Giuseppe Pezzella",
          "slug": "giuseppe-pezzella",
          "shortName": "G. Pezzella",
          "position": "M",
          "jerseyNumber": "3",
          "height": 187,
          "userCount": 607,
          "gender": "M",
          "id": 814589,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 880761600,
          "proposedMarketValueRaw": {
            "value": 2600000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "جوزيبي بيزيلا",
              "hi": "गिउसेप्पे पेज़ेला",
              "bn": "জিউসেপ পেজেলা"
            },
            "shortNameTranslation": {
              "ar": "ج. بيزيلا",
              "hi": "जी. पेज़ेला",
              "bn": "জি. পেজেলা"
            }
          }
        },
        "teamId": 2761,
        "shirtNumber": 97,
        "jerseyNumber": "97",
        "position": "D",
        "substitute": false,
        "statistics": {
          "totalPass": 20,
          "accuratePass": 12,
          "totalLongBalls": 2,
          "accurateOwnHalfPasses": 8,
          "totalOwnHalfPasses": 15,
          "accurateOppositionHalfPasses": 4,
          "totalOppositionHalfPasses": 7,
          "totalCross": 2,
          "duelLost": 3,
          "duelWon": 5,
          "challengeLost": 1,
          "dispossessed": 1,
          "totalContest": 5,
          "wonContest": 4,
          "totalClearance": 10,
          "interceptionWon": 2,
          "ballRecovery": 5,
          "unsuccessfulTouch": 2,
          "wasFouled": 1,
          "minutesPlayed": 90,
          "touches": 59,
          "rating": 6.7,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Bruno Henrique",
          "slug": "bruno-henrique",
          "shortName": "B. Henrique",
          "position": "M",
          "jerseyNumber": "8",
          "height": 180,
          "userCount": 911,
          "gender": "M",
          "id": 345113,
          "country": {
            "alpha2": "BR",
            "alpha3": "BRA",
            "name": "Brazil",
            "slug": "brazil"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 624931200,
          "proposedMarketValueRaw": {
            "value": 440000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "برونو هنريكي",
              "hi": "ब्रूनो हेनरिक",
              "bn": "ব্রুনো হেনরিক"
            },
            "shortNameTranslation": {
              "ar": "ب. هنريكي",
              "hi": "ब्रूनो हेनरिक",
              "bn": "ব্রুনো হেনরিক"
            }
          }
        },
        "teamId": 1966,
        "shirtNumber": 25,
        "jerseyNumber": "25",
        "position": "M",
        "substitute": false,
        "statistics": {
          "totalPass": 29,
          "accuratePass": 24,
          "totalLongBalls": 5,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 14,
          "totalOwnHalfPasses": 16,
          "accurateOppositionHalfPasses": 10,
          "totalOppositionHalfPasses": 14,
          "totalCross": 1,
          "aerialLost": 1,
          "duelLost": 7,
          "duelWon": 3,
          "challengeLost": 2,
          "dispossessed": 1,
          "totalContest": 2,
          "wonContest": 1,
          "totalClearance": 1,
          "interceptionWon": 3,
          "ballRecovery": 8,
          "totalTackle": 1,
          "unsuccessfulTouch": 4,
          "wasFouled": 1,
          "fouls": 2,
          "minutesPlayed": 90,
          "touches": 46,
          "rating": 6.6,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Mato Jajalo",
          "slug": "mato-jajalo",
          "shortName": "M. Jajalo",
          "position": "M",
          "height": 180,
          "userCount": 71,
          "gender": "M",
          "id": 42147,
          "country": {
            "alpha2": "BA",
            "alpha3": "BIH",
            "name": "Bosnia & Herzegovina",
            "slug": "bosnia-and-herzegovina"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 580521600,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ماتو جاجالو",
              "hi": "माटो जाजालो",
              "bn": "মাতো জাজলো"
            },
            "shortNameTranslation": {
              "ar": "م. جاجالو",
              "hi": "एम. जाजालो",
              "bn": "এম. জাজালো"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 28,
        "jerseyNumber": "28",
        "position": "M",
        "substitute": false,
        "statistics": {
          "totalPass": 31,
          "accuratePass": 25,
          "totalLongBalls": 5,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 20,
          "totalOwnHalfPasses": 24,
          "accurateOppositionHalfPasses": 5,
          "totalOppositionHalfPasses": 7,
          "duelLost": 4,
          "duelWon": 4,
          "challengeLost": 1,
          "dispossessed": 1,
          "totalContest": 1,
          "wonContest": 1,
          "totalClearance": 4,
          "outfielderBlock": 1,
          "interceptionWon": 1,
          "ballRecovery": 9,
          "totalTackle": 2,
          "unsuccessfulTouch": 3,
          "wasFouled": 1,
          "fouls": 2,
          "minutesPlayed": 90,
          "touches": 48,
          "rating": 6.3,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Ivajlo Chochev",
          "firstName": "",
          "lastName": "",
          "slug": "ivajlo-chochev",
          "shortName": "I. Chochev",
          "position": "M",
          "jerseyNumber": "18",
          "height": 187,
          "userCount": 250,
          "gender": "M",
          "id": 94681,
          "country": {
            "alpha2": "BG",
            "alpha3": "BGR",
            "name": "Bulgaria",
            "slug": "bulgaria"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 729993600,
          "proposedMarketValueRaw": {
            "value": 2200000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "إيفايلو تشوتشيف",
              "hi": "इवाज्लो चोचेव",
              "bn": "ইভাজলো চোচেভ"
            },
            "shortNameTranslation": {
              "ar": "إ. تشوتشيف",
              "hi": "आई. चोचेव",
              "bn": "আই. চোচেভ"
            }
          }
        },
        "teamId": 43840,
        "shirtNumber": 18,
        "jerseyNumber": "18",
        "position": "M",
        "substitute": false,
        "statistics": {
          "totalPass": 25,
          "accuratePass": 18,
          "totalLongBalls": 4,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 11,
          "totalOwnHalfPasses": 13,
          "accurateOppositionHalfPasses": 7,
          "totalOppositionHalfPasses": 12,
          "aerialWon": 3,
          "duelLost": 3,
          "duelWon": 10,
          "totalContest": 6,
          "wonContest": 3,
          "onTargetScoringAttempt": 1,
          "blockedScoringAttempt": 1,
          "totalClearance": 3,
          "interceptionWon": 6,
          "ballRecovery": 6,
          "totalTackle": 4,
          "minutesPlayed": 90,
          "touches": 47,
          "rating": 7.8,
          "totalShots": 2,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Aleksandar Trajkovski",
          "slug": "aleksandar-trajkovski",
          "shortName": "A. Trajkovski",
          "position": "F",
          "height": 179,
          "userCount": 777,
          "gender": "M",
          "id": 60184,
          "country": {
            "alpha2": "MK",
            "alpha3": "MKD",
            "name": "North Macedonia",
            "slug": "north-macedonia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 715651200,
          "proposedMarketValueRaw": {
            "value": 170000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ألكسندر ترايكوفسكي",
              "hi": "एलेक्ज़ेंडर ट्रैजकोव्स्की",
              "bn": "আলেকসান্ডার ট্রাজকোভস্কি"
            },
            "shortNameTranslation": {
              "ar": "أ. ترايكوفسكي",
              "hi": "ए. ट्रैजकोव्स्की",
              "bn": "এ. ট্রাজকোভস্কি"
            }
          }
        },
        "teamId": 4777,
        "shirtNumber": 8,
        "jerseyNumber": "8",
        "position": "F",
        "substitute": false,
        "statistics": {
          "totalPass": 16,
          "accuratePass": 12,
          "accurateOwnHalfPasses": 8,
          "totalOwnHalfPasses": 10,
          "accurateOppositionHalfPasses": 4,
          "totalOppositionHalfPasses": 7,
          "totalCross": 1,
          "aerialLost": 1,
          "duelLost": 3,
          "duelWon": 2,
          "challengeLost": 1,
          "dispossessed": 1,
          "interceptionWon": 1,
          "ballRecovery": 7,
          "totalTackle": 1,
          "unsuccessfulTouch": 5,
          "wasFouled": 1,
          "totalOffside": 1,
          "minutesPlayed": 60,
          "touches": 28,
          "rating": 6.4,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Ilija Nestorovski",
          "slug": "ilija-nestorovski",
          "shortName": "I. Nestorovski",
          "position": "F",
          "jerseyNumber": "90",
          "height": 182,
          "userCount": 220,
          "gender": "M",
          "id": 55314,
          "country": {
            "alpha2": "MK",
            "alpha3": "MKD",
            "name": "North Macedonia",
            "slug": "north-macedonia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 637200000,
          "proposedMarketValueRaw": {
            "value": 190000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "إيليا نيستوروفسكي",
              "hi": "इलिया नेस्टरोव्स्की",
              "bn": "ইলিজা নেস্টোরভস্কি"
            },
            "shortNameTranslation": {
              "ar": "إ. نيستوروفسكي",
              "hi": "आई. नेस्टरोव्स्की",
              "bn": "আই. নেস্টোরভস্কি"
            }
          }
        },
        "teamId": 2042,
        "shirtNumber": 30,
        "jerseyNumber": "30",
        "position": "F",
        "substitute": false,
        "statistics": {
          "totalPass": 25,
          "accuratePass": 15,
          "totalLongBalls": 5,
          "accurateLongBalls": 2,
          "accurateOwnHalfPasses": 12,
          "totalOwnHalfPasses": 14,
          "accurateOppositionHalfPasses": 3,
          "totalOppositionHalfPasses": 11,
          "aerialLost": 7,
          "aerialWon": 4,
          "duelLost": 12,
          "duelWon": 6,
          "challengeLost": 1,
          "dispossessed": 2,
          "totalContest": 2,
          "wonContest": 1,
          "onTargetScoringAttempt": 1,
          "goals": 1,
          "totalClearance": 2,
          "interceptionWon": 1,
          "ballRecovery": 2,
          "unsuccessfulTouch": 4,
          "wasFouled": 1,
          "fouls": 1,
          "totalOffside": 2,
          "minutesPlayed": 89,
          "touches": 41,
          "rating": 7.2,
          "keyPass": 1,
          "totalShots": 1,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Robin Quaison",
          "firstName": "",
          "lastName": "",
          "slug": "robin-quaison",
          "shortName": "R. Quaison",
          "position": "F",
          "height": 184,
          "userCount": 393,
          "gender": "M",
          "id": 163565,
          "country": {
            "alpha2": "SE",
            "alpha3": "SWE",
            "name": "Sweden",
            "slug": "sweden"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 750124800,
          "proposedMarketValueRaw": {
            "value": 825000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "روبن كوايسون",
              "hi": "रॉबिन क्वाइसन",
              "bn": "রবিন কোয়েসন"
            },
            "shortNameTranslation": {
              "ar": "ر. كوايسون",
              "hi": "आर. क्वाइसन",
              "bn": "আর. কোয়েসন"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 21,
        "jerseyNumber": "21",
        "position": "F",
        "substitute": false,
        "statistics": {
          "totalPass": 18,
          "accuratePass": 13,
          "accurateOwnHalfPasses": 9,
          "totalOwnHalfPasses": 12,
          "accurateOppositionHalfPasses": 4,
          "totalOppositionHalfPasses": 6,
          "aerialLost": 5,
          "aerialWon": 1,
          "duelLost": 8,
          "duelWon": 1,
          "dispossessed": 1,
          "totalContest": 1,
          "totalClearance": 4,
          "interceptionWon": 2,
          "ballRecovery": 2,
          "unsuccessfulTouch": 1,
          "fouls": 1,
          "totalOffside": 1,
          "minutesPlayed": 88,
          "touches": 31,
          "rating": 6.3,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Haitam Aleesami",
          "slug": "haitam-aleesami",
          "shortName": "H. Aleesami",
          "position": "D",
          "jerseyNumber": "5",
          "height": 181,
          "userCount": 84,
          "gender": "M",
          "id": 250319,
          "country": {
            "alpha2": "NO",
            "alpha3": "NOR",
            "name": "Norway",
            "slug": "norway"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 680918400,
          "proposedMarketValueRaw": {
            "value": 230000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "هيثم العيسمي",
              "hi": "हेतम अलेसामी",
              "bn": "হাইথাম আলিসামি"
            },
            "shortNameTranslation": {
              "ar": "ه. العيسمي",
              "hi": "एच. अलेसामी",
              "bn": "এইচ. আলিসামি"
            }
          }
        },
        "teamId": 656,
        "shirtNumber": 19,
        "jerseyNumber": "19",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalPass": 5,
          "accuratePass": 2,
          "accurateOwnHalfPasses": 2,
          "totalOwnHalfPasses": 2,
          "totalOppositionHalfPasses": 3,
          "aerialLost": 1,
          "duelLost": 4,
          "challengeLost": 1,
          "dispossessed": 1,
          "totalContest": 1,
          "totalClearance": 1,
          "ballRecovery": 2,
          "minutesPlayed": 30,
          "touches": 9,
          "rating": 6,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Michel Morganella",
          "slug": "michel-morganella",
          "shortName": "M. Morganella",
          "position": "D",
          "jerseyNumber": "4",
          "height": 184,
          "userCount": 10,
          "gender": "M",
          "id": 21942,
          "country": {
            "alpha2": "CH",
            "alpha3": "CHE",
            "name": "Switzerland",
            "slug": "switzerland"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 611366400,
          "proposedMarketValueRaw": {
            "value": 27000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ميشيل مورغانلا"
            },
            "shortNameTranslation": {
              "ar": "م. مورغانلا"
            }
          }
        },
        "teamId": 2495,
        "shirtNumber": 89,
        "jerseyNumber": "89",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalPass": 1,
          "accuratePass": 1,
          "accurateOwnHalfPasses": 1,
          "totalOwnHalfPasses": 1,
          "minutesPlayed": 2,
          "touches": 1,
          "rating": 6.5,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Toni Šunjić",
          "slug": "toni-sunjic",
          "shortName": "T. Šunjić",
          "position": "D",
          "height": 192,
          "userCount": 94,
          "gender": "M",
          "id": 55392,
          "country": {
            "alpha2": "BA",
            "alpha3": "BIH",
            "name": "Bosnia & Herzegovina",
            "slug": "bosnia-and-herzegovina"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 598147200,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "توني سونيجيتش",
              "hi": "टोनी सुंजिक",
              "bn": "টনি সুনজিচ"
            },
            "shortNameTranslation": {
              "ar": "ت. سونيجيتش",
              "hi": "टी. सुंजिक",
              "bn": "টি. সুনজিচ"
            }
          }
        },
        "teamId": 485806,
        "shirtNumber": 44,
        "jerseyNumber": "44",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalPass": 1,
          "accuratePass": 1,
          "accurateOwnHalfPasses": 1,
          "totalOwnHalfPasses": 1,
          "minutesPlayed": 1,
          "touches": 1,
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Leonardo Marson",
          "firstName": "",
          "lastName": "",
          "slug": "leonardo-marson",
          "shortName": "L. Marson",
          "position": "G",
          "jerseyNumber": "77",
          "height": 194,
          "userCount": 33,
          "gender": "M",
          "id": 816837,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 883958400,
          "proposedMarketValueRaw": {
            "value": 245000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ليوناردو مارسون"
            },
            "shortNameTranslation": {
              "ar": "ل. مارسون"
            }
          }
        },
        "teamId": 2730,
        "shirtNumber": 55,
        "jerseyNumber": "55",
        "position": "G",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Samuele Guddo",
          "slug": "samuele-guddo",
          "shortName": "S. Guddo",
          "position": "G",
          "jerseyNumber": "25",
          "height": 188,
          "userCount": 7,
          "gender": "M",
          "id": 859459,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 930873600,
          "proposedMarketValueRaw": {
            "value": 105000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "سامويل غودو"
            },
            "shortNameTranslation": {
              "ar": "س. غودو"
            }
          }
        },
        "teamId": 518330,
        "shirtNumber": 57,
        "jerseyNumber": "57",
        "position": "G",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Thiago Cionek",
          "slug": "thiago-cionek",
          "shortName": "T. Cionek",
          "position": "D",
          "jerseyNumber": "15",
          "height": 184,
          "userCount": 483,
          "gender": "M",
          "id": 44932,
          "country": {
            "alpha2": "PL",
            "alpha3": "POL",
            "name": "Poland",
            "slug": "poland"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 514425600,
          "proposedMarketValueRaw": {
            "value": 23000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "تياغو كيونيك",
              "hi": "थियागो सियोनेक",
              "bn": "থিয়াগো সিওনেক"
            },
            "shortNameTranslation": {
              "ar": "ت. كيونيك",
              "hi": "टी. सियोनेक",
              "bn": "টি. সিওনেক"
            }
          }
        },
        "teamId": 6361,
        "shirtNumber": 15,
        "jerseyNumber": "15",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Siniša Anđelković",
          "firstName": "",
          "lastName": "",
          "slug": "sinisa-andelkovic",
          "shortName": "S. Anđelković",
          "position": "D",
          "height": 186,
          "userCount": 10,
          "id": 49038,
          "country": {
            "alpha2": "SI",
            "alpha3": "SVN",
            "name": "Slovenia",
            "slug": "slovenia"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 508636800,
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "سينيشا أنجيليكوفيتش",
              "hi": "सिनिसा एन्डेलकोविच",
              "bn": "সিনিশা আনদেলকোভিচ"
            },
            "shortNameTranslation": {
              "ar": "س. أنجيليكوفيتش",
              "hi": "एस. एन्डेलकोविच",
              "bn": "এস. আনদেলকোভিচ"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 4,
        "jerseyNumber": "4",
        "position": "D",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Norbert Balogh",
          "firstName": "",
          "lastName": "",
          "slug": "norbert-balogh",
          "shortName": "N. Balogh",
          "position": "F",
          "height": 197,
          "userCount": 154,
          "gender": "M",
          "id": 576324,
          "country": {
            "alpha2": "HU",
            "alpha3": "HUN",
            "name": "Hungary",
            "slug": "hungary"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 824860800,
          "proposedMarketValueRaw": {
            "value": 49000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "نوربرت بالوغ",
              "hi": "नॉर्बर्ट बालोग",
              "bn": "নরবার্ট বালোগ"
            },
            "shortNameTranslation": {
              "ar": "ن. بالوغ",
              "hi": "एन. बालोग",
              "bn": "এন. বালোগ"
            }
          }
        },
        "teamId": 241802,
        "shirtNumber": 22,
        "jerseyNumber": "22",
        "position": "F",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Simone Lo Faso",
          "slug": "lo-faso-simone",
          "shortName": "S. L. Faso",
          "position": "F",
          "jerseyNumber": "77",
          "height": 180,
          "userCount": 93,
          "gender": "M",
          "id": 795054,
          "country": {
            "alpha2": "IT",
            "alpha3": "ITA",
            "name": "Italy",
            "slug": "italy"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 887673600,
          "proposedMarketValueRaw": {
            "value": 47000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "سيموني لو فاسو"
            },
            "shortNameTranslation": {
              "ar": "س. فاسو",
              "hi": "एस. एल. फासो",
              "bn": "এস. এল. ফাসো"
            }
          }
        },
        "teamId": 940319,
        "shirtNumber": 98,
        "jerseyNumber": "98",
        "position": "F",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Carlos Embalo",
          "slug": "carlos-embalo",
          "shortName": "C. Embalo",
          "position": "F",
          "jerseyNumber": "77",
          "height": 177,
          "userCount": 27,
          "gender": "M",
          "id": 345849,
          "country": {
            "alpha2": "GW",
            "alpha3": "GNB",
            "name": "Guinea-Bissau",
            "slug": "guinea-bissau"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 785721600,
          "proposedMarketValueRaw": {
            "value": 48000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "كارلوس إمبالو"
            },
            "shortNameTranslation": {
              "ar": "ك. إمبالو"
            }
          }
        },
        "teamId": 170612,
        "shirtNumber": 11,
        "jerseyNumber": "11",
        "position": "M",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      },
      {
        "player": {
          "name": "Stefan Silva",
          "slug": "stefan-silva",
          "shortName": "S. Silva",
          "position": "F",
          "jerseyNumber": "10",
          "height": 182,
          "userCount": 17,
          "gender": "M",
          "id": 384166,
          "country": {
            "alpha2": "SE",
            "alpha3": "SWE",
            "name": "Sweden",
            "slug": "sweden"
          },
          "marketValueCurrency": "EUR",
          "dateOfBirthTimestamp": 637027200,
          "proposedMarketValueRaw": {
            "value": 96000,
            "currency": "EUR"
          },
          "fieldTranslations": {
            "nameTranslation": {
              "ar": "ستيفان سيلفا"
            },
            "shortNameTranslation": {
              "ar": "س. سيلفا"
            }
          }
        },
        "teamId": 1074854,
        "shirtNumber": 9,
        "jerseyNumber": "9",
        "position": "F",
        "substitute": true,
        "statistics": {
          "totalShots": 0,
          "statisticsType": {
            "sportSlug": "football",
            "statisticsType": "player"
          }
        }
      }
    ],
    "supportStaff": [],
    "formation": "4-3-3",
    "playerColor": {
      "primary": "ffb3d9",
      "number": "000000",
      "outline": "ffb3d9",
      "fancyNumber": "222226"
    },
    "goalkeeperColor": {
      "primary": "ff99ff",
      "number": "000000",
      "outline": "ff99ff",
      "fancyNumber": "222226"
    }
  },
  "statisticalVersion": null
}</pre><div class="json-formatter-container"></div>`

	incidentsJSON = `xxxxx{
  "incidents": [
    {
      "text": "FT",
      "homeScore": 1,
      "awayScore": 1,
      "isLive": false,
      "time": 90,
      "addedTime": 999,
      "timeSeconds": 5400,
      "reversedPeriodTime": 1,
      "reversedPeriodTimeSeconds": 0,
      "periodTimeSeconds": 2700,
      "incidentType": "period"
    },
    {
      "player": {
        "name": "Giancarlo González",
        "slug": "giancarlo-gonzalez",
        "shortName": "G. González",
        "position": "D",
        "jerseyNumber": "26",
        "height": 186,
        "userCount": 100,
        "gender": "M",
        "id": 117793,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 571276800,
        "proposedMarketValueRaw": {
          "value": 52000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "جيانكارلو غونزاليس",
            "hi": "जियानकार्लो गोंजालेज",
            "bn": "জিয়ানকার্লো গনজালেজ"
          },
          "shortNameTranslation": {
            "ar": "ج. غونزاليس",
            "hi": "जी. गोंजालेज",
            "bn": "জি. গনজালেজ"
          }
        }
      },
      "reason": "Argument",
      "rescinded": false,
      "id": 118585170,
      "time": 90,
      "addedTime": 2,
      "isHome": false,
      "incidentClass": "yellow",
      "reversedPeriodTime": 1,
      "incidentType": "card"
    },
    {
      "playerIn": {
        "name": "Toni Šunjić",
        "slug": "toni-sunjic",
        "shortName": "T. Šunjić",
        "position": "D",
        "height": 192,
        "userCount": 94,
        "gender": "M",
        "id": 55392,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 598147200,
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "توني سونيجيتش",
            "hi": "टोनी सुंजिक",
            "bn": "টনি সুনজিচ"
          },
          "shortNameTranslation": {
            "ar": "ت. سونيجيتش",
            "hi": "टी. सुंजिक",
            "bn": "টি. সুনজিচ"
          }
        }
      },
      "playerOut": {
        "name": "Ilija Nestorovski",
        "slug": "ilija-nestorovski",
        "shortName": "I. Nestorovski",
        "position": "F",
        "jerseyNumber": "90",
        "height": 182,
        "userCount": 220,
        "gender": "M",
        "id": 55314,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 637200000,
        "proposedMarketValueRaw": {
          "value": 190000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "إيليا نيستوروفسكي",
            "hi": "इलिया नेस्टरोव्स्की",
            "bn": "ইলিজা নেস্টোরভস্কি"
          },
          "shortNameTranslation": {
            "ar": "إ. نيستوروفسكي",
            "hi": "आई. नेस्टरोव्स्की",
            "bn": "আই. নেস্টোরভস্কি"
          }
        }
      },
      "id": 118515522,
      "time": 90,
      "addedTime": 1,
      "injury": false,
      "isHome": false,
      "incidentClass": "regular",
      "reversedPeriodTime": 1,
      "incidentType": "substitution"
    },
    {
      "length": 5,
      "time": 90,
      "addedTime": 0,
      "reversedPeriodTime": 1,
      "incidentType": "injuryTime"
    },
    {
      "player": {
        "name": "Edoardo Goldaniga",
        "firstName": "",
        "lastName": "",
        "slug": "edoardo-goldaniga",
        "shortName": "E. Goldaniga",
        "position": "D",
        "jerseyNumber": "5",
        "height": 188,
        "userCount": 303,
        "gender": "M",
        "id": 295133,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 752198400,
        "proposedMarketValueRaw": {
          "value": 2900000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "إدواردو جولدانيجا",
            "hi": "एडोआर्डो गोल्डानिगा",
            "bn": "এদোয়ার্দো গোল্ডানিগা"
          },
          "shortNameTranslation": {
            "ar": "إ. جولدانيجا",
            "hi": "ई. गोल्डानिगा",
            "bn": "ই. গোল্ডানিগা"
          }
        }
      },
      "reason": "Foul",
      "rescinded": false,
      "id": 118585165,
      "time": 90,
      "isHome": false,
      "incidentClass": "red",
      "reversedPeriodTime": 1,
      "incidentType": "card"
    },
    {
      "playerIn": {
        "name": "Michel Morganella",
        "slug": "michel-morganella",
        "shortName": "M. Morganella",
        "position": "D",
        "jerseyNumber": "4",
        "height": 184,
        "userCount": 10,
        "gender": "M",
        "id": 21942,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 611366400,
        "proposedMarketValueRaw": {
          "value": 27000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "ميشيل مورغانلا"
          },
          "shortNameTranslation": {
            "ar": "م. مورغانلا"
          }
        }
      },
      "playerOut": {
        "name": "Robin Quaison",
        "firstName": "",
        "lastName": "",
        "slug": "robin-quaison",
        "shortName": "R. Quaison",
        "position": "F",
        "height": 184,
        "userCount": 393,
        "gender": "M",
        "id": 163565,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 750124800,
        "proposedMarketValueRaw": {
          "value": 825000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "روبن كوايسون",
            "hi": "रॉबिन क्वाइसन",
            "bn": "রবিন কোয়েসন"
          },
          "shortNameTranslation": {
            "ar": "ر. كوايسون",
            "hi": "आर. क्वाइसन",
            "bn": "আর. কোয়েসন"
          }
        }
      },
      "id": 118515520,
      "time": 88,
      "injury": false,
      "isHome": false,
      "incidentClass": "regular",
      "reversedPeriodTime": 3,
      "incidentType": "substitution"
    },
    {
      "homeScore": 1,
      "awayScore": 1,
      "player": {
        "name": "Dries Mertens",
        "slug": "dries-mertens",
        "shortName": "D. Mertens",
        "position": "M",
        "height": 169,
        "userCount": 11877,
        "gender": "M",
        "sofascoreId": "driesmertens",
        "id": 32493,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 547257600,
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "دريس ميرتينز",
            "hi": "ड्रीस मेर्टेंस",
            "bn": "ড্রাইস মার্টেনস"
          },
          "shortNameTranslation": {
            "ar": "د. ميرتنز",
            "hi": "डी. मेर्टेंस",
            "bn": "ডি. মার্টেনস"
          }
        }
      },
      "assist1": {
        "name": "Piotr Zieliński",
        "slug": "piotr-zielinski",
        "shortName": "P. Zieliński",
        "position": "M",
        "jerseyNumber": "7",
        "height": 180,
        "userCount": 18007,
        "gender": "M",
        "sofascoreId": "zielu_94",
        "id": 138605,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 769392000,
        "proposedMarketValueRaw": {
          "value": 10300000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "بيوتر زيلينسكي",
            "hi": "पिओट्र ज़िलिंस्की",
            "bn": "পিওতর জিলিনিস্কি"
          },
          "shortNameTranslation": {
            "ar": "ب. زيلينسكي",
            "hi": "पी. ज़िलिंस्की",
            "bn": "পি. জিলিনিস্কি"
          }
        }
      },
      "footballPassingNetworkAction": [],
      "id": 120826224,
      "time": 66,
      "isHome": true,
      "incidentClass": "regular",
      "reversedPeriodTime": 25,
      "incidentType": "goal"
    },
    {
      "player": {
        "name": "Mato Jajalo",
        "slug": "mato-jajalo",
        "shortName": "M. Jajalo",
        "position": "M",
        "height": 180,
        "userCount": 71,
        "gender": "M",
        "id": 42147,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 580521600,
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "ماتو جاجالو",
            "hi": "माटो जाजालो",
            "bn": "মাতো জাজলো"
          },
          "shortNameTranslation": {
            "ar": "م. جاجالو",
            "hi": "एम. जाजालो",
            "bn": "এম. জাজালো"
          }
        }
      },
      "reason": "Foul",
      "rescinded": false,
      "id": 118585163,
      "time": 65,
      "isHome": false,
      "incidentClass": "yellow",
      "reversedPeriodTime": 26,
      "incidentType": "card"
    },
    {
      "playerIn": {
        "name": "Leonardo Pavoletti",
        "slug": "leonardo-pavoletti",
        "shortName": "L. Pavoletti",
        "position": "F",
        "jerseyNumber": "30",
        "height": 188,
        "userCount": 523,
        "gender": "M",
        "id": 132704,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 596505600,
        "proposedMarketValueRaw": {
          "value": 640000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "ليوناردو بافوليتي",
            "hi": "लियोनार्डो पावोलेटी",
            "bn": "লিওনার্দো পাভোলেটি"
          },
          "shortNameTranslation": {
            "ar": "ل. بافوليتي",
            "hi": "एल. पावोलेटी",
            "bn": "এল. পাভোলেটি"
          }
        }
      },
      "playerOut": {
        "name": "Jorginho",
        "slug": "jorginho",
        "shortName": "Jorginho",
        "position": "M",
        "jerseyNumber": "21",
        "height": 180,
        "userCount": 30309,
        "gender": "M",
        "sofascoreId": "jorginhofrello",
        "id": 132874,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 693187200,
        "proposedMarketValueRaw": {
          "value": 6500000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "جورجينيو",
            "hi": "जोर्गिनहो",
            "bn": "জর্গিনহো"
          },
          "shortNameTranslation": {
            "ar": "جورجينيو",
            "hi": "जोर्गिनहो",
            "bn": "জর্গিনহো"
          }
        }
      },
      "id": 118515519,
      "time": 63,
      "injury": false,
      "isHome": true,
      "incidentClass": "regular",
      "reversedPeriodTime": 28,
      "incidentType": "substitution"
    },
    {
      "playerIn": {
        "name": "Haitam Aleesami",
        "slug": "haitam-aleesami",
        "shortName": "H. Aleesami",
        "position": "D",
        "jerseyNumber": "5",
        "height": 181,
        "userCount": 85,
        "gender": "M",
        "id": 250319,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 680918400,
        "proposedMarketValueRaw": {
          "value": 230000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "هيثم العيسمي",
            "hi": "हेतम अलेसामी",
            "bn": "হাইথাম আলিসামি"
          },
          "shortNameTranslation": {
            "ar": "ه. العيسمي",
            "hi": "एच. अलेसामी",
            "bn": "এইচ. আলিসামি"
          }
        }
      },
      "playerOut": {
        "name": "Aleksandar Trajkovski",
        "slug": "aleksandar-trajkovski",
        "shortName": "A. Trajkovski",
        "position": "F",
        "height": 179,
        "userCount": 777,
        "gender": "M",
        "id": 60184,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 715651200,
        "proposedMarketValueRaw": {
          "value": 170000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "ألكسندر ترايكوفسكي",
            "hi": "एलेक्ज़ेंडर ट्रैजकोव्स्की",
            "bn": "আলেকসান্ডার ট্রাজকোভস্কি"
          },
          "shortNameTranslation": {
            "ar": "أ. ترايكوفسكي",
            "hi": "ए. ट्रैजकोव्स्की",
            "bn": "এ. ট্রাজকোভস্কি"
          }
        }
      },
      "id": 118515512,
      "time": 60,
      "injury": false,
      "isHome": false,
      "incidentClass": "regular",
      "reversedPeriodTime": 31,
      "incidentType": "substitution"
    },
    {
      "playerIn": {
        "name": "Piotr Zieliński",
        "slug": "piotr-zielinski",
        "shortName": "P. Zieliński",
        "position": "M",
        "jerseyNumber": "7",
        "height": 180,
        "userCount": 18007,
        "gender": "M",
        "sofascoreId": "zielu_94",
        "id": 138605,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 769392000,
        "proposedMarketValueRaw": {
          "value": 10300000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "بيوتر زيلينسكي",
            "hi": "पिओट्र ज़िलिंस्की",
            "bn": "পিওতর জিলিনিস্কি"
          },
          "shortNameTranslation": {
            "ar": "ب. زيلينسكي",
            "hi": "पी. ज़िलिंस्की",
            "bn": "পি. জিলিনিস্কি"
          }
        }
      },
      "playerOut": {
        "name": "Allan",
        "slug": "allan",
        "shortName": "Allan",
        "position": "M",
        "jerseyNumber": "25",
        "height": 175,
        "userCount": 2951,
        "gender": "M",
        "id": 158277,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 663292800,
        "proposedMarketValueRaw": {
          "value": 1900000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "ألان",
            "hi": "एलन",
            "bn": "অ্যালান"
          },
          "shortNameTranslation": {
            "ar": "ألان",
            "hi": "एलन",
            "bn": "অ্যালান"
          }
        }
      },
      "id": 118515510,
      "time": 54,
      "injury": false,
      "isHome": true,
      "incidentClass": "regular",
      "reversedPeriodTime": 37,
      "incidentType": "substitution"
    },
    {
      "text": "HT",
      "homeScore": 0,
      "awayScore": 1,
      "isLive": false,
      "time": 45,
      "addedTime": 999,
      "timeSeconds": 2700,
      "reversedPeriodTime": 1,
      "reversedPeriodTimeSeconds": 0,
      "periodTimeSeconds": 2700,
      "incidentType": "period"
    },
    {
      "length": 2,
      "time": 45,
      "addedTime": 0,
      "reversedPeriodTime": 1,
      "incidentType": "injuryTime"
    },
    {
      "player": {
        "name": "Bruno Henrique",
        "slug": "bruno-henrique",
        "shortName": "B. Henrique",
        "position": "M",
        "jerseyNumber": "8",
        "height": 180,
        "userCount": 908,
        "gender": "M",
        "id": 345113,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 624931200,
        "proposedMarketValueRaw": {
          "value": 440000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "برونو هنريكي",
            "hi": "ब्रूनो हेनरिक",
            "bn": "ব্রুনো হেনরিক"
          },
          "shortNameTranslation": {
            "ar": "ب. هنريكي",
            "hi": "ब्रूनो हेनरिक",
            "bn": "ব্রুনো হেনরিক"
          }
        }
      },
      "reason": "Foul",
      "rescinded": false,
      "id": 118585144,
      "time": 32,
      "isHome": false,
      "incidentClass": "yellow",
      "reversedPeriodTime": 14,
      "incidentType": "card"
    },
    {
      "player": {
        "name": "Robin Quaison",
        "firstName": "",
        "lastName": "",
        "slug": "robin-quaison",
        "shortName": "R. Quaison",
        "position": "F",
        "height": 184,
        "userCount": 393,
        "gender": "M",
        "id": 163565,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 750124800,
        "proposedMarketValueRaw": {
          "value": 825000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "روبن كوايسون",
            "hi": "रॉबिन क्वाइसन",
            "bn": "রবিন কোয়েসন"
          },
          "shortNameTranslation": {
            "ar": "ر. كوايسون",
            "hi": "आर. क्वाइसन",
            "bn": "আর. কোয়েসন"
          }
        }
      },
      "reason": "Foul",
      "rescinded": false,
      "id": 118585139,
      "time": 23,
      "isHome": false,
      "incidentClass": "yellow",
      "reversedPeriodTime": 23,
      "incidentType": "card"
    },
    {
      "from": "heading",
      "homeScore": 0,
      "awayScore": 1,
      "player": {
        "name": "Ilija Nestorovski",
        "slug": "ilija-nestorovski",
        "shortName": "I. Nestorovski",
        "position": "F",
        "jerseyNumber": "90",
        "height": 182,
        "userCount": 220,
        "gender": "M",
        "id": 55314,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 637200000,
        "proposedMarketValueRaw": {
          "value": 190000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "إيليا نيستوروفسكي",
            "hi": "इलिया नेस्टरोव्स्की",
            "bn": "ইলিজা নেস্টোরভস্কি"
          },
          "shortNameTranslation": {
            "ar": "إ. نيستوروفسكي",
            "hi": "आई. नेस्टरोव्स्की",
            "bn": "আই. নেস্টোরভস্কি"
          }
        }
      },
      "assist1": {
        "name": "Andrea Rispoli",
        "slug": "andrea-rispoli",
        "shortName": "A. Rispoli",
        "position": "D",
        "height": 188,
        "userCount": 30,
        "gender": "M",
        "id": 43910,
        "marketValueCurrency": "EUR",
        "dateOfBirthTimestamp": 591494400,
        "proposedMarketValueRaw": {
          "value": 105000,
          "currency": "EUR"
        },
        "fieldTranslations": {
          "nameTranslation": {
            "ar": "أندريا ريسبولي",
            "hi": "एंड्रिया रिस्पोली",
            "bn": "আন্দ্রেয়া রিসপোলি"
          },
          "shortNameTranslation": {
            "ar": "أ. ريسبولي",
            "hi": "ए. रिस्पोली",
            "bn": "এ. রিসপোলি"
          }
        }
      },
      "footballPassingNetworkAction": [],
      "id": 120825871,
      "time": 6,
      "isHome": false,
      "incidentClass": "regular",
      "reversedPeriodTime": 40,
      "incidentType": "goal"
    }
  ],
  "home": {
    "goalkeeperColor": {
      "primary": "ff9900",
      "number": "000000",
      "outline": "ff9900",
      "fancyNumber": "222226"
    },
    "playerColor": {
      "primary": "ffffff",
      "number": "008cea",
      "outline": "ffffff",
      "fancyNumber": "222226"
    }
  },
  "away": {
    "goalkeeperColor": {
      "primary": "ff99ff",
      "number": "000000",
      "outline": "ff99ff",
      "fancyNumber": "222226"
    },
    "playerColor": {
      "primary": "ffb3d9",
      "number": "000000",
      "outline": "ffb3d9",
      "fancyNumber": "222226"
    }
  }
}</pre><div class="json-formatter-container"></div>`
)
