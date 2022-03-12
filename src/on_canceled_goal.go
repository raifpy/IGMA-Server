package soccer

import (
	"log"
	"soccerapi/src/scores"
)

//! Zannedersem çalışmıyor :)
func (s *Soccer) OnCanceledGoal(_ scores.GameScore, old scores.GameScore) {
	plain := old.ToPlain()
	log.Printf("\033[31mMaç %s %s gol iptal\033[0m\n", plain.HomeName, plain.AwayName)
	clips, err := s.GetMatchClips(int64(plain.Id))
	if err != nil {
		log.Println("Gol iptal olmuş ya, onu canceled_goal yapacaktım ama zaten yazılmamış :')")
		return
	}

	for _, v := range clips.Clips {
		if v.Client.HomeScore == plain.HomeScore && v.Client.AwayScore == plain.AwayScore {
			v.Client.Status = "canceled_goal"
		}
	}

	s.ReplaceMatchClips(clips)
}
