package scores

func (s *Scores) CleanGameScoreCaches() {
	for _, i := range s.GameIds {
		s.Database.CleanGameScore(i)
	}
}

func (s *Scores) CleanGameScoreCache(id string) {

	s.Database.CleanGameScore(id)

}
