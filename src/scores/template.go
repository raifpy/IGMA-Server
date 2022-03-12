package scores

type HomeCompetitor struct {
	ID             int     `json:"id"`
	CountryID      int     `json:"countryId"`
	SportID        int     `json:"sportId"`
	Name           string  `json:"name"`
	Score          float64 `json:"score"`
	IsQualified    bool    `json:"isQualified"`
	ToQualify      bool    `json:"toQualify"`
	IsWinner       bool    `json:"isWinner"`
	RedCards       int     `json:"redCards"`
	NameForURL     string  `json:"nameForURL"`
	Type           int     `json:"type"`
	PopularityRank int     `json:"popularityRank"`
	Outcome        int     `json:"outcome"`
	ImageVersion   int     `json:"imageVersion"`
}
type AwayCompetitor struct {
	ID             int     `json:"id"`
	CountryID      int     `json:"countryId"`
	SportID        int     `json:"sportId"`
	Name           string  `json:"name"`
	Score          float64 `json:"score"`
	IsQualified    bool    `json:"isQualified"`
	ToQualify      bool    `json:"toQualify"`
	IsWinner       bool    `json:"isWinner"`
	RedCards       int     `json:"redCards"`
	NameForURL     string  `json:"nameForURL"`
	Type           int     `json:"type"`
	PopularityRank int     `json:"popularityRank"`
	Outcome        int     `json:"outcome"`
	ImageVersion   int     `json:"imageVersion"`
}

type Competitions []Competition
type Competition struct {
	ID                 int    `json:"id"`
	CountryID          int    `json:"countryId"`
	SportID            int    `json:"sportId"`
	Name               string `json:"name"`
	HasStandings       bool   `json:"hasStandings,omitempty"`
	HasStandingsGroups bool   `json:"hasStandingsGroups,omitempty"`
	HasBrackets        bool   `json:"hasBrackets"`
	NameForURL         string `json:"nameForURL"`
	TotalGames         int    `json:"totalGames"`
	LiveGames          int    `json:"liveGames"`
	PopularityRank     int    `json:"popularityRank"`
	HasActiveGames     bool   `json:"hasActiveGames,omitempty"`
	ImageVersion       int    `json:"imageVersion"`
	HasLiveStandings   bool   `json:"hasLiveStandings,omitempty"`
	ShortName          string `json:"shortName,omitempty"`
	LongName           string `json:"longName,omitempty"`
}

type PlainGameScore struct {
	Id          int    `json:"id"`
	MatchStatus string `json:"match_status ,omitempty"` // live - not_started ? finished
	HomeName    string `json:"home_name"`
	AwayName    string `json:"away_name"`

	LeagueName string `json:"league_name"`

	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`

	Minute string `json:"minute"`

	Finished bool `json:"finished"`

	UpdateID int64 `json:"update_id"`
}
