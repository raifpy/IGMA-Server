package scores

import (
	"log"
	"time"
)

func (s *Scores) Loop() {
	s.AllGameScoresLoop()
	s.GameScoreLoop()
}

func (s *Scores) AllGameScoresLoop() {
	go func() {
		for range time.NewTicker(time.Second * 30).C {
			r, err := s.RequestAllScores(false)
			if err != nil {
				s.OnRequestError(err)
				continue
			}

			s.Database.StoreAllScore(r)
		}
	}()
}

func (s *Scores) GameScoreLoop() {
	go func() {
		for range time.NewTicker(time.Second * 2).C {

			for _, id := range s.GameIds {
				response, err := s.RequestGameScore(id)
				if err != nil {
					s.OnRequestError(err)
					continue
				}
				//fmt.Println(response)
				template := response.ToPlain()

				updatetemplatefunc := func() {
					s.Updates.UpdateID = RandStringRunes(20)
					s.Updates.RMutex.RLock()
					defer s.Updates.RMutex.RUnlock()
					s.Updates.PlainScore[template.Id] = template
				}

				if response.GameGame.StatusText == "Sonuç" ||
					response.GameGame.StatusText == "Bırakıldı" ||
					response.GameGame.StatusText == "Ertelenen" ||
					response.GameGame.StatusText == "İptal Edildi" ||
					response.GameGame.StatusText == "Penaltılardan Sonra" {

					//fmt.Printf("maç %s bitmiş", id)

					s.GameIds = s.GameIds.Delete(id)
					s.Database.CleanGameScore(id)

					template.Finished = true

					go updatetemplatefunc()

					continue
				}

				oldresponse, ok := s.Database.GetGameScore(id)
				if !ok {
					s.Database.StoreGameScore(response)
					go updatetemplatefunc()
					continue
				}

				if oldresponse.LastUpdateID == response.LastUpdateID {
					continue
				}

				//go updatetemplatefunc()
				s.Database.StoreGameScore(response)
				go updatetemplatefunc()

				var oldhomescore = int(oldresponse.GameGame.HomeCompetitor.Score)
				var oldawayscore = int(oldresponse.GameGame.AwayCompetitor.Score)

				var homescore = int(response.GameGame.HomeCompetitor.Score)
				var awayscore = int(response.GameGame.AwayCompetitor.Score)

				if awayscore == 0 && homescore == 0 {
					continue
				}

				if oldhomescore == homescore && oldawayscore == awayscore {
					continue
				}

				if (oldhomescore > homescore) || (oldawayscore > awayscore) {
					log.Printf("Anlaşılan gol iptal: %s %s : %d %d\n", response.GameGame.HomeCompetitor.Name, response.GameGame.AwayCompetitor.Name, homescore, awayscore)
					go s.OnWhenGoalCanceled(response, oldresponse)
					go updatetemplatefunc()
					continue
				}

				if awayscore > oldhomescore || homescore > oldawayscore {
					go updatetemplatefunc()
					log.Printf("Gol olmuş olmalı: %s %s : %d %d\n", response.GameGame.HomeCompetitor.Name, response.GameGame.AwayCompetitor.Name, homescore, awayscore)
					kanal := make(chan GameScore, 1)
					go func() {
						defer close(kanal) //!!Olası panic
						guncelres, err := s.RequestGameScore(id)
						if err != nil {
							s.OnRequestError(err)
							return
						}
						kanal <- guncelres

					}()
					go s.OnWhenGoal(response, kanal)
				}

				time.Sleep(time.Millisecond * 300)

			}

			//fmt.Printf("s.GameIds: %v\n", s.GameIds)
		}
	}()
}
