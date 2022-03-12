package scores

import (
	"encoding/json"
	"time"
)

type GameScore struct {
	LastUpdateID      int64    `json:"lastUpdateId"`
	RequestedUpdateID int      `json:"requestedUpdateId"`
	GameGame          GameGame `json:"game"`
}

func (g GameScore) ToPlain() (p PlainGameScore) {
	p.HomeName = g.GameGame.HomeCompetitor.Name
	p.AwayName = g.GameGame.AwayCompetitor.Name

	p.HomeScore = int(g.GameGame.HomeCompetitor.Score)
	p.AwayScore = int(g.GameGame.AwayCompetitor.Score)

	p.Id = g.GameGame.ID
	p.UpdateID = g.LastUpdateID
	p.LeagueName = g.GameGame.CompetitionDisplayName
	p.Minute = g.GameGame.GameTimeDisplay

	if g.GameGame.StatusText == "Sonuç" ||
		g.GameGame.StatusText == "Bırakıldı" ||
		g.GameGame.StatusText == "Ertelenen" ||
		g.GameGame.StatusText == "İptal Edildi" ||
		g.GameGame.StatusText == "Penaltılardan Sonra" {
		p.Finished = true
		//p.MatchStatus = "finished"
	} else {

		p.MatchStatus = "live"
		if p.HomeScore == -1 || p.Minute == "" {
			p.MatchStatus = "not_started"
		}
	}

	return
}

func (a GameScore) ToJson() []byte {
	b, _ := json.Marshal(a)
	return b
}

func (a GameScore) ToPrettyJson() []byte {
	b, _ := json.MarshalIndent(a, "", " ")
	return b
}

type PreciseGameTime struct {
	Minutes        int  `json:"minutes"`
	Seconds        int  `json:"seconds"`
	AutoProgress   bool `json:"autoProgress"`
	ClockDirection int  `json:"clockDirection"`
}

type EventType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	SubTypeID   int    `json:"subTypeId"`
	SubTypeName string `json:"subTypeName"`
}
type Events []Event
type Event struct {
	CompetitorID                 int       `json:"competitorId"`
	StatusID                     int       `json:"statusId"`
	StageID                      int       `json:"stageId"`
	Order                        int       `json:"order"`
	Num                          int       `json:"num"`
	GameTime                     float64   `json:"gameTime"`
	AddedTime                    int       `json:"addedTime"`
	GameTimeDisplay              string    `json:"gameTimeDisplay"`
	GameTimeAndStatusDisplayType int       `json:"gameTimeAndStatusDisplayType"`
	PlayerID                     int       `json:"playerId"`
	IsMajor                      bool      `json:"isMajor"`
	EventType                    EventType `json:"eventType,omitempty"`

	ExtraPlayers []int `json:"extraPlayers,omitempty"`
}

type GameGame struct {
	ID                           int             `json:"id"`
	SportID                      int             `json:"sportId"`
	CompetitionID                int             `json:"competitionId"`
	SeasonNum                    int             `json:"seasonNum"`
	StageNum                     int             `json:"stageNum"`
	RoundNum                     int             `json:"roundNum"`
	CompetitionDisplayName       string          `json:"competitionDisplayName"`
	StartTime                    time.Time       `json:"startTime"`
	StatusGroup                  int             `json:"statusGroup"`
	StatusText                   string          `json:"statusText"`
	ShortStatusText              string          `json:"shortStatusText"`
	GameTimeAndStatusDisplayType int             `json:"gameTimeAndStatusDisplayType"`
	JustEnded                    bool            `json:"justEnded"`
	GameTime                     float64         `json:"gameTime"`
	GameTimeDisplay              string          `json:"gameTimeDisplay"`
	PreciseGameTime              PreciseGameTime `json:"preciseGameTime"`
	HasLineups                   bool            `json:"hasLineups"`
	HasMissingPlayers            bool            `json:"hasMissingPlayers"`
	HasFieldPositions            bool            `json:"hasFieldPositions"`
	HomeCompetitor               HomeCompetitor  `json:"homeCompetitor"`
	AwayCompetitor               AwayCompetitor  `json:"awayCompetitor"`

	Events Events `json:"events"`

	PreviousMeetings []int `json:"previousMeetings"`
}
