package scores

import (
	"net/http"
	"sync"
)

var DefaultAllScoresApiURl ReplaceString = "https://webws.365scores.com/web/games/allscores/?langId=33&timezoneName=Europe/Istanbul&userCountryId=12&sports=1&onlyLiveGames={live}"
var DefaultGameApiURL ReplaceString = "https://webws.365scores.com/web/game/?appTypeId=5&langId=33&timezoneName=Europe/Istanbul&userCountryId=12&gameId={gameid}"

type Scores struct {
	Updates            *MatchUpdates
	OnWhenGoal         func(GameScore, chan GameScore)
	OnWhenGoalCanceled func(new GameScore, old GameScore)
	OnRequestError     func(error)
	GameIds            EasyList
	Database           Database
	Client             *http.Client
	AllScoresApiURL    ReplaceString
	GameScoreApiURL    ReplaceString
}

func NewScores() *Scores {
	return &Scores{
		Database:        NewSyncMapDatabase(),
		Client:          http.DefaultClient,
		AllScoresApiURL: DefaultAllScoresApiURl,
		GameScoreApiURL: DefaultGameApiURL,
		Updates: &MatchUpdates{
			RMutex: &sync.RWMutex{},
		},
	}
}
