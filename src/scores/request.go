package scores

import (
	"encoding/json"
	"fmt"
)

func (s *Scores) RequestAllScores(live bool) (sc AllScore, err error) {
	res, err := s.Client.Get(s.AllScoresApiURL.Replace(map[string]string{"live": fmt.Sprint(live)}))
	if err != nil {
		return sc, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&sc)
	/*f, _ := os.Create("outraw.json")
	io.Copy(f, res.Body)
	f.Close()*/

	return
}

func (s *Scores) RequestGameScore(gameid string) (gc GameScore, err error) {

	res, err := s.Client.Get(s.GameScoreApiURL.Replace(map[string]string{
		"gameid": gameid,
	}))

	if err != nil {
		return GameScore{}, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&gc)

	return
}
