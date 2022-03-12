package soccer

import (
	"fmt"
	"log"
	"soccerapi/src/scores"
	"soccerapi/src/sporxmac"
	"strings"
	"time"
)

type GameScoreWatcher struct {
	ChannelList []string

	//Futbol    sporekraniparser.SporEkraniEvents
	//Futbol    tvmaclari.Matchs
	Futbol    sporxmac.Matchs
	AllScores scores.AllScore
	Scores    *scores.Scores
	GameIDs   map[int]string
}

func (g *GameScoreWatcher) Set() error {
	//res, err := tvmaclari.Get()
	res, err := sporxmac.Get()
	if err != nil {
		return err
	}

	g.Futbol = res.FilterWithSportName("futbol").FilterWithChannelsName(g.ChannelList) // Filter
	log.Println("Filtreli tvMaçları: ", g.Futbol)

	if g.AllScores, err = g.Scores.RequestAllScores(false); err != nil {
		return err
	}

	g.setgameids()

	fmt.Printf("g.Scores.GameIds: %v\n", g.Scores.GameIds)

	return nil
}

func (g *GameScoreWatcher) setgameids() {
	g.GameIDs = map[int]string{}
	g.Scores.GameIds = []string{}

	g.Scores.Updates.RMutex.RLock() // Var olan veri türüne (map?) veri eklemek yerine değişkene yeniden başka bir tanımlama yapıldığı için RLock ile bekleyen bütün Read isteklerini bekletmek daha doğru olacaktır.
	g.Scores.Updates.PlainScore = make(map[int]scores.PlainGameScore)
	g.Scores.Updates.RMutex.RUnlock()

	for _, game := range g.AllScores.Games {
		for _, kanalgame := range g.Futbol {

			if strings.EqualFold(game.HomeCompetitor.Name, kanalgame.HomeName) || strings.EqualFold(game.AwayCompetitor.Name, kanalgame.AwayName) {
				// Maç eşleşti, kanalı eşleştirelim
				for _, channel := range kanalgame.Channels {
					var end = false
					for _, bizimkanal := range g.ChannelList {
						if channel == bizimkanal {
							// Herşey tamam dayı
							end = true
							g.GameIDs[game.ID] = bizimkanal
							break
						}
					}
					if end {
						break
					}
				}
			}
		}
	}

	for key := range g.GameIDs {
		g.Scores.GameIds = append(g.Scores.GameIds, fmt.Sprint(key))
	}

	//g.Scores.GameIds
	/*for _, ga := range g.AllScores.Games {
		g.Scores.GameIds = append(g.Scores.GameIds, fmt.Sprint(ga.ID))
	}*/
	//fmt.Printf("g.Scores.GameIds: %v\n", g.Scores.GameIds)

}

func (g *GameScoreWatcher) loop() {

	g.Scores.Loop()

	for range time.NewTicker(time.Hour).C {
		if err := g.Set(); err != nil {
			time.Sleep(time.Second * 5)
			g.Set()
		}
	}

}

type GameScoreWatcherOptions struct {
	IpTvChannelList []string
	OnGoal          func(scores.GameScore, chan scores.GameScore)
	OnCanceledGoal  func(new scores.GameScore, old scores.GameScore)
	OnRequestError  func(error)
}

func NewGameScoreWatcher(o GameScoreWatcherOptions) (gsw *GameScoreWatcher, err error) {
	gsw = &GameScoreWatcher{
		ChannelList: o.IpTvChannelList,
		GameIDs:     map[int]string{},
		Scores:      scores.NewScores(),
	}

	gsw.Scores.OnWhenGoal = o.OnGoal
	gsw.Scores.OnRequestError = o.OnRequestError
	gsw.Scores.OnWhenGoalCanceled = o.OnCanceledGoal

	if err := gsw.Set(); err != nil {
		return nil, err
	}

	go gsw.loop()
	return
}

//func (g *GameScoreWatcher) update() {}
