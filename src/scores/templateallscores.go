package scores

import (
	"encoding/json"
	"time"
)

type AllScore struct {
	LastUpdateID      int64        `json:"lastUpdateId"`
	RequestedUpdateID int          `json:"requestedUpdateId"`
	TTL               int          `json:"ttl"`
	Sports            Sports       `json:"sports"`
	Countries         Countries    `json:"countries"`
	Competitions      Competitions `json:"competitions"`
	Competitors       Competitors  `json:"competitors"`
	Games             Games        `json:"games"`
	LiveGamesCount    int          `json:"liveGamesCount"`
}

func (a AllScore) ToJson() []byte {
	b, _ := json.Marshal(a)
	return b
}

func (a AllScore) ToPrettyJson() []byte {
	b, _ := json.MarshalIndent(a, "", " ")
	return b
}

type Sports []Sport
type Sport struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NameForURL   string `json:"nameForURL"`
	DrawSupport  bool   `json:"drawSupport"`
	TotalGames   int    `json:"totalGames"`
	LiveGames    int    `json:"liveGames"`
	ImageVersion int    `json:"imageVersion"`
}
type Countries []Countrie
type Countrie struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TotalGames   int    `json:"totalGames,omitempty"`
	LiveGames    int    `json:"liveGames,omitempty"`
	NameForURL   string `json:"nameForURL"`
	SportTypes   []int  `json:"sportTypes"`
	ImageVersion int    `json:"imageVersion"`
}

type Competitors []Competitor
type Competitor struct {
	ID             int    `json:"id"`
	CountryID      int    `json:"countryId"`
	SportID        int    `json:"sportId"`
	Name           string `json:"name"`
	NameForURL     string `json:"nameForURL"`
	Type           int    `json:"type"`
	PopularityRank int    `json:"popularityRank"`
	ImageVersion   int    `json:"imageVersion"`
	LongName       string `json:"longName,omitempty"`
}

type Games []Game
type Game struct {
	ID                           int            `json:"id"`
	SportID                      int            `json:"sportId"`
	CompetitionID                int            `json:"competitionId"`
	SeasonNum                    int            `json:"seasonNum"`
	StageNum                     int            `json:"stageNum"`
	RoundNum                     int            `json:"roundNum"`
	CompetitionDisplayName       string         `json:"competitionDisplayName"`
	StartTime                    time.Time      `json:"startTime"`
	StatusGroup                  int            `json:"statusGroup"`
	StatusText                   string         `json:"statusText"`
	ShortStatusText              string         `json:"shortStatusText"`
	GameTimeAndStatusDisplayType int            `json:"gameTimeAndStatusDisplayType"`
	JustEnded                    bool           `json:"justEnded"`
	GameTime                     float64        `json:"gameTime"`
	GameTimeDisplay              string         `json:"gameTimeDisplay"`
	HasLineups                   bool           `json:"hasLineups,omitempty"`
	HasMissingPlayers            bool           `json:"hasMissingPlayers,omitempty"`
	HasFieldPositions            bool           `json:"hasFieldPositions,omitempty"`
	HasTVNetworks                bool           `json:"hasTVNetworks"`
	HasBetsTeaser                bool           `json:"hasBetsTeaser,omitempty"`
	AwayCompetitor               AwayCompetitor `json:"awayCompetitor,omitempty"`
	HomeCompetitor               HomeCompetitor `json:"homeCompetitor,omitempty"`
}
